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
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/iothub/cluster"
)

type Iothub struct {
	sync.Once
	config     core.Config
	clustermgr cluster.ClusterManager
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
	}
	return nil
}

// getIothub return global iothub instance used in iothub packet
func getIothub() *Iothub {
	return iothub
}

// addProduct add tenant to iothub
func (p *Iothub) addProduct(tid, pid string, count int32) ([]string, error) {
	return p.clustermgr.CreateBrokers(tid, pid, count)
}

// deleteProduct delete tenant from iothub
func (p *Iothub) deleteProduct(tid string, pid string) error {
	return p.clustermgr.DeleteBrokers(tid, pid)
}

// startProduct start product's brokers
func (p *Iothub) startProduct(tid string, pid string) error {
	return p.clustermgr.StartBrokers(tid, pid)
}

// stopProduct stop product's brokers
func (p *Iothub) stopProduct(tid string, pid string) error {
	return p.clustermgr.StopBrokers(tid, pid)
}
