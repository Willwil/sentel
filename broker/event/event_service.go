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
	ServiceName         = "event"
	fmtOfMqttEventBus   = "broker-%s-%s-event-mqtt"
	fmtOfBrokerEventBus = "broker-%s-%s-event-broker"
)

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

type eventService struct {
	base.ServiceBase
	brokerId    string
	consumer    map[string]sarama.Consumer // kafka client endpoint
	producer    sarama.AsyncProducer
	eventChan   chan *Event                    // Event notify channel
	subscribers map[uint32][]subscriberContext // All subscriber's contex
	mutex       sync.Mutex                     // Mutex to lock subscribers
	tenant      string
	product     string
}

// subscriberContext hold subscribre's handler and context
type subscriberContext struct {
	handler EventHandler
	ctx     interface{}
}

// New create global broker
func New(c core.Config, quit chan os.Signal) (base.Service, error) {
	// Retrieve tenant and product
	tenant := c.MustString("broker", "tenant")
	product := c.MustString("broker", "product")

	var consumer = make(map[string]sarama.Consumer)
	var producer sarama.AsyncProducer
	sarama.Logger = logger
	if khosts, err := core.GetServiceEndpoint(c, "broker", "kafka"); err == nil && khosts != "" {
		config := sarama.NewConfig()
		config.ClientID = base.GetBrokerId() + "_MQTT"
		// config.Consumer.MaxWaitTime = time.Duration(5 * time.Second)
		// config.Consumer.Offsets.CommitInterval = 1 * time.Second
		// config.Consumer.Offsets.Initial = sarama.OffsetNewest
		cm, err := sarama.NewConsumer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to connect with kafka server(MQTT) '%s'", khosts)
		}
		consumer[fmt.Sprintf(fmtOfMqttEventBus, tenant, product)] = cm

		config = sarama.NewConfig()
		config.ClientID = base.GetBrokerId() + "_Broker"
		cb, err := sarama.NewConsumer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to connect with kafka server(Broker) '%s'", khosts)
		}
		consumer[fmt.Sprintf(fmtOfBrokerEventBus, tenant, product)] = cb
		glog.Infof("event service connected with kafka '%s' with broker id '%s'", khosts, base.GetBrokerId())

		config = sarama.NewConfig()
		config.ClientID = base.GetBrokerId()
		config.Producer.RequiredAcks = sarama.WaitForLocal
		config.Producer.Compression = sarama.CompressionSnappy
		config.Producer.Partitioner = sarama.NewRandomPartitioner
		config.Producer.Return.Successes = true
		// config.Producer.Retry.Max = 10
		// config.Producer.Timeout = 5 * time.Second

		pc, err := sarama.NewAsyncProducer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to create kafka producer:%s", err.Error())
		}

		go func(pc sarama.AsyncProducer) {
			for err := range pc.Errors() {
				glog.Errorf("event service failed to send message: %s", err)
			}
		}(pc)

		producer = pc
	}

	return &eventService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumer:    consumer,
		producer:    producer,
		eventChan:   make(chan *Event),
		subscribers: make(map[uint32][]subscriberContext),
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
	// subscribe topic from kafka
	if p.consumer != nil {
		// mqttEventBus := fmt.Sprintf(fmtOfMqttEventBus, p.tenant, p.product)
		// brokerEventBus := fmt.Sprintf(fmtOfBrokerEventBus, p.tenant, p.product)
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
	for topic, consumer := range p.consumer {
		if partitionList, err := consumer.Partitions(topic); err == nil {
			for partition := range partitionList {
				pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetOldest)
				if err != nil {
					return fmt.Errorf("event service subscribe kafka topic '%s' failed:%s", topic, err.Error())
				}
				p.WaitGroup.Add(1)

				go func(p *eventService, pc sarama.PartitionConsumer) {
					defer p.WaitGroup.Done()
					for msg := range pc.Messages() {
						// glog.Infof("event service receive message: Key='%s'", msg.Key)
						obj := &Event{}
						if err := json.Unmarshal(msg.Value, obj); err == nil {
							p.handleKafkaEvent(obj)
						}
					}
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
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

func (p *eventService) publishKafkaMsg(topic string, value sarama.Encoder) error {
	if p.producer != nil {
		p.producer.Input() <- &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder("sentel-broker-kafka"),
			Value: value,
		}
		glog.Infof("event service notify event '%s' to kafka", topic)
	}
	return nil
}

// Stop
func (p *eventService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	if p.consumer != nil {
		for _, v := range p.consumer {
			v.Close()
		}
	}

	if p.producer != nil {
		p.producer.Close()
	}
	p.WaitGroup.Wait()
	close(p.eventChan)
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
