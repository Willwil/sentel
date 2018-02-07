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

package conductor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/cloustone/sentel/iotmanager/mgrdb"
	"github.com/cloustone/sentel/pkg/cluster"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	sd "github.com/cloustone/sentel/pkg/service-discovery"
	"github.com/golang/glog"
)

type conductorService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	clustermgr  cluster.ClusterManager
	tenants     map[string]*mgrdb.Tenant
	mutex       sync.Mutex
	consumer    message.Consumer
	dbconn      mgrdb.ManagerDB
	recoverChan chan interface{}
	quitChan    chan interface{}
}

var (
	logger = log.New(os.Stderr, "[kafka]", log.LstdFlags)
)

const SERVICE_NAME = "conductor"

type ServiceFactory struct{}

func (m ServiceFactory) New(c config.Config) (service.Service, error) {
	clustermgr, err := cluster.New(c)
	if err != nil {
		return nil, err
	}
	dbconn, err := mgrdb.New(c)
	if err != nil {
		return nil, err
	}
	discovery, err := sd.New(c)
	if err != nil {
		return nil, err
	}
	consumer, err := message.NewConsumer(c, "iotmanager")
	if err != nil {
		return nil, err
	}
	clustermgr.SetServiceDiscovery(discovery)
	return &conductorService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		clustermgr:  clustermgr,
		tenants:     make(map[string]*mgrdb.Tenant),
		mutex:       sync.Mutex{},
		dbconn:      dbconn,
		recoverChan: make(chan interface{}),
		quitChan:    make(chan interface{}),
		consumer:    consumer,
	}, nil
}

// Name
func (p *conductorService) Name() string { return SERVICE_NAME }

// Initialize load iotmanager data and recovery from scratch
func (p *conductorService) Initialize() error {
	tenants := p.dbconn.GetAllTenants()
	if len(tenants) > 0 {
		for _, t := range tenants {
			p.tenants[t.TenantId] = &t
		}
		p.recoverChan <- true
	}
	return nil
}

// Start
func (p *conductorService) Start() error {
	// subscribe topic and start message consumer
	err1 := p.consumer.Subscribe(message.TopicNameTenant, p.messageHandlerFunc, nil)
	err2 := p.consumer.Subscribe(message.TopicNameProduct, p.messageHandlerFunc, nil)
	if err1 != nil || err2 != nil {
		return errors.New("failed to subscribe topic from message server")
	}
	if err := p.consumer.Start(); err != nil {
		return err
	}
	p.waitgroup.Add(1)
	go func(p *conductorService) {
		defer p.waitgroup.Done()
		for {
			select {
			case <-p.quitChan:
				return
			case <-p.recoverChan:
				p.recoverStartup()
			}
		}
	}(p)
	return nil
}

// Stop
func (p *conductorService) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	p.consumer.Close()
	close(p.quitChan)
	close(p.recoverChan)
}

// recoverStartup load hub data and confirm service status
func (p *conductorService) recoverStartup() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	// create default network if network is specified
	network, err := p.config.String("network")
	if err == nil && network != "" {
		p.clustermgr.CreateNetwork(network)
	}

	// load all tenants form mgrdb and create services if tenant service no exist
	retries := []*mgrdb.Tenant{}
	for tid, t := range p.tenants {
		if t.ServiceState != cluster.ServiceStateNone {
			if _, err := p.clustermgr.IntrospectService(tid); err != nil {
				env := []string{
					"KAFKA_HOST=kafka:9092",
					"MONGO_HOST=mongo:27017",
					fmt.Sprintf("BROKER_TENANT=%s", tid),
				}
				spec := cluster.ServiceSpec{
					TenantId:    tid,
					NetworkId:   network,
					Replicas:    t.InstanceReplicas,
					Environment: env,
					Image:       "sentel/broker",
					ServiceName: fmt.Sprintf("tenant_%s", tid),
				}
				// remember all failed recovery and try again at later
				if _, err := p.clustermgr.CreateService(spec); err != nil {
					retries = append(retries, t)
				}
			}
		}
	}
	// retry to recover again
	for _, t := range retries {
		spec := cluster.ServiceSpec{
			TenantId:  t.TenantId,
			NetworkId: network,
			Replicas:  t.InstanceReplicas,
		}
		if _, err := p.clustermgr.CreateService(spec); err != nil {
			glog.Errorf("service '%s' recovery failed", t.TenantId)
		}
	}
}

func (p *conductorService) messageHandlerFunc(msg message.Message, ctx interface{}) error {
	topic := msg.Topic()
	glog.Infof("receive message from topic '%s'", topic)

	switch topic {
	case message.TopicNameTenant:
		return p.handleTenantNotify(msg)
	case message.TopicNameProduct:
		return p.handleProductNotify(msg)
	default:
		return fmt.Errorf("invalid topic '%s'", topic)
	}
}

// handleProductNotify handle notification about product from api server
func (p *conductorService) handleProductNotify(msg message.Message) error {
	tf, ok := msg.(*message.Product)
	if !ok || tf == nil {
		return errors.New("invalid product notification")
	}
	switch tf.Action {
	case message.ActionCreate:
		_, err := p.createProduct(tf.TenantId, tf.ProductId, tf.Replicas)
		return err
	case message.ActionRemove:
		return p.removeProduct(tf.TenantId, tf.ProductId)
	}
	return nil
}

// handleTenantNotify handle notification about tenant from api server
func (p *conductorService) handleTenantNotify(msg message.Message) error {
	tf, ok := msg.(*message.Tenant)
	if !ok || tf == nil {
		return errors.New("invalid product notification")
	}

	switch tf.Action {
	case message.ActionCreate:
		p.createTenant(tf.TenantId)
	case message.ActionRemove:
		return p.removeTenant(tf.TenantId)
	}
	return nil
}

// addTenant add tenant to iotmanager
func (p *conductorService) createTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		p.tenants[tid] = &mgrdb.Tenant{
			TenantId:     tid,
			CreatedAt:    time.Now(),
			Products:     make(map[string]*mgrdb.Product),
			ServiceState: cluster.ServiceStateNone,
		}
		p.dbconn.AddTenant(p.tenants[tid])
		return nil
	}
	return fmt.Errorf("tenant '%s' already existed in iotmanager", tid)
}

// deleteTenant remove tenant from iotmanager
func (p *conductorService) removeTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		return fmt.Errorf("tenant '%s' doesn't exist in iotmanager", tid)
	}
	t := p.tenants[tid]
	// remove service created early
	if t.ServiceState != cluster.ServiceStateNone {
		p.clustermgr.RemoveService(t.ServiceId)
	}
	// Delete all products
	for name := range t.Products {
		if err := p.removeProduct(tid, name); err != nil {
			glog.Errorf("remove tenant '%s' product '%s' failed", tid, name)
		}
	}
	// Remove service
	delete(p.tenants, tid)
	p.dbconn.RemoveTenant(tid)
	return nil
}

func (p *conductorService) isProductExist(tid, pid string) bool {
	if t, found := p.tenants[tid]; found {
		if _, found := t.Products[pid]; found {
			return true
		}
	}
	return false
}

// addProduct add product to iotmanager
func (p *conductorService) createProduct(tid string, pid string, replicas int32) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.isProductExist(tid, pid) {
		return "", fmt.Errorf("product '%s' of '%s' already exist in iotmanager", pid, tid)
	}
	t := p.tenants[tid]
	if t.ServiceState == cluster.ServiceStateNone {
		network, err := p.config.String("network")
		if err != nil {
			network = ""
		}
		spec := cluster.ServiceSpec{
			TenantId:  tid,
			NetworkId: network,
			Replicas:  replicas,
			Image:     "sentel/broker",
		}
		serviceId, err := p.clustermgr.CreateService(spec)
		if err != nil {
			return "", err
		}
		t.ServiceState = cluster.ServiceStateStarted
		t.ServiceId = serviceId
		t.ServiceName = tid
	}
	product := &mgrdb.Product{ProductId: pid, CreatedAt: time.Now()}
	t.Products[pid] = product
	p.dbconn.AddProduct(tid, pid)
	return t.ServiceId, nil
}

// deleteProduct delete product from iotmanager
func (p *conductorService) removeProduct(tid string, pid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isProductExist(tid, pid) {
		return fmt.Errorf("product '%s' of '%s' does not exist in iotmanager", pid, tid)
	}
	t := p.tenants[tid]
	delete(t.Products, pid)
	if len(t.Products) == 0 {
		p.clustermgr.RemoveService(t.ServiceId)
		t.ServiceState = cluster.ServiceStateNone
	}
	p.dbconn.RemoveProduct(tid, pid)
	return nil
}
