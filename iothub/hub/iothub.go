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

	mgo "gopkg.in/mgo.v2"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/iothub/cluster"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

type iothubService struct {
	config     config.Config
	waitgroup  sync.WaitGroup
	clustermgr cluster.ClusterManager
	tenants    map[string]*tenant
	mutex      sync.Mutex
	consumers  []sarama.Consumer
}
type tenant struct {
	tid       string
	createdAt time.Time
	products  map[string]*product
	networkId string
}

type product struct {
	pid         string
	tid         string
	createdAt   time.Time
	serviceName string
}

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

const SERVICE_NAME = "iothub"

type ServiceFactory struct{}

func (m ServiceFactory) New(c config.Config) (service.Service, error) {
	sarama.Logger = logger
	clustermgr, err := cluster.New(c)
	if err != nil {
		return nil, err
	}
	// try connect with mongo db
	addr := c.MustString("iothub", "mongo")
	session, err := mgo.DialWithTimeout(addr, 1*time.Second)
	if err != nil {
		return nil, fmt.Errorf("iothub connect with mongo '%s'failed: '%s'", addr, err.Error())
	}
	session.Close()

	return &iothubService{
		config:     c,
		waitgroup:  sync.WaitGroup{},
		consumers:  []sarama.Consumer{},
		clustermgr: clustermgr,
		tenants:    make(map[string]*tenant),
		mutex:      sync.Mutex{},
	}, nil
}

// Name
func (p *iothubService) Name() string      { return SERVICE_NAME }
func (p *iothubService) Initialize() error { return nil }

// GetiothubService return iothub service instance
func GetIothub() *iothubService {
	mgr := service.GetServiceManager()
	return mgr.GetService(SERVICE_NAME).(*iothubService)
}

// Start
func (p *iothubService) Start() error {
	tc, subError1 := p.subscribeTopic(message.TopicNameTenant)
	pc, subError2 := p.subscribeTopic(message.TopicNameProduct)
	if subError1 != nil || subError2 != nil {
		return errors.New("iothub failed to subsribe tenant and product topic from kafka")
	}
	p.consumers = append(p.consumers, tc)
	p.consumers = append(p.consumers, pc)
	return nil
}

// Stop
func (p *iothubService) Stop() {
	for _, consumer := range p.consumers {
		consumer.Close()
	}
	p.waitgroup.Wait()
}

// subscribeTopc subscribe topics from apiserver
func (p *iothubService) subscribeTopic(topic string) (sarama.Consumer, error) {
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

			go func(p *iothubService, pc sarama.PartitionConsumer) {
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
func (p *iothubService) handleNotification(topic string, value []byte) error {
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
func (p *iothubService) handleProductNotify(value []byte) error {
	tf := message.ProductTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}
	switch tf.Action {
	case message.ObjectActionRegister:
		_, err := p.CreateProduct(tf.TenantId, tf.ProductId, tf.Replicas)
		return err
	case message.ObjectActionDelete:
		return p.RemoveProduct(tf.TenantId, tf.ProductId)
	}
	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *iothubService) handleTenantNotify(value []byte) error {
	tf := message.TenantTopic{}
	if err := json.Unmarshal(value, &tf); err != nil {
		return err
	}

	switch tf.Action {
	case message.ObjectActionRegister:
		p.CreateTenant(tf.TenantId)
	case message.ObjectActionDelete:
		return p.RemoveTenant(tf.TenantId)
	}
	return nil
}

// addTenant add tenant to iothub
func (p *iothubService) CreateTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		// Create netwrok for each tenant
		networkId, err := p.clustermgr.CreateNetwork(tid)
		if err != nil {
			return err
		}
		p.tenants[tid] = &tenant{
			tid:       tid,
			createdAt: time.Now(),
			products:  make(map[string]*product),
			networkId: networkId,
		}
		return nil
	}
	return fmt.Errorf("tenant '%s' already existed in iothub", tid)
}

// deleteTenant remove tenant from iothub
func (p *iothubService) RemoveTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		return fmt.Errorf("tenant '%s' doesn't exist in iothub", tid)
	}
	t := p.tenants[tid]
	// Remove network
	if err := p.clustermgr.RemoveNetwork(t.networkId); err != nil {
		return err
	}
	// Delete all products
	for name, _ := range t.products {
		if err := p.RemoveProduct(tid, name); err != nil {
			glog.Errorf("iothub remove tenant '%s' product '%s' failed", tid, name)
			// TODO: trying to delete again if failure
		}
	}
	delete(p.tenants, tid)
	return nil
}

func (p *iothubService) isProductExist(tid, pid string) bool {
	if t, found := p.tenants[tid]; found {
		if _, found := t.products[pid]; found {
			return true
		}
	}
	return false
}

// addProduct add product to iothub
func (p *iothubService) CreateProduct(tid, pid string, replicas int32) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.isProductExist(tid, pid) {
		return "", fmt.Errorf("product '%s' of '%s' already exist in iothub", pid, tid)
	}
	serviceName, err := p.clustermgr.CreateService(tid, replicas)
	if err != nil {
		return "", err
	} else {
		t := p.tenants[tid]
		product := &product{tid: tid, pid: pid, createdAt: time.Now(), serviceName: serviceName}
		t.products[pid] = product
	}
	return serviceName, nil
}

// deleteProduct delete product from iothub
func (p *iothubService) RemoveProduct(tid string, pid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isProductExist(tid, pid) {
		return fmt.Errorf("product '%s' of '%s' does not exist in iothub", pid, tid)
	}
	t := p.tenants[tid]
	product := t.products[pid]
	if err := p.clustermgr.RemoveService(product.serviceName); err != nil {
		return err
	}
	delete(t.products, pid)
	return nil
}
