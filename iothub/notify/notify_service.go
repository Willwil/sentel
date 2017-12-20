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

package notify

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/iothub/hub"
	"github.com/golang/glog"
)

type notifyService struct {
	com.ServiceBase
	consumers []sarama.Consumer
}

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

type ServiceFactory struct{}

// New create apiService service factory
func (m ServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
	sarama.Logger = logger
	return &notifyService{
		ServiceBase: com.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumers: []sarama.Consumer{},
	}, nil
}

// Name
func (p *notifyService) Name() string { return "notify" }

// Start
func (p *notifyService) Start() error {
	tc, subError1 := p.subscribeTopic(com.TopicNameTenant)
	pc, subError2 := p.subscribeTopic(com.TopicNameProduct)
	if subError1 != nil || subError2 != nil {
		return errors.New("iothub failed to subsribe tenant and product topic from kafka")
	}
	p.consumers = append(p.consumers, tc)
	p.consumers = append(p.consumers, pc)
	return nil
}

// Stop
func (p *notifyService) Stop() {
	for _, consumer := range p.consumers {
		consumer.Close()
	}
	p.WaitGroup.Wait()
	close(p.Quit)
}

// subscribeTopc subscribe topics from apiserver
func (p *notifyService) subscribeTopic(topic string) (sarama.Consumer, error) {
	endpoint := p.Config.MustString("iothub", "kafka")
	glog.Infof("iothub get kafka service endpoint: %s", endpoint)

	config := sarama.NewConfig()
	config.ClientID = "iothub-notify-" + topic
	consumer, err := sarama.NewConsumer(strings.Split(endpoint, ","), config)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", endpoint)
	}

	partitionList, err := consumer.Partitions(topic)
	if err != nil {
		return nil, err
	}
	for partition := range partitionList {
		if pc, err := consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest); err != nil {
			return nil, fmt.Errorf("iothub subscribe kafka topic '%s' failed:%s", topic, err.Error())
		} else {
			p.WaitGroup.Add(1)

			go func(p *notifyService, pc sarama.PartitionConsumer) {
				defer p.WaitGroup.Done()
				for msg := range pc.Messages() {
					if err := p.handleNotification(topic, msg.Value); err != nil {
						glog.Errorf("iothub failed to handle topic '%s': '%s'", topic, err.Error())
					}
				}
				glog.Info("iothub message receiver stop")
			}(p, pc)
		}
	}
	return consumer, nil
}

// handleNotify handle notification from kafka
func (p *notifyService) handleNotification(topic string, value []byte) error {
	glog.Infof("iothub receive message from topic '%s'", topic)
	var err error
	switch topic {
	case com.TopicNameTenant:
		err = p.handleTenantNotify(value)
	case com.TopicNameProduct:
		err = p.handleProductNotify(value)
	default:
	}
	return err
}

// handleProductNotify handle notification about product from api server
func (p *notifyService) handleProductNotify(value []byte) error {
	hub := hub.GetIothub()
	tf := com.ProductTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}
	switch tf.Action {
	case com.ObjectActionRegister:
		_, err := hub.CreateProduct(tf.TenantId, tf.ProductId, tf.Replicas)
		return err
	case com.ObjectActionDelete:
		return hub.RemoveProduct(tf.TenantId, tf.ProductId)
	}
	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *notifyService) handleTenantNotify(value []byte) error {
	hub := hub.GetIothub()
	tf := com.TenantTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}

	switch tf.Action {
	case com.ObjectActionRegister:
		hub.CreateTenant(tf.TenantId)
	case com.ObjectActionDelete:
		return hub.RemoveTenant(tf.TenantId)
	}
	return nil
}
