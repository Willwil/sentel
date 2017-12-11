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

	"github.com/cloustone/sentel/core"
)

type ClusterManager interface {
	// CreateBrokers create a number of brokers for tenant and product
	CreateBrokers(tid string, pid string, replicas int32) ([]string, error)
	// StartBroker start the specified broker
	StartBroker(bid string) error
	// StopBroker stop the specified broker
	StopBroker(bid string) error
	// StartBrokers start product's all brokers
	StartBrokers(tid, pid string) error
	// StopBrokers stop product's all brokers
	StopBrokers(tid, pid string) error
	// DeleteBrokers delete all product's brokers
	DeleteBrokers(tid, pid string) error
	// DeleteBroker stop and delete specified broker
	DeleteBroker(bid string) error
	// RollbackBrokers rollback tenant's brokers
	RollbackBrokers(tid, pid string, replicas int32) error
}

// New retrieve clustermanager instance connected with clustermgr
func New(c core.Config) (ClusterManager, error) {
	v, err := c.String("iothub", "cluster")
	if err != nil {
		return nil, errors.New("iothub cluster manager is not specified")
	}
	switch v {
	case "k8s":
		return newK8sClusterManager(c)
	}
	return nil, errors.New("iothub cluster manager is not specified")
}
