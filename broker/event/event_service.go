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
	producer    sarama.SyncProducer
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
	var producer sarama.SyncProducer
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

		pc, err := sarama.NewSyncProducer(strings.Split(khosts, ","), config)
		if err != nil {
			return nil, fmt.Errorf("event service failed to create kafka producer:%s", err.Error())
		}

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
	switch e.Common.Type {
	case SessionCreate, SessionDestroy, TopicPublish, TopicSubscribe, TopicUnsubscribe:
		return fmt.Sprintf(fmtOfMqttEventBus, p.tenant, p.product)
	default:
		return fmt.Sprintf(fmtOfBrokerEventBus, p.tenant, p.product)
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
	if p.consumer != nil {
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
				pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
				if err != nil {
					return fmt.Errorf("event service subscribe kafka topic '%s' failed:%s", topic, err.Error())
				}
				p.WaitGroup.Add(1)

				go func(p *eventService, pc sarama.PartitionConsumer) {
					defer p.WaitGroup.Done()
					for msg := range pc.Messages() {
						glog.Infof("event service receive message: Key='%s'", msg.Key)
						obj := &serEvent{}
						if err := json.Unmarshal(msg.Value, obj); err == nil {
							p.handleKafkaEvent(obj)
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
				if _, found := p.subscribers[e.Common.Type]; found {
					subscribers := p.subscribers[e.Common.Type]
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

func (p *eventService) publishKafkaMsg(topic string, e *Event) error {
	if p.producer != nil {
		// p.producer.Input() <- &sarama.ProducerMessage{
		// 	Topic: topic,
		// 	Key:   sarama.StringEncoder("sentel-broker-kafka"),
		// 	Value: sarama.ByteEncoder(value),
		// }
		se := &serEvent{}
		se.Common, _ = json.Marshal(e.Common)
		se.Detail, _ = json.Marshal(e.Detail)
		value, _ := json.Marshal(se)

		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder("sentel-broker-kafka"),
			Value: sarama.ByteEncoder(value),
		}

		partition, offset, err := p.producer.SendMessage(msg)
		if err != nil {
			glog.Errorf("Failed to store your data:, %s", err)
		} else {
			// The tuple (topic, partition, offset) can be used as a unique identifier
			// for a message in a Kafka cluster.
			glog.Errorf("Your data is stored with unique identifier important/%d/%d", partition, offset)
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

func (p *eventService) unmarshalDetail(e *Event, detail json.RawMessage) (interface{}, error) {
	switch e.Common.Type {
	case SessionCreate:
		sc := &SessionCreateDetail{}
		if err := json.Unmarshal(detail, sc); err != nil {
			return nil, fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		return sc, nil
	case SessionDestroy:
		sd := &SessionDestroyDetail{}
		if err := json.Unmarshal(detail, sd); err != nil {
			return nil, fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		return sd, nil
	case TopicSubscribe:
		ts := &TopicSubscribeDetail{}
		if err := json.Unmarshal(detail, ts); err != nil {
			return nil, fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		return ts, nil
	case TopicUnsubscribe:
		tu := &TopicUnsubscribeDetail{}
		if err := json.Unmarshal(detail, tu); err != nil {
			return nil, fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		return tu, nil
	case TopicPublish:
		tp := &TopicPublishDetail{}
		if err := json.Unmarshal(detail, tp); err != nil {
			return nil, fmt.Errorf("event service unmarshal event detail failed:%s", err.Error())
		}
		return tp, nil
	}

	return nil, fmt.Errorf("event service unmarshal event detail invalid type:%d", e.Common.Type)
}

// handleEvent handle mqtt event from other service
func (p *eventService) handleKafkaEvent(se *serEvent) {
	// cluster event manager only handle kafka event from other broker
	// Iterate all subscribers to notify
	e := &Event{}
	if err := json.Unmarshal(se.Common, e.Common); err != nil {
		return
	}

	if e.Common.BrokerId != base.GetBrokerId() {
		if _, found := p.subscribers[e.Common.Type]; found {
			subscribers := p.subscribers[e.Common.Type]
			var detail interface{}
			detail, err := p.unmarshalDetail(e, se.Detail)
			if err != nil {
				glog.Errorf(err.Error())
				return
			}
			e.Detail = detail
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
