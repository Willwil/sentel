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
)

type clusterEventManager struct {
	config         core.Config
	brokerId       string
	consumer       sarama.Consumer                // kafka client endpoint
	eventChan      chan *Event                    // Event notify channel
	subscribers    map[uint32][]subscriberContext // All subscriber's contex
	mutex          sync.Mutex                     // Mutex to lock subscribers
	quit           chan os.Signal
	waitgroup      sync.WaitGroup
	product        string
	mqttEventBus   string
	brokerEventBus string
}

const (
	mqttEventBusName   = "broker/event/mqtt"
	brokerEventBusName = "broker/event/broker"
)

// newclusterEventManager create global broker
func newClusterEventManager(product string, c core.Config) (*clusterEventManager, error) {
	// kafka
	khosts, _ := core.GetServiceEndpoint(c, "broker", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &clusterEventManager{
		config:         c,
		consumer:       consumer,
		eventChan:      make(chan *Event),
		subscribers:    make(map[uint32][]subscriberContext),
		quit:           make(chan os.Signal),
		waitgroup:      sync.WaitGroup{},
		mqttEventBus:   product + "/event/mqtt",
		brokerEventBus: product + "/event/broker",
	}, nil
}

// initialize
func (p *clusterEventManager) initialize(c core.Config) error {
	// subscribe topic from kafka
	if err := p.subscribeKafkaTopic(p.mqttEventBus); err != nil {
		return err
	}
	if err := p.subscribeKafkaTopic(p.brokerEventBus); err != nil {
		return err
	}
	return nil
}

// Start
func (p *clusterEventManager) run() error {
	go func(p *clusterEventManager) {
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
					core.AsyncProduceMessage(p.config, "event", mqttEventBusName, e)
				} else {
					core.AsyncProduceMessage(p.config, "event", brokerEventBusName, e)
				}
			case <-p.quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *clusterEventManager) stop() {
	close(p.eventChan)
	p.consumer.Close()
}

// handleEvent handle broker event from other service
func (p *clusterEventManager) handleclusterEventManagerEvent(e *Event) {
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
func (p *clusterEventManager) handleMqttEvent(e *Event) {
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

// subscribeKafkaTopc subscribe topics from apiserver
func (p *clusterEventManager) subscribeKafkaTopic(topic string) error {
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
func (p *clusterEventManager) handleNotifications(topic string, value []byte) error {
	obj := &Event{}
	if err := json.Unmarshal(value, obj); err != nil {
		return err
	}
	switch topic {
	case mqttEventBusName:
		p.handleMqttEvent(obj)
	case brokerEventBusName:
		p.handleclusterEventManagerEvent(obj)
	}
	return nil
}

// publish publish event to event service
func (p *clusterEventManager) notify(e *Event) {
	p.eventChan <- e
}

// subscribe subcribe event from event service
func (p *clusterEventManager) subscribe(t uint32, handler EventHandler, ctx interface{}) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.subscribers[t]; !ok {
		// Create handler list if not existed
		p.subscribers[t] = []subscriberContext{}
	}
	p.subscribers[t] = append(p.subscribers[t], subscriberContext{handler: handler, ctx: ctx})
}
