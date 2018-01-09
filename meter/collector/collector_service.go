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
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type collectorService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	consumer  sarama.Consumer
}

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
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", hosts)
	}

	return &collectorService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		consumer:  consumer,
	}, nil

}

// Name
func (p *collectorService) Name() string {
	return "collector"
}

func (p *collectorService) Initialize() error { return nil }

// Start
func (p *collectorService) Start() error {
	p.subscribeTopic(TopicNameNode)
	p.subscribeTopic(TopicNameClient)
	p.subscribeTopic(TopicNameSession)
	p.subscribeTopic(TopicNameSubscription)
	p.subscribeTopic(TopicNamePublish)
	p.subscribeTopic(TopicNameMetric)
	p.subscribeTopic(TopicNameStats)
	return nil
}

// Stop
func (p *collectorService) Stop() {
	p.consumer.Close()
	p.waitgroup.Wait()

}

// handleNotifications handle notification from kafka
func (p *collectorService) handleNotifications(topic string, value []byte) error {
	if err := handleTopicObject(p, context.Background(), topic, value); err != nil {
		glog.Error(err)
		return err
	}
	return nil
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

// subscribeTopc subscribe topics from apiserver
func (p *collectorService) subscribeTopic(topic string) error {
	partitionList, err := p.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
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
