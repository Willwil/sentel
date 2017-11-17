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
	"errors"
	"fmt"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

type Iothub struct {
	sync.Once
	config     core.Config
	mutex      sync.Mutex
	tenants    map[string]*Tenant
	brokers    map[string]*Broker
	clustermgr *clusterManager
}

type Tenant struct {
	id           string    `json:"tenantId"`
	createdAt    time.Time `json:"createdAt"`
	brokersCount int32     `json:"brokersCount"`
	brokers      map[string]*Broker
}

var (
	_iothub *Iothub
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

	clustermgr, err := newClusterManager(c)
	if err != nil {
		return err
	}
	_iothub = &Iothub{
		config:     c,
		mutex:      sync.Mutex{},
		tenants:    make(map[string]*Tenant),
		brokers:    make(map[string]*Broker),
		clustermgr: clustermgr,
	}
	return nil
}

// getIothub return global iothub instance used in iothub packet
func getIothub() *Iothub {
	return _iothub
}

// getTenantFromDatabase retrieve tenant from database
func (p *Iothub) getTenantFromDatabase(tid string) (*Tenant, error) {
	hosts, err := p.config.String("iothub", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	tenant := Tenant{}
	c := session.DB("iothub").C("tenants")
	if err := c.Find(bson.M{"tenantId": tid}).One(&tenant); err != nil {
		return nil, err
	}
	return &tenant, nil
}

// addTenant add tenant to iothub
func (p *Iothub) addTenant(tid string) error {
	// check wether the tenant already exist
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.tenants[tid]; ok {
		return fmt.Errorf("tenant(%s) already exist", tid)
	}

	// retrieve tenant from database
	tenant, err := p.getTenantFromDatabase(tid)
	if err != nil {
		return err
	}
	// save new tenant into iothub
	p.tenants[tid] = tenant

	// create brokers accroding to tenant's request
	brokers, err := p.clustermgr.createBrokers(tid, tenant.brokersCount)
	if err != nil {
		// TODO should we update database status
		glog.Fatalf("Failed to created brokers for tenant(%s)", tid)
	}
	for _, broker := range brokers {
		p.brokers[broker.bid] = broker
	}

	return nil
}

// deleteTenant delete tenant from iothub
func (p *Iothub) deleteTenant(tid string) error {
	// check wether the tenant already exist
	p.mutex.Lock()
	defer p.mutex.Unlock()
	if _, ok := p.tenants[tid]; ok {
		return fmt.Errorf("tenant(%s) already exist", tid)
	}

	if err := p.clustermgr.deleteBrokers(tid); err != nil {
		return err
	}

	// delete tenant from iothub
	tenant := p.tenants[tid]
	for bid, _ := range tenant.brokers {
		delete(p.brokers, bid)
	}
	delete(p.tenants, tid)

	// remove tenant from database
	hosts, err := p.config.String("iothub", "mongo")
	session, err := mgo.Dial(hosts)
	if err != nil {
		return err
	}
	defer session.Close()

	c := session.DB("iothub").C("tenants")
	if err := c.Remove(bson.M{"tenantId": tid}); err != nil {
		return err
	}
	return nil
}

// updateTenant update a local tenant with database
func (p *Iothub) updateTenant(tid string) error {
	// check wether the tenant already exist
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.tenants[tid]; ok {
		return fmt.Errorf("tenant(%s) already exist", tid)
	}

	lt := p.tenants[tid]

	// retrieve tenant from database
	tenant, err := p.getTenantFromDatabase(tid)
	if err != nil {
		return err
	}

	// rollback brokers if local tenant and tenant in database is not same
	if lt.brokersCount != tenant.brokersCount {
		return p.clustermgr.rollbackTenantBrokers(tenant)
	}

	return nil
}

// addBroker add broker to tenant
func (p *Iothub) addBroker(tid string, b *Broker) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.tenants[tid]; !ok {
		return fmt.Errorf("invalid tenant(%s)", tid)
	}
	tenant := p.tenants[tid]
	tenant.brokers[b.bid] = b
	p.brokers[b.bid] = b

	return nil
}

// deleteBroker delete broker from iothub
func (p *Iothub) deleteBroker(bid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.brokers[bid]; !ok {
		return fmt.Errorf("invalid broker(%s)", bid)
	}
	broker := p.brokers[bid]
	tenant := p.tenants[broker.tid]

	// stop broker when deleted
	if err := p.clustermgr.deleteBroker(broker); err != nil {
		return err
	}
	delete(p.brokers, bid)
	delete(tenant.brokers, bid)

	return nil
}

// startBroker start a tenant's broker
func (p *Iothub) startBroker(bid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.brokers[bid]; !ok {
		return fmt.Errorf("invalid broker(%s)", bid)
	}

	broker := p.brokers[bid]
	if err := p.clustermgr.startBroker(broker); err != nil {
		return err
	}
	broker.status = BrokerStatusStarted
	return nil
}

// stopBroker stop a broker
func (p *Iothub) stopBroker(bid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.brokers[bid]; !ok {
		return fmt.Errorf("invalid broker(%s)", bid)
	}

	broker := p.brokers[bid]
	if err := p.clustermgr.stopBroker(broker); err != nil {
		return err
	}
	broker.status = BrokerStatusStoped

	return nil
}

// startTenantBrokers start tenant's all broker
func (p *Iothub) startTenantBrokers(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.tenants[tid]; !ok {
		return fmt.Errorf("invalid tenant(%s)", tid)
	}

	tenant := p.tenants[tid]
	for bid, broker := range tenant.brokers {
		if broker.status != BrokerStatusStarted {
			if err := p.clustermgr.startBroker(broker); err != nil {
				glog.Errorf("Failed to start broker(%s) for tenant(%s)", bid, tid)
				continue
			}
		}
	}
	return nil
}

// stopTenantBrokers stop tenant's all brokers
func (p *Iothub) stopTenantBrokers(tid string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.tenants[tid]; !ok {
		return fmt.Errorf("invalid tenant(%s)", tid)
	}

	tenant := p.tenants[tid]
	for bid, broker := range tenant.brokers {
		if broker.status != BrokerStatusStoped {
			if err := p.clustermgr.stopBroker(broker); err != nil {
				glog.Errorf("Failed to stop broker(%s) for tenant(%s)", bid, tid)
				continue
			}
		}
	}
	return nil
}

// getTenant retrieve a tenant by id
func (p *Iothub) getTenant(id string) *Tenant {
	return p.tenants[id]
}

// getBroker retrieve a broker by id
func (p *Iothub) getBroker(id string) *Broker {
	return p.brokers[id]
}

// setBrokerStatus set broker's status
func (p *Iothub) setBrokerStatus(bid string, status BrokerStatus) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if _, ok := p.brokers[bid]; !ok {
		return fmt.Errorf("Invalid broker(%s)", bid)
	}
	broker := p.brokers[bid]
	if broker.status != status {
		var err error
		switch status {
		case BrokerStatusStarted:
			err = p.clustermgr.startBroker(broker)
		case BrokerStatusStoped:
			err = p.clustermgr.stopBroker(broker)
		default:
			err = fmt.Errorf("Invalid broker status to set for broker(%s)", bid)
		}
		if err != nil {
			return err
		}
	}
	broker.status = status
	return nil
}

// getBrokerStatus retrieve broker's status
func (p *Iothub) getBrokerStatus(bid string) BrokerStatus {
	if _, ok := p.brokers[bid]; !ok {
		return BrokerStatusInvalid
	}
	return p.brokers[bid].status
}
