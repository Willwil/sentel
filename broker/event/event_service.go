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
//  under the License.

package event

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/core"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type EventService struct {
	core.ServiceBase
	consumer    sarama.Consumer               // kafka client endpoint
	eventChan   chan *Event                   // Event notify channel
	subscribers map[uint8][]subscriberContext // All subscriber's contex
	mutex       sync.Mutex                    // Mutex to lock subscribers
}

// subscriberContext hold subscriber handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

const (
	ServiceName = "event"
)

// EventServiceFactory
type EventServiceFactory struct{}

// New create metadata service factory
func (p *EventServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// kafka
	khosts, _ := core.GetServiceEndpoint(c, "broker", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &EventService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumer:    consumer,
		eventChan:   make(chan *Event),
		subscribers: make(map[uint8][]subscriberContext),
	}, nil

}

// Name
func (p *EventService) Name() string {
	return ServiceName
}

// Start
func (p *EventService) Start() error {
	if err := p.subscribeTopic(EventBusName); err != nil {
		return err
	}
	go func(p *EventService) {
		for {
			select {
			case e := <-p.eventChan:
				core.AsyncProduceMessage(p.Config, "event", EventBusName, e)
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *EventService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	close(p.eventChan)
	p.consumer.Close()
}

// handleEvent handle event from other service
func (p *EventService) handleEvent(e *Event) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.subscribers[e.Type]; !ok {
		return
	}
	subscribers := p.subscribers[e.Type]
	for _, subscriber := range subscribers {
		go subscriber.handler(e, subscriber.ctx)
	}
}

// subscribeTopc subscribe topics from apiserver
func (p *EventService) subscribeTopic(topic string) error {
	partitionList, err := p.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := p.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		p.WaitGroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.WaitGroup.Done()
			for msg := range pc.Messages() {
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}

// handleNotifications handle notification from kafka
func (p *EventService) handleNotifications(topic string, value []byte) error {
	switch topic {
	case EventBusName:
		obj := &Event{}
		if err := json.Unmarshal(value, obj); err != nil {
			return err
		}
		p.handleEvent(obj)
	}

	return nil
}

// publish publish event to event service
func (p *EventService) notify(e *Event) {
	p.eventChan <- e
}

// subscribe subcribe event from event service
func (p *EventService) subscribe(t uint8, handler EventHandler, ctx interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.subscribers[t]; !ok {
		// Create handler list if not existed
		p.subscribers[t] = []subscriberContext{}
	}
	p.subscribers[t] = append(p.subscribers[t], subscriberContext{handler: handler, ctx: ctx})
}
