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
	"strings"
	"sync"
	"time"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/iothub/cluster"
	sd "github.com/cloustone/sentel/iothub/service-discovery"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

type hubService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	clustermgr  cluster.ClusterManager
	tenants     map[string]*tenant
	mutex       sync.Mutex
	consumers   []sarama.Consumer
	hubdb       *hubDB
	recoverChan chan interface{}
	quitChan    chan interface{}
}

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

const SERVICE_NAME = "iothub"

type ServiceFactory struct{}

func (m ServiceFactory) New(c config.Config) (service.Service, error) {
	sarama.Logger = logger
	clustermgr, cerr := cluster.New(c)
	hubdb, nerr := newHubDB(c)
	discovery, derr := sd.New(c, sd.BackendZookeeper)
	if cerr != nil || nerr != nil || derr != nil {
		return nil, errors.New("service backend initialization failed")
	}
	clustermgr.SetServiceDiscovery(discovery)
	return &hubService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		consumers:   []sarama.Consumer{},
		clustermgr:  clustermgr,
		tenants:     make(map[string]*tenant),
		mutex:       sync.Mutex{},
		hubdb:       hubdb,
		recoverChan: make(chan interface{}),
		quitChan:    make(chan interface{}),
	}, nil
}

// Name
func (p *hubService) Name() string { return SERVICE_NAME }

// Initialize load iothub data and recovery from scratch
func (p *hubService) Initialize() error {
	tenants := p.hubdb.getAllTenants()
	if len(tenants) > 0 {
		for _, t := range tenants {
			p.tenants[t.TenantId] = &t
		}
		p.recoverChan <- true
	}
	return nil
}

// Start
func (p *hubService) Start() error {
	tc, subError1 := p.subscribeTopic(message.TopicNameTenant)
	pc, subError2 := p.subscribeTopic(message.TopicNameProduct)
	if subError1 != nil || subError2 != nil {
		return errors.New("iothub failed to subsribe tenant and product topic from kafka")
	}
	p.consumers = append(p.consumers, tc)
	p.consumers = append(p.consumers, pc)
	p.waitgroup.Add(1)
	go func(p *hubService) {
		defer p.waitgroup.Done()
		for {
			select {
			case <-p.quitChan:
			case <-p.recoverChan:
				p.recoverStartup()
			}
		}
	}(p)
	return nil
}

// Stop
func (p *hubService) Stop() {
	p.quitChan <- true
	for _, consumer := range p.consumers {
		consumer.Close()
	}
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.recoverChan)
}

// recoverStartup load hub data and confirm service status
func (p *hubService) recoverStartup() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	retries := []*tenant{}
	for tid, t := range p.tenants {
		if t.ServiceState != cluster.ServiceStateNone {
			if _, err := p.clustermgr.IntrospectService(tid); err != nil {
				if _, err := p.clustermgr.CreateService(tid, t.InstanceReplicas); err != nil {
					retries = append(retries, t)
				}
			}
		}
	}
	// retry to recover again
	for _, t := range retries {
		if _, err := p.clustermgr.CreateService(t.TenantId, t.InstanceReplicas); err != nil {
			glog.Error("service '%s' recovery failed", t.TenantId)
		}
	}
}

// subscribeTopc subscribe topics from apiserver
func (p *hubService) subscribeTopic(topic string) (sarama.Consumer, error) {
	endpoint := p.config.MustString("iothub", "kafka")
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
			p.waitgroup.Add(1)

			go func(p *hubService, pc sarama.PartitionConsumer) {
				defer p.waitgroup.Done()
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
func (p *hubService) handleNotification(topic string, value []byte) error {
	glog.Infof("iothub receive message from topic '%s'", topic)
	var err error
	switch topic {
	case message.TopicNameTenant:
		err = p.handleTenantNotify(value)
	case message.TopicNameProduct:
		err = p.handleProductNotify(value)
	default:
	}
	return err
}

// handleProductNotify handle notification about product from api server
func (p *hubService) handleProductNotify(value []byte) error {
	tf := message.ProductTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}
	switch tf.Action {
	case message.ObjectActionRegister:
		_, err := p.createProduct(tf.TenantId, tf.ProductId, tf.Replicas)
		return err
	case message.ObjectActionDelete:
		return p.removeProduct(tf.TenantId, tf.ProductId)
	}
	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *hubService) handleTenantNotify(value []byte) error {
	tf := message.TenantTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}

	switch tf.Action {
	case message.ObjectActionRegister:
		p.createTenant(tf.TenantId)
	case message.ObjectActionDelete:
		return p.removeTenant(tf.TenantId)
	}
	return nil
}

// addTenant add tenant to iothub
func (p *hubService) createTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		p.tenants[tid] = &tenant{
			TenantId:     tid,
			CreatedAt:    time.Now(),
			Products:     make(map[string]*product),
			ServiceState: cluster.ServiceStateNone,
		}
		p.hubdb.createTenant(p.tenants[tid])
		return nil
	}
	return fmt.Errorf("tenant '%s' already existed in iothub", tid)
}

// deleteTenant remove tenant from iothub
func (p *hubService) removeTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		return fmt.Errorf("tenant '%s' doesn't exist in iothub", tid)
	}
	t := p.tenants[tid]
	// remove service created early
	if t.ServiceState != cluster.ServiceStateNone {
		p.clustermgr.RemoveService(t.ServiceId)
	}
	// Delete all products
	for name, _ := range t.Products {
		if err := p.removeProduct(tid, name); err != nil {
			glog.Errorf("iothub remove tenant '%s' product '%s' failed", tid, name)
		}
	}
	// Remove service
	delete(p.tenants, tid)
	p.hubdb.removeTenant(tid)
	return nil
}

func (p *hubService) isProductExist(tid, pid string) bool {
	if t, found := p.tenants[tid]; found {
		if _, found := t.Products[pid]; found {
			return true
		}
	}
	return false
}

// addProduct add product to iothub
func (p *hubService) createProduct(tid, pid string, replicas int32) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.isProductExist(tid, pid) {
		return "", fmt.Errorf("product '%s' of '%s' already exist in iothub", pid, tid)
	}
	t := p.tenants[tid]
	if t.ServiceState == cluster.ServiceStateNone {
		serviceId, err := p.clustermgr.CreateService(tid, replicas)
		if err != nil {
			return "", err
		}
		t.ServiceState = cluster.ServiceStateStarted
		t.ServiceId = serviceId
		t.ServiceName = tid
	}
	product := &product{ProductId: pid, CreatedAt: time.Now()}
	t.Products[pid] = product
	p.hubdb.createProduct(tid, pid)
	return t.ServiceId, nil
}

// deleteProduct delete product from iothub
func (p *hubService) removeProduct(tid string, pid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isProductExist(tid, pid) {
		return fmt.Errorf("product '%s' of '%s' does not exist in iothub", pid, tid)
	}
	t := p.tenants[tid]
	delete(t.Products, pid)
	if len(t.Products) == 0 {
		p.clustermgr.RemoveService(t.ServiceId)
		t.ServiceState = cluster.ServiceStateNone
	}
	p.hubdb.removeProduct(tid, pid)
	return nil
}
