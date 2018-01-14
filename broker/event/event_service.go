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
	fmtOfMqttEventBus   = "broker-%s-event-mqtt"
	fmtOfBrokerEventBus = "broker-%s-event-broker"
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
	eventChan   chan *Event                    // Event notify channel
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
	tenant := c.MustString("broker", "tenant")
	khosts, err := c.String("broker", "kafka")
	if err != nil || khosts == "" {
		return nil, errors.New("message server is not rightly configed")
	}
	consumer, err1 := message.NewConsumer(khosts, base.GetBrokerId())
	producer, err2 := message.NewProducer(khosts, base.GetBrokerId(), true)
	if err1 != nil || err2 != nil {
		return nil, errors.New("message client endpoint failed")
	}
	return &eventService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		consumer:    consumer,
		producer:    producer,
		eventChan:   make(chan *Event),
		subscribers: make(map[uint32][]subscriberContext),
		tenant:      tenant,
		quitChan:    make(chan interface{}),
	}, nil
}

func (p *eventService) nameOfEventBus(e *Event) string {
	switch e.Type {
	case SessionCreate, SessionDestroy, TopicPublish, TopicSubscribe, TopicUnsubscribe:
		return fmt.Sprintf(fmtOfMqttEventBus, p.tenant)
	default:
		return fmt.Sprintf(fmtOfBrokerEventBus, p.tenant)
	}
}

// initialize
func (p *eventService) Initialize() error { return nil }

func (p *eventService) Name() string { return ServiceName }
func (p *eventService) bootstrap() error {
	return nil
}

// messageHandlerFunc handle mqtt event from other service
func (p *eventService) messageHandlerFunc(topic string, value []byte, ctx interface{}) {
	re := &RawEvent{}
	if err := json.Unmarshal(value, re); err == nil {
		e := &Event{}
		if err := json.Unmarshal(re.Header, &e.EventHeader); err != nil {
			glog.Errorf("event service unmarshal event common failed:%s", err.Error())
			return
		}
		// cluster event manager only handle kafka event from other broker
		// Iterate all subscribers to notify
		if e.BrokerId != base.GetBrokerId() {
			if _, found := p.subscribers[e.Type]; found {
				subscribers := p.subscribers[e.Type]
				err := p.unmarshalEventPayload(e, re.Payload)
				if err != nil {
					glog.Errorf(err.Error())
					return
				}
				for _, subscriber := range subscribers {
					subscriber.handler(e, subscriber.ctx)
				}
			}
		}
	}
}

// Start
func (p *eventService) Start() error {
	err1 := p.consumer.Subscribe(fmt.Sprintf(fmtOfMqttEventBus, p.tenant), p.messageHandlerFunc, nil)
	err2 := p.consumer.Subscribe(fmt.Sprintf(fmtOfBrokerEventBus, p.tenant), p.messageHandlerFunc, nil)
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
				if _, found := p.subscribers[e.Type]; found {
					subscribers := p.subscribers[e.Type]
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

func (p *eventService) publishMsg(topic string, e *Event) error {
	re := &RawEvent{}
	re.Header, _ = json.Marshal(e.EventHeader)
	re.Payload, _ = json.Marshal(e.Detail)
	err := p.producer.SendMessage("iot-broker", topic, &re)
	if err != nil {
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

func (p *eventService) unmarshalEventPayload(e *Event, payload json.RawMessage) error {
	switch e.Type {
	case SessionCreate:
		r := SessionCreateDetail{}
		if err := json.Unmarshal(payload, &r); err != nil {
			return fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = r
	case SessionDestroy:
		r := SessionDestroyDetail{}
		if err := json.Unmarshal(payload, &r); err != nil {
			return fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = r
	case TopicSubscribe:
		r := TopicSubscribeDetail{}
		if err := json.Unmarshal(payload, &r); err != nil {
			return fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = r
	case TopicUnsubscribe:
		r := TopicUnsubscribeDetail{}
		if err := json.Unmarshal(payload, &r); err != nil {
			return fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = r
	case TopicPublish:
		r := TopicPublishDetail{}
		if err := json.Unmarshal(payload, &r); err != nil {
			return fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		e.Detail = r
	default:
		return fmt.Errorf("invalid event type '%d'", e.Type)
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
