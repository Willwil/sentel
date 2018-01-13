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
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/pkg/config"
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
	consumers   map[string]sarama.Consumer // kafka client endpoint
	producer    sarama.SyncProducer
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

	var consumers = make(map[string]sarama.Consumer)
	var producer sarama.SyncProducer
	sarama.Logger = logger
	if khosts := c.MustString("broker", "kafka"); khosts != "" {
		config := sarama.NewConfig()
		config.ClientID = base.GetBrokerId() + "_MQTT"
		// config.Consumer.MaxWaitTime = time.Duration(5 * time.Second)
		// config.Consumer.Offsets.CommitInterval = 1 * time.Second
		// config.Consumer.Offsets.Initial = sarama.OffsetNewest
		cm, err := sarama.NewConsumer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to connect with kafka server(MQTT) '%s'", khosts)
		}
		consumers[fmt.Sprintf(fmtOfMqttEventBus, tenant)] = cm

		config = sarama.NewConfig()
		config.ClientID = base.GetBrokerId() + "_Broker"
		cb, err := sarama.NewConsumer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to connect with kafka server(Broker) '%s'", khosts)
		}
		consumers[fmt.Sprintf(fmtOfBrokerEventBus, tenant)] = cb
		glog.Infof("event service connected with kafka '%s' with broker id '%s'", khosts, base.GetBrokerId())

		config = sarama.NewConfig()
		config.ClientID = base.GetBrokerId()
		config.Producer.RequiredAcks = sarama.WaitForLocal
		config.Producer.Compression = sarama.CompressionSnappy
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		config.Producer.Return.Successes = true
		// config.Producer.Retry.Max = 10
		// config.Producer.Timeout = 5 * time.Second

		pc, err := sarama.NewSyncProducer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to create kafka producer:%s", err.Error())
		}

		producer = pc
	}

	return &eventService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		consumers:   consumers,
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
func (p *eventService) Initialize() error {
	// if p.producer != nil {
	// 	go func(pc sarama.AsyncProducer) {
	// 		for err := range pc.Errors() {
	// 			glog.Errorf("event service failed to send message: %s", err)
	// 		}
	// 		glog.Errorf("event service send stop")
	// 	}(p.producer)
	// }

	// subscribe topic from kmafka
	if len(p.consumers) > 0 {
		return p.subscribeKafkaTopic()
	}
	return nil
}

func (p *eventService) Name() string { return ServiceName }
func (p *eventService) bootstrap() error {
	return nil
}

// subscribeKafkaTopc subscribe topics from apiserver
func (p *eventService) subscribeKafkaTopic() error {
	for topic, consumer := range p.consumers {
		if partitionList, err := consumer.Partitions(topic); err == nil {
			for partition := range partitionList {
				pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
				if err != nil {
					return fmt.Errorf("event service subscribe kafka topic '%s' failed:%s", topic, err.Error())
				}
				p.waitgroup.Add(1)

				go func(p *eventService, pc sarama.PartitionConsumer) {
					defer p.waitgroup.Done()
					for msg := range pc.Messages() {
						glog.Infof("event service receive message: Key='%s', Value:%s", msg.Key, msg.Value)
						re := &RawEvent{}
						if err := json.Unmarshal(msg.Value, re); err == nil {
							p.handleKafkaEvent(re)
						}
					}
					glog.Errorf("event service receive message stop")
				}(p, pc)
			}
		}
	}
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
				p.publishKafkaMsg(p.nameOfEventBus(e), e)
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

func (p *eventService) publishKafkaMsg(topic string, e *Event) error {
	if p.producer != nil {
		// p.producer.Input() <- &sarama.ProducerMessage{
		// 	Topic: topic,
		// 	Key:   sarama.StringEncoder("sentel-broker-kafka"),
		// 	Value: sarama.ByteEncoder(value),
		// }
		re := &RawEvent{}
		re.Header, _ = json.Marshal(e.EventHeader)
		re.Payload, _ = json.Marshal(e.Detail)
		value, _ := json.Marshal(re)

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder("sentel-broker-kafka"),
			Value: sarama.ByteEncoder(value),
		}

		_, _, err := p.producer.SendMessage(msg)
		if err != nil {
			glog.Errorf("Failed to store your data:, %s", err)
		}
		// else {
		// 	// The tuple (topic, partition, offset) can be used as a unique identifier
		// 	// for a message in a Kafka cluster.
		// 	 glog.Infof("Your data is stored with unique identifier important/%d/%d", partition, offset)
		// }
		// glog.Infof("event service notify event '%s' to kafka", topic)
	}

	return nil
}

// Stop
func (p *eventService) Stop() {
	p.quitChan <- true
	for _, v := range p.consumers {
		v.Close()
	}

	if p.producer != nil {
		p.producer.Close()
	}
	p.waitgroup.Wait()
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

// handleEvent handle mqtt event from other service
func (p *eventService) handleKafkaEvent(re *RawEvent) {
	// cluster event manager only handle kafka event from other broker
	// Iterate all subscribers to notify
	e := &Event{}
	if err := json.Unmarshal(re.Header, &e.EventHeader); err != nil {
		glog.Errorf("event service unmarshal event common failed:%s", err.Error())
		return
	}

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
