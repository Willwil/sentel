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

package iothub

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	apiserver "github.com/cloustone/sentel/apiserver/util"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

const (
	notifyActionCreate = "create"
	notifyActionDelete = "delete"
	notifyActionUpdate = "update"
)

// TenantNofiy is notification object from api server by kafka
type tenantNotify struct {
	action      string `json:"action"`
	id          string `json:"tenantid"`
	brokerCount string `json:"brokerCount"`
}

type NotifyService struct {
	core.ServiceBase
	consumer sarama.Consumer
}

// NotifyServiceFactory
type NotifyServiceFactory struct{}

// New create apiService service factory
func (m *NotifyServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// kafka
	khosts, err := core.GetServiceEndpoint(c, "iothub", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	return &NotifyService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumer: consumer,
	}, nil
}

// Name
func (p *NotifyService) Name() string {
	return "notify"
}

// Start
func (p *NotifyService) Start() error {
	if err := p.subscribeTopic(apiserver.TopicNameTenant); err != nil {
		return err
	}
	if err := p.subscribeTopic(apiserver.TopicNameProduct); err != nil {
		return err
	}
	p.WaitGroup.Wait()
	return nil
}

// Stop
func (p *NotifyService) Stop() {
	p.consumer.Close()
	p.WaitGroup.Wait()
}

// subscribeTopc subscribe topics from apiserver
func (p *NotifyService) subscribeTopic(topic string) error {
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

// handleNotifications handle notification from kafka
func (p *NotifyService) handleNotifications(topic string, value []byte) error {
	switch topic {
	case apiserver.TopicNameTenant:
		obj := &tenantNotify{}
		if err := json.Unmarshal(value, obj); err != nil {
			return err
		}
		return p.handleTenantNotify(obj)
	}

	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *NotifyService) handleTenantNotify(tf *tenantNotify) error {
	glog.Infof("iothub-notifyservice: tenant(%s) notification received", tf.id)

	hub := getIothub()

	switch tf.action {
	case notifyActionCreate:
		return hub.addTenant(tf.id)
	case notifyActionDelete:
		return hub.deleteTenant(tf.id)
	case notifyActionUpdate:
		return hub.updateTenant(tf.id)
	}
	return nil
}
