//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations

package event

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

const (
	ServiceName = "event"
)

type eventService struct {
	base.ServiceBase
	brokerId    string
	consumer    sarama.Consumer                // kafka client endpoint
	eventChan   chan *Event                    // Event notify channel
	subscribers map[uint32][]subscriberContext // All subscriber's contex
	mutex       sync.Mutex                     // Mutex to lock subscribers
	quit        chan os.Signal
	waitgroup   sync.WaitGroup
	product     string
}

// subscriberContext hold subscribre's handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

func busNameOfEvent(e *Event) string {
	switch e.Type {
	case SessionCreate, SessionDestroy, TopicPublish, TopicSubscribe, TopicUnsubscribe:
		return "broker/event/mqtt"
	default:
		return "broker/event/broker"
	}
}

// neweventService create global broker
func New(c core.Config, quit chan os.Signal) (base.Service, error) {
	khosts, _ := core.GetServiceEndpoint(c, "broker", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &eventService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumer:    consumer,
		eventChan:   make(chan *Event),
		subscribers: make(map[uint32][]subscriberContext),
		quit:        make(chan os.Signal),
		waitgroup:   sync.WaitGroup{},
	}, nil
}

// initialize
func (p *eventService) initialize(c core.Config) error {
	// TODO: eventService must be associated with tenant/product
	mqttEventBus := "/event/mqtt"
	brokerEventBus := "/event/broker"
	// subscribe topic from kafka
	if err := p.subscribeKafkaTopic(mqttEventBus); err != nil {
		return err
	}
	if err := p.subscribeKafkaTopic(brokerEventBus); err != nil {
		return err
	}
	return nil
}

// Name
func (p *eventService) Name() string {
	return ServiceName
}

func (p *eventService) Initialize() error {
	// bootstrap
	return nil
}

// Start
func (p *eventService) Start() error {
	go func(p *eventService) {
		for {
			select {
			case e := <-p.eventChan:
				// When event is received from local service, we should
				// transeverse other service to process it at first
				if _, found := p.subscribers[e.Type]; found {
					subscribers := p.subscribers[e.Type]
					for _, subscriber := range subscribers {
						subscriber.handler(e, subscriber.ctx)
					}
				}
				// notify kafka for broker synchronization
				core.AsyncProduceMessage(p.Config, "event", busNameOfEvent(e), e)
			case <-p.quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *eventService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	close(p.eventChan)
	p.consumer.Close()
	p.WaitGroup.Wait()
}

// handleEvent handle mqtt event from other service
func (p *eventService) handleKafkaEvent(e *Event) {
	// cluster event manager only handle kafka event from other broker
	// Iterate all subscribers to notify
	if e.BrokerId != base.GetBrokerId() {
		if _, found := p.subscribers[e.Type]; found {
			subscribers := p.subscribers[e.Type]
			for _, subscriber := range subscribers {
				subscriber.handler(e, subscriber.ctx)
			}
		}
	}
}

// subscribeKafkaTopc subscribe topics from apiserver
func (p *eventService) subscribeKafkaTopic(topic string) error {
	partitionList, err := p.consumer.Partitions(topic)
	if err != nil {
		return err
	}

	for partition := range partitionList {
		pc, err := p.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		p.waitgroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.waitgroup.Done()
			for msg := range pc.Messages() {
				obj := &Event{}
				if err := json.Unmarshal(msg.Value, obj); err == nil {
					p.handleKafkaEvent(obj)
				}
			}
		}(pc)
	}
	return nil
}

// publish publish event to event service
func (p *eventService) notify(e *Event) {
	p.eventChan <- e
}

// subscribe subcribe event from event service
func (p *eventService) subscribe(t uint32, handler EventHandler, ctx interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.subscribers[t]; !ok {
		// Create handler list if not existed
		p.subscribers[t] = []subscriberContext{}
	}
	p.subscribers[t] = append(p.subscribers[t], subscriberContext{handler: handler, ctx: ctx})
}
