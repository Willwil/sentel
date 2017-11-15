//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
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
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type CollectorService struct {
	config     core.Config
	quit       chan os.Signal
	wg         sync.WaitGroup
	consumer   sarama.Consumer
	mongoHosts string // mongo hosts
}

// CollectorServiceFactory
type CollectorServiceFactory struct{}

// New create apiService service factory
func (m *CollectorServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, err := c.String("ceilometer", "mongo")
	if err != nil || hosts == "" {
		return nil, errors.New("Invalid mongo configuration")
	}
	// try connect with mongo db
	session, err := mgo.Dial(hosts)
	if err != nil {
		return nil, err
	}
	session.Close()

	// kafka
	khosts := c.MustString("collector", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", hosts)
	}

	return &CollectorService{
		config:     c,
		wg:         sync.WaitGroup{},
		quit:       quit,
		consumer:   consumer,
		mongoHosts: hosts,
	}, nil

}

// Name
func (this *CollectorService) Name() string {
	return "collector"
}

// Start
func (this *CollectorService) Start() error {
	this.subscribeTopic(TopicNameNode)
	this.subscribeTopic(TopicNameClient)
	this.subscribeTopic(TopicNameSession)
	this.subscribeTopic(TopicNameSubscription)
	this.subscribeTopic(TopicNamePublish)
	this.subscribeTopic(TopicNameMetric)
	this.subscribeTopic(TopicNameStats)
	this.wg.Wait()
	return nil
}

// Stop
func (this *CollectorService) Stop() {
	this.consumer.Close()
}

// handleNotifications handle notification from kafka
func (this *CollectorService) handleNotifications(topic string, value []byte) error {
	if err := handleTopicObject(this, context.Background(), topic, value); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func (this *CollectorService) getDatabase() (*mgo.Database, error) {
	session, err := mgo.DialWithTimeout(this.mongoHosts, 2*time.Second)
	if err != nil {
		glog.Fatalf("Failed to connect with mongo:%s", this.mongoHosts)
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("iothub"), nil
}

// subscribeTopc subscribe topics from apiserver
func (this *CollectorService) subscribeTopic(topic string) error {
	partitionList, err := this.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := this.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		this.wg.Add(1)

		go func(sarama.PartitionConsumer) {
			defer this.wg.Done()
			for msg := range pc.Messages() {
				this.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}
