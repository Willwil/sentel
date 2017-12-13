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

package cluster

import (
	"errors"
	"time"

	"github.com/cloustone/sentel/core"
)

type ClusterManager interface {
	// CreateBrokers create a number of brokers for tenant and product
	CreateBrokers(tid string, pid string, replicas int32) ([]string, error)
	// StartBroker start the specified broker
	StartBroker(bid string) error
	// StopBroker stop the specified broker
	StopBroker(bid string) error
	// DeleteBroker stop and delete specified broker
	DeleteBroker(bid string) error
	// RollbackBrokers rollback tenant's brokers
	RollbackBrokers(tid, pid string, replicas int32) error
}

type broker struct {
	bid         string       // broker identifier
	tid         string       // tenant identifier
	pid         string       // product identifier
	ip          string       // broker ip address
	status      brokerStatus // broker status
	createdAt   time.Time    // created time for broker
	lastUpdated time.Time    // last updated time for broker
	context     interface{}  // the attached context
}

const (
	brokerStatusInvalid = "invalid"
	brokerStatusCreated = "created"
	brokerStatusStarted = "started"
	brokerStatusStoped  = "stoped"
)

type brokerStatus string

// New retrieve clustermanager instance connected with clustermgr
func New(c core.Config) (ClusterManager, error) {
	v, err := c.String("iothub", "cluster")
	if err != nil {
		return nil, errors.New("iothub cluster manager is not specified")
	}
	switch v {
	case "k8s":
		return newK8sCluster(c)
	case "swarm":
		return newSwarmCluster(c)
	}
	return nil, errors.New("iothub cluster manager is not specified")
}
