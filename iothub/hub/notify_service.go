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

package hub

import (
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	apiserver "github.com/cloustone/sentel/apiserver/util"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

const (
	notifyActionCreate = "create"
	notifyActionDelete = "delete"
	notifyActionUpdate = "update"
	notifyActionStart  = "start"
	notifyActionStop   = "stop"
)

// TenantNofiy is notification object from api server by kafka
type productNotify struct {
	Action         string `json:"action"`
	TenantId       string `json:"tenantId"`
	ProductId      string `json:"productId"`
	BrokerReplicas int32  `json:"brokerReplicas"`
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
	if err := p.subscribeTopic(apiserver.TopicNameProduct); err != nil {
		return err
	}
	go func(p *NotifyService) {
		for {
			select {
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *NotifyService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.consumer.Close()
	p.WaitGroup.Wait()
	close(p.Quit)
}

// subscribeTopc subscribe topics from apiserver
func (p *NotifyService) subscribeTopic(topic string) error {
	if partitionList, err := p.consumer.Partitions(topic); err == nil {
		for partition := range partitionList {
			pc, err := p.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
			if err != nil {
				return fmt.Errorf("event service subscribe kafka topic '%s' failed:%s", topic, err.Error())
			}
			p.WaitGroup.Add(1)

			go func(p *NotifyService, pc sarama.PartitionConsumer) {
				defer p.WaitGroup.Done()
				for msg := range pc.Messages() {
					glog.Infof("event service receive message: Key='%s', Value:%s", msg.Key, msg.Value)
					obj := &productNotify{}
					if err := json.Unmarshal(msg.Value, obj); err == nil {
						p.handleProductNotification(obj)
					}
				}
				glog.Errorf("event service receive message stop")
			}(p, pc)
		}
	}
	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *NotifyService) handleProductNotification(tf *productNotify) error {
	glog.Infof("iothub-notifyservice: tenant(%s) notification received", tf.TenantId)

	hub := getIothub()

	switch tf.Action {
	case notifyActionCreate:
		_, err := hub.addProduct(tf.TenantId, tf.ProductId, tf.BrokerReplicas)
		return err
	case notifyActionDelete:
		return hub.deleteProduct(tf.TenantId, tf.ProductId)
	case notifyActionStart:
		return hub.startProduct(tf.TenantId, tf.ProductId)
	case notifyActionStop:
		return hub.stopProduct(tf.TenantId, tf.ProductId)
	}
	return nil
}
