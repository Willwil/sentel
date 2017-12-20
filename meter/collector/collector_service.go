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
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/common"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type CollectorService struct {
	com.ServiceBase
	consumer sarama.Consumer
}

// CollectorServiceFactory
type CollectorServiceFactory struct{}

// New create apiService service factory
func (m CollectorServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
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

	return &CollectorService{
		ServiceBase: com.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumer: consumer,
	}, nil

}

// Name
func (p *CollectorService) Name() string {
	return "collector"
}

// Start
func (p *CollectorService) Start() error {
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
func (p *CollectorService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	p.consumer.Close()
}

// handleNotifications handle notification from kafka
func (p *CollectorService) handleNotifications(topic string, value []byte) error {
	if err := handleTopicObject(p, context.Background(), topic, value); err != nil {
		glog.Error(err)
		return err
	}
	return nil
}

func (p *CollectorService) getDatabase() (*mgo.Database, error) {
	// check mongo db configuration
	hosts := p.Config.MustString("meter", "mongo")
	timeout := p.Config.MustInt("meter", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("iothub"), nil
}

// subscribeTopc subscribe topics from apiserver
func (p *CollectorService) subscribeTopic(topic string) error {
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
		p.WaitGroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.WaitGroup.Done()
			for msg := range pc.Messages() {
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}
