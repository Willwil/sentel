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

package collector

import (
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
)

type collectorService struct {
	config   config.Config
	consumer message.Consumer
}

const SERVICE_NAME = "collector"

// collectorServiceFactory
type ServiceFactory struct{}

func (m ServiceFactory) New(c config.Config) (service.Service, error) {
	consumer, _ := message.NewConsumer(c, "iotmanager")
	return &collectorService{
		config:   c,
		consumer: consumer,
	}, nil

}

// Name
func (p *collectorService) Name() string      { return SERVICE_NAME }
func (p *collectorService) Initialize() error { return nil }

// MessageFactory interface implementation
func (p *collectorService) CreateMessage(topic string) message.Message {
	switch topic {
	case TopicNameNode:
		return &NodeTopic{TopicName: topic}
	case TopicNameClient:
		return &ClientTopic{TopicName: topic}
	case TopicNameSession:
		return &SessionTopic{TopicName: topic}
	case TopicNameSubscription:
		return &SubscriptionTopic{TopicName: topic}
	case TopicNamePublish:
		return &PublishTopic{TopicName: topic}
	case TopicNameMetric:
		return &MetricTopic{TopicName: topic}
	case TopicNameStats:
		return &StatsTopic{TopicName: topic}
	default:
		return nil
	}
}

// Start
func (p *collectorService) Start() error {
	p.consumer.Subscribe(TopicNameNode, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameClient, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameSession, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameSubscription, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNamePublish, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameMetric, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameStats, p.messageHandlerFunc, nil)
	p.consumer.SetMessageFactory(p)
	p.consumer.Start()
	return nil
}

// Stop
func (p *collectorService) Stop() {
	p.consumer.Close()
}

// handleNotifications handle notification from kafka
func (p *collectorService) messageHandlerFunc(msg message.Message, ctx interface{}) {
	if handler, ok := msg.(topicHandler); ok {
		handler.handleTopic(p.config, nil)
	}
}
