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
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

type NotifyService struct {
	core.ServiceBase
	consumers []sarama.Consumer
}

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

// NotifyServiceFactory
type NotifyServiceFactory struct{}

// New create apiService service factory
func (m *NotifyServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	sarama.Logger = logger
	return &NotifyService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		consumers: []sarama.Consumer{},
	}, nil
}

// Name
func (p *NotifyService) Name() string {
	return "notify"
}

// Start
func (p *NotifyService) Start() error {
	tc, subError1 := p.subscribeTopic(core.TopicNameTenant)
	pc, subError2 := p.subscribeTopic(core.TopicNameProduct)
	if subError1 != nil || subError2 != nil {
		return errors.New("iothub failed to subsribe tenant and product topic from kafka")
	}
	p.consumers = append(p.consumers, tc)
	p.consumers = append(p.consumers, pc)
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
	for _, consumer := range p.consumers {
		consumer.Close()
	}
	p.WaitGroup.Wait()
	close(p.Quit)
}

// subscribeTopc subscribe topics from apiserver
func (p *NotifyService) subscribeTopic(topic string) (sarama.Consumer, error) {
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

			go func(p *NotifyService, pc sarama.PartitionConsumer) {
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
func (p *NotifyService) handleNotification(topic string, value []byte) error {
	glog.Infof("iothub receive message from topic '%s'", topic)
	var err error
	switch topic {
	case core.TopicNameTenant:
		err = p.handleTenantNotify(value)
	case core.TopicNameProduct:
		err = p.handleProductNotify(value)
	default:
	}
	return err
}

// handleProductNotify handle notification about product from api server
func (p *NotifyService) handleProductNotify(value []byte) error {
	hub := getIothub()
	tf := core.ProductTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}
	switch tf.Action {
	case core.ObjectActionRegister:
		_, err := hub.createProduct(tf.TenantId, tf.ProductId, tf.Replicas)
		return err
	case core.ObjectActionDelete:
		return hub.removeProduct(tf.TenantId, tf.ProductId)
	}
	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *NotifyService) handleTenantNotify(value []byte) error {
	hub := getIothub()
	tf := core.TenantTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}

	switch tf.Action {
	case core.ObjectActionRegister:
		hub.createTenant(tf.TenantId)
	case core.ObjectActionDelete:
		return hub.removeTenant(tf.TenantId)
	}
	return nil
}
