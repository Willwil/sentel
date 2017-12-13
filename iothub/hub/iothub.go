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
	"errors"
	"fmt"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/iothub/cluster"
	"github.com/golang/glog"
)

type Iothub struct {
	sync.Once
	config     core.Config
	clustermgr cluster.ClusterManager
	tenants    map[string]*tenant
	mutex      sync.Mutex
}
type tenant struct {
	tid       string              `json:"tenantId"`
	createdAt time.Time           `json:"createdAt"`
	products  map[string]*product `json:"products"`
}

type product struct {
	pid       string    `json:"productId"`
	tid       string    `json:"tenantId"`
	createdAt time.Time `json:"createdAt"`
	brokers   []string  "json:‘brokers"
}

var (
	iothub *Iothub
)

// InitializeIothub create iothub global instance at startup time
func InitializeIothub(c core.Config) error {
	// check mongo db configuration
	hosts, err := c.String("iothub", "mongo")
	if err != nil || hosts == "" {
		return errors.New("Invalid mongo configuration")
	}
	// try connect with mongo db
	session, err := mgo.DialWithTimeout(hosts, 5*time.Second)
	if err != nil {
		return err
	}
	session.Close()

	clustermgr, err := cluster.New(c)
	if err != nil {
		return err
	}
	iothub = &Iothub{
		config:     c,
		clustermgr: clustermgr,
		tenants:    make(map[string]*tenant),
		mutex:      sync.Mutex{},
	}
	return nil
}

// getIothub return global iothub instance used in iothub packet
func getIothub() *Iothub {
	return iothub
}

// addTenant add tenant to iothub
func (p *Iothub) addTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		p.tenants[tid] = &tenant{
			tid:       tid,
			createdAt: time.Now(),
			products:  make(map[string]*product),
		}
		return nil
	}
	return fmt.Errorf("tenant '%s' already existed in iothub")
}

// deleteTenant remove tenant from iothub
func (p *Iothub) DeleteTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		return fmt.Errorf("tenant '%s' doesn't exist in iothub")
	}
	// Delete all products
	t := p.tenants[tid]
	for name, _ := range t.products {
		if err := p.deleteProduct(tid, name); err != nil {
			glog.Errorf("iothub remove tenant '%s' product '%s' failed", tid, name)
			// TODO: trying to delete again if failure
		}
	}
	delete(p.tenants, tid)
	return nil
}

func (p *Iothub) isProductExist(tid, pid string) bool {
	return false
}

// addProduct add product to iothub
func (p *Iothub) addProduct(tid, pid string, replicas int32) ([]string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.isProductExist(tid, pid) {
		return nil, fmt.Errorf("product '%s' of '%s' already exist in iothub", pid, tid)
	}
	brokers, err := p.clustermgr.CreateBrokers(tid, pid, replicas)
	if err != nil {
		return nil, err
	} else {
		t := p.tenants[tid]
		product := &product{tid: tid, pid: pid, createdAt: time.Now(), brokers: brokers}
		t.products[pid] = product
	}
	return brokers, nil
}

// deleteProduct delete product from iothub
func (p *Iothub) deleteProduct(tid string, pid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isProductExist(tid, pid) {
		return fmt.Errorf("product '%s' of '%s' does not exist in iothub", pid, tid)
	}
	t := p.tenants[tid]
	product := t.products[pid]
	for _, bid := range product.brokers {
		if err := p.clustermgr.DeleteBroker(bid); err != nil {
			return err
		}
	}
	delete(t.products, pid)
	return nil
}

// startProduct start product's brokers
func (p *Iothub) startProduct(tid string, pid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isProductExist(tid, pid) {
		return fmt.Errorf("product '%s' of '%s' does not exist in iothub", pid, tid)
	}
	t := p.tenants[tid]
	product := t.products[pid]
	for _, bid := range product.brokers {
		if err := p.clustermgr.StartBroker(bid); err != nil {
			return err
		}
	}
	return nil
}

// stopProduct stop product's brokers
func (p *Iothub) stopProduct(tid string, pid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if !p.isProductExist(tid, pid) {
		return fmt.Errorf("product '%s' of '%s' does not exist in iothub", pid, tid)
	}
	t := p.tenants[tid]
	product := t.products[pid]
	for _, bid := range product.brokers {
		if err := p.clustermgr.StopBroker(bid); err != nil {
			return err
		}
	}
	return nil

}
