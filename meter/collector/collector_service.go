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
	"context"
	"sync"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type collectorService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	consumer  message.Consumer
}

const SERVICE_NAME = "meter"

// collectorServiceFactory
type ServiceFactory struct{}

// New create apiService service factory
func (m ServiceFactory) New(c config.Config) (service.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("meter", "mongo")
	timeout := c.MustInt("meter", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	// kafka
	khosts := c.MustString("meter", "kafka")
	consumer, _ := message.NewConsumer(khosts, "meter")
	return &collectorService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		consumer:  consumer,
	}, nil

}

// Name
func (p *collectorService) Name() string      { return SERVICE_NAME }
func (p *collectorService) Initialize() error { return nil }

// Start
func (p *collectorService) Start() error {
	p.consumer.Subscribe(TopicNameNode, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameClient, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameSession, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameSubscription, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNamePublish, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameMetric, p.messageHandlerFunc, nil)
	p.consumer.Subscribe(TopicNameStats, p.messageHandlerFunc, nil)
	return nil
}

// Stop
func (p *collectorService) Stop() {
	p.consumer.Close()
	p.waitgroup.Wait()

}

// handleNotifications handle notification from kafka
func (p *collectorService) messageHandlerFunc(topic string, value []byte, ctx interface{}) {
	if err := handleTopicObject(p, context.Background(), topic, value); err != nil {
		glog.Error(err)
	}
}

func (p *collectorService) getDatabase() (*mgo.Database, error) {
	// check mongo db configuration
	hosts := p.config.MustString("meter", "mongo")
	timeout := p.config.MustInt("meter", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("iothub"), nil
}
