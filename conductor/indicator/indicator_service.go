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

package indicator

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cloustone/sentel/conductor/executor"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"

	"github.com/Shopify/sarama"
	"gopkg.in/mgo.v2"
)

type indicatorService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	consumer  sarama.Consumer
}

// indicatorServiceFactory
type ServiceFactory struct{}

// New create apiService service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("conductor", "mongo")
	timeout := c.MustInt("conductor", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	return &indicatorService{
		config:    c,
		waitgroup: sync.WaitGroup{},
	}, nil

}

// Name
func (p *indicatorService) Name() string      { return "indicator" }
func (p *indicatorService) Initialize() error { return nil }

// Start
func (p *indicatorService) Start() error {
	consumer, err := p.subscribeTopic(message.TopicNameRule)
	if err != nil {
		return fmt.Errorf("conductor failed to subscribe kafka event : %s", err.Error())
	}
	p.consumer = consumer
	return nil
}

// Stop
func (p *indicatorService) Stop() {
	if p.consumer != nil {
		p.consumer.Close()
	}
	p.waitgroup.Wait()
}

// subscribeTopc subscribe topics from apiserver
func (p *indicatorService) subscribeTopic(topic string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.ClientID = "sentel_conductor_indicator"
	khosts, _ := p.config.String("conductor", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}
	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		consumer.Close()
		return nil, fmt.Errorf("Failed to get list of partions:%s", err.Error())
	}

	for partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			consumer.Close()
			return nil, fmt.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
		}
		p.waitgroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.waitgroup.Done()
			for msg := range pc.Messages() {
				glog.Infof("donductor recevied topic '%s': '%s'", string(msg.Topic), msg.Value)
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return consumer, nil
}

// handleNotifications handle notification from kafka
func (p *indicatorService) handleNotifications(topic string, value []byte) error {
	r := message.RuleTopic{}
	if err := json.Unmarshal(value, &r); err != nil {
		glog.Errorf("conductor failed to resolve topic from kafka: '%s'", err)
		return err
	}
	return executor.HandleRuleNotification(&r)
}
