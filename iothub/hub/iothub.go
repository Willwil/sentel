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
	network   string              `json:"network"`
}

type product struct {
	pid         string    `json:"productId"`
	tid         string    `json:"tenantId"`
	createdAt   time.Time `json:"createdAt"`
	serviceName string    `json:"serviceName"`
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
func (p *Iothub) createTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		// Create netwrok for each tenant
		network, err := p.clustermgr.CreateNetwork(tid)
		if err != nil {
			return err
		}
		p.tenants[tid] = &tenant{
			tid:       tid,
			createdAt: time.Now(),
			products:  make(map[string]*product),
			network:   network,
		}
		return nil
	}
	return fmt.Errorf("tenant '%s' already existed in iothub")
}

// deleteTenant remove tenant from iothub
func (p *Iothub) removeTenant(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, found := p.tenants[tid]; !found {
		return fmt.Errorf("tenant '%s' doesn't exist in iothub")
	}
	t := p.tenants[tid]
	// Remove network
	if err := p.clustermgr.RemoveNetwork(t.network); err != nil {
		return err
	}
	// Delete all products
	for name, _ := range t.products {
		if err := p.removeProduct(tid, name); err != nil {
			glog.Errorf("iothub remove tenant '%s' product '%s' failed", tid, name)
			// TODO: trying to delete again if failure
		}
	}
	delete(p.tenants, tid)
	return nil
}

func (p *Iothub) isProductExist(tid, pid string) bool {
	if t, found := p.tenants[tid]; found {
		if _, found := t.products[pid]; found {
			return true
		}
	}
	return false
}

// addProduct add product to iothub
func (p *Iothub) createProduct(tid, pid string, replicas int32) (string, error) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if p.isProductExist(tid, pid) {
		return "", fmt.Errorf("product '%s' of '%s' already exist in iothub", pid, tid)
	}
	serviceName, err := p.clustermgr.CreateService(tid, pid, replicas)
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
func (p *Iothub) removeProduct(tid string, pid string) error {
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
