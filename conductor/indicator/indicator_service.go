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

package indicator

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/cloustone/sentel/conductor/executor"
	"github.com/cloustone/sentel/core"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
	"gopkg.in/mgo.v2"
)

type IndicatorService struct {
	core.ServiceBase
	consumer   sarama.Consumer
	mongoHosts string // mongo hosts
}

// IndicatorServiceFactory
type IndicatorServiceFactory struct{}

// New create apiService service factory
func (this *IndicatorServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("conductor", "mongo")
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// kafka
	khosts := c.MustString("conductor", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &IndicatorService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumer:   consumer,
		mongoHosts: hosts,
	}, nil

}

// Name
func (this *IndicatorService) Name() string {
	return "indicator"
}

// Start
func (this *IndicatorService) Start() error {
	this.subscribeTopic("rule")
	this.WaitGroup.Wait()
	return nil
}

// Stop
func (this *IndicatorService) Stop() {
	this.consumer.Close()
}

// subscribeTopc subscribe topics from apiserver
func (this *IndicatorService) subscribeTopic(topic string) error {
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
		this.WaitGroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer this.WaitGroup.Done()
			for msg := range pc.Messages() {
				this.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}

type ruleTopic struct {
	RuleName  string `json:"ruleName"`
	RuleId    string `json:"ruleId"`
	ProductId string `json:"productId"`
	Action    string `json:"action"`
}

// handleNotifications handle notification from kafka
func (this *IndicatorService) handleNotifications(topic string, value []byte) error {
	rule := ruleTopic{}
	if err := json.Unmarshal(value, &topic); err != nil {
		return err
	}
	r := &executor.Rule{
		RuleName:  rule.RuleName,
		RuleId:    rule.RuleId,
		ProductId: rule.ProductId,
		Action:    rule.Action,
	}
	return executor.HandleRuleNotification(r)
}
