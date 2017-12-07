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
)

const (
	ServiceName         = "event"
	fmtOfMqttEventBus   = "broker-%s-%s-event-mqtt"
	fmtOfBrokerEventBus = "broker-%s-%s-event-broker"
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
	tenant      string
	product     string
	deploy      string
}

// subscriberContext hold subscribre's handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

// neweventService create global broker
func New(c core.Config, quit chan os.Signal) (base.Service, error) {
	// Retrieve tenant and product
	tenant := c.MustString("broker", "tenant")
	product := c.MustString("broker", "product")

	var consumer sarama.Consumer
	khosts, err := core.GetServiceEndpoint(c, "broker", "kafka")
	if err == nil && khosts != "" {
		cc, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
		if err != nil {
			return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
		}
		consumer = cc
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
		tenant:      tenant,
		product:     product,
	}, nil
}

func (p *eventService) nameOfEventBus(e *Event) string {
	switch e.Type {
	case SessionCreate, SessionDestroy, TopicPublish, TopicSubscribe, TopicUnsubscribe:
		return fmt.Sprintf(fmtOfMqttEventBus, p.tenant, p.product)
	default:
		return fmt.Sprintf(fmtOfBrokerEventBus, p.tenant, p.product)
	}
}

// initialize
func (p *eventService) Initialize() error {
	mqttEventBus := fmt.Sprintf(fmtOfMqttEventBus, p.tenant, p.product)
	brokerEventBus := fmt.Sprintf(fmtOfBrokerEventBus, p.tenant, p.product)
	// subscribe topic from kafka
	if p.consumer != nil {
		err := p.subscribeKafkaTopic(mqttEventBus)
		err = p.subscribeKafkaTopic(brokerEventBus)
		return err
	}
	return nil
}

func (p *eventService) Name() string { return ServiceName }
func (p *eventService) bootstrap() error {
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
				if p.consumer != nil {
					core.AsyncProduceMessage(p.Config, "event", p.nameOfEventBus(e), e)
				}
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
	if p.consumer != nil {
		p.consumer.Close()
	}
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
	if partitionList, err := p.consumer.Partitions(topic); err == nil {
		for partition, _ := range partitionList {
			pc, err := p.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				return fmt.Errorf("event service subscribe kafka topic failed:%s", err.Error())
			}
			p.waitgroup.Add(1)

			go func(sarama.PartitionConsumer) {
				defer pc.AsyncClose()
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
	return fmt.Errorf("event service failed to subscribe kafka topic '%s'", topic)
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
