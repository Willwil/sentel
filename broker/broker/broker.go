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

package broker

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

type Broker struct {
	core.ServiceManager
	brokerId    string
	consumer    sarama.Consumer                // kafka client endpoint
	eventChan   chan *Event                    // Event notify channel
	subscribers map[uint32][]subscriberContext // All subscriber's contex
	mutex       sync.Mutex                     // Mutex to lock subscribers
	quit        chan os.Signal
	waitgroup   sync.WaitGroup
}

// subscriberContext hold subscriber handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

const (
	BrokerVersion      = "0.1"
	mqttEventBusName   = "broker/event/mqtt"
	brokerEventBusName = "broker/event/broker"
)

var (
	_broker *Broker
)

// newBroker create global broker
func NewBroker(c core.Config) (*Broker, error) {
	if _broker != nil {
		panic("Global broker had already been created")
	}

	// kafka
	khosts, _ := core.GetServiceEndpoint(c, "broker", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	serviceMgr, err := core.NewServiceManager("broker", c)
	if err != nil {
		return nil, err
	}
	_broker = &Broker{
		ServiceManager: *serviceMgr,
		brokerId:       uuid.NewV4().String(),
		consumer:       consumer,
		eventChan:      make(chan *Event),
		subscribers:    make(map[uint32][]subscriberContext),
		quit:           make(chan os.Signal),
		waitgroup:      sync.WaitGroup{},
	}
	return _broker, nil
}

// GetBroker create service manager and all supported service
// The function should be called in service
func GetBroker() *Broker { return _broker }

// Version
func GetVersion() string {
	return BrokerVersion
}

// GetBrokerId return broker's identifier
func GetId() string {
	return _broker.brokerId
}

// GetService return specified service instance
func GetService(name string) core.Service {
	broker := GetBroker()
	return broker.GetServiceByName(name)
}

// GetConfig return broker's configuration
func (p *Broker) GetConfig() core.Config {
	return p.Config
}

// GetServicesByName return service instance by name, or matched by part of name
func (p *Broker) GetServiceByName(name string) core.Service {
	if _, ok := p.Services[name]; !ok {
		panic(fmt.Sprintf("Failed to find service '%s' in broker", name))
	}
	return p.Services[name]
}

// Start
func (p *Broker) Run() error {
	// launch other services at first
	if err := p.ServiceManager.Run(); err != nil {
		return err
	}
	// subscribe topic from kafka
	if err := p.subscribeTopic(mqttEventBusName); err != nil {
		return err
	}
	if err := p.subscribeTopic(brokerEventBusName); err != nil {
		return err
	}

	go func(p *Broker) {
		for {
			select {
			case e := <-p.eventChan:
				// When event is received from local service, we should
				// transeverse other service to process it at first
				if _, found := p.subscribers[e.Type]; found {
					subscribers := p.subscribers[e.Type]
					for _, subscriber := range subscribers {
						go subscriber.handler(e, subscriber.ctx)
					}
				}
				// notify kafka for broker synchronization
				if (e.Type & 0xff) != 0 {
					core.AsyncProduceMessage(p.Config, "event", mqttEventBusName, e)
				} else {
					core.AsyncProduceMessage(p.Config, "event", brokerEventBusName, e)
				}
			case <-p.quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *Broker) Stop() {
	close(p.eventChan)
	p.consumer.Close()
}

// handleEvent handle broker event from other service
func (p *Broker) handleBrokerEvent(e *Event) {
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

// handleEvent handle mqtt event from other service
func (p *Broker) handleMqttEvent(e *Event) {
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
func (p *Broker) subscribeTopic(topic string) error {
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
		p.waitgroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.waitgroup.Done()
			for msg := range pc.Messages() {
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}

// handleNotifications handle notification from kafka
func (p *Broker) handleNotifications(topic string, value []byte) error {
	obj := &Event{}
	if err := json.Unmarshal(value, obj); err != nil {
		return err
	}
	switch topic {
	case mqttEventBusName:
		p.handleMqttEvent(obj)
	case brokerEventBusName:
		p.handleBrokerEvent(obj)
	}
	return nil
}

// publish publish event to event service
func (p *Broker) notify(e *Event) {
	p.eventChan <- e
}

// subscribe subcribe event from event service
func (p *Broker) subscribe(t uint32, handler EventHandler, ctx interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.subscribers[t]; !ok {
		// Create handler list if not existed
		p.subscribers[t] = []subscriberContext{}
	}
	p.subscribers[t] = append(p.subscribers[t], subscriberContext{handler: handler, ctx: ctx})
}
