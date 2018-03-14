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
	"errors"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

const (
	ServiceName         = "event"
	FmtOfMqttEventBus   = "iot-%s-event-mqtt"
	FmtOfBrokerEventBus = "iot-%s-event-broker"
)

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

type eventService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	brokerId    string
	consumer    message.Consumer
	producer    message.Producer
	eventChan   chan Event                     // Event notify channel
	subscribers map[uint32][]subscriberContext // All subscriber's contex
	mutex       sync.Mutex                     // Mutex to lock subscribers
	tenant      string
	quitChan    chan interface{}
}

// subscriberContext hold subscribre's handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

type ServiceFactory struct{}

// New create global broker
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// Retrieve tenant and product
	tenant := c.MustString("tenant")
	consumer, err1 := message.NewConsumer(c, base.GetBrokerId())
	producer, err2 := message.NewProducer(c, base.GetBrokerId(), true)
	if err1 != nil || err2 != nil {
		return nil, errors.New("message client endpoint failed")
	}
	return &eventService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		consumer:    consumer,
		producer:    producer,
		eventChan:   make(chan Event),
		subscribers: make(map[uint32][]subscriberContext),
		tenant:      tenant,
		quitChan:    make(chan interface{}),
	}, nil
}

func (p *eventService) nameOfEventBus(e Event) string {
	switch e.GetType() {
	case SessionCreate, SessionDestroy, TopicPublish, TopicSubscribe, TopicUnsubscribe:
		return fmt.Sprintf(FmtOfMqttEventBus, p.tenant)
	default:
		return fmt.Sprintf(FmtOfBrokerEventBus, p.tenant)
	}
}

// initialize
func (p *eventService) Initialize() error { return nil }
func (p *eventService) Name() string      { return ServiceName }
func (p *eventService) bootstrap() error  { return nil }

// messageHandlerFunc handle mqtt event from other service
func (p *eventService) messageHandlerFunc(msg message.Message, ctx interface{}) {
	if e, err := Decode(msg, JSONCodec); err == nil {
		// cluster event manager only handle kafka event from other broker
		// Iterate all subscribers to notify
		etype := e.GetType()
		if e.GetBrokerId() != base.GetBrokerId() {
			if _, found := p.subscribers[etype]; found {
				subscribers := p.subscribers[etype]
				for _, subscriber := range subscribers {
					subscriber.handler(e, subscriber.ctx)
				}
			}
		}
	}
}

// Start
func (p *eventService) Start() error {
	err1 := p.consumer.Subscribe(fmt.Sprintf(FmtOfMqttEventBus, p.tenant), p.messageHandlerFunc, nil)
	err2 := p.consumer.Subscribe(fmt.Sprintf(FmtOfBrokerEventBus, p.tenant), p.messageHandlerFunc, nil)
	if err1 != nil || err2 != nil {
		return errors.New("subscribe topic failed")
	}
	if err := p.consumer.Start(); err != nil {
		return err
	}
	p.waitgroup.Add(1)
	go func(p *eventService) {
		defer p.waitgroup.Done()
		for {
			select {
			case e := <-p.eventChan:
				// When event is received from local service, we should
				// transeverse other service to process it at first
				if _, found := p.subscribers[e.GetType()]; found {
					subscribers := p.subscribers[e.GetType()]
					for _, subscriber := range subscribers {
						subscriber.handler(e, subscriber.ctx)
					}
				}
				// notify kafka for broker synchronization
				p.publishMsg(p.nameOfEventBus(e), e)
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

func (p *eventService) publishMsg(topic string, e Event) error {
	value, _ := Encode(e, JSONCodec)
	msg := &message.Broker{TopicName: topic, Payload: value}
	if err := p.producer.SendMessage(msg); err != nil {
		glog.Errorf("Failed to store your data:, %s", err)
		return err
	}
	return nil
}

// Stop
func (p *eventService) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	p.producer.Close()
	p.consumer.Close()
	close(p.eventChan)
	close(p.quitChan)
}

// publish publish event to event service
func (p *eventService) notify(e Event) {
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
