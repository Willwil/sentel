//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package broker

import (
	"github.com/cloustone/sentel/broker/auth"
	"github.com/cloustone/sentel/broker/broker"
	"github.com/cloustone/sentel/broker/http"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/mqtt"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/broker/quto"
	"github.com/cloustone/sentel/broker/rpc"
	sub "github.com/cloustone/sentel/broker/subtree"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

// RunWithConfigFile create and start broker
func RunWithConfigFile(fileName string) error {
	glog.Infof("Starting 'broker' server...")

	// Check all registered service
	if err := core.CheckAllRegisteredServices(); err != nil {
		return err
	}
	// Get configuration
	config, err := core.NewConfigWithFile(fileName)
	if err != nil {
		return err
	}
	// Create service manager according to the configuration
	broker, err := broker.NewBroker(config)
	if err != nil {
		return err
	}
	return broker.Run()

}

// init initialize default configurations and services before startup
func init() {
	core.RegisterConfigGroup(defaultConfigs)
	core.RegisterService(rpc.ServiceName, &rpc.ApiServiceFactory{})
	core.RegisterService(metric.ServiceName, &metric.MetricServiceFactory{})
	core.RegisterService(metadata.ServiceName, &metadata.MetadataServiceFactory{})
	core.RegisterService(quto.ServiceName, &quto.QutoServiceFactory{})
	core.RegisterService(queue.ServiceName, &queue.QueueServiceFactory{})
	core.RegisterService(auth.ServiceName, &auth.AuthServiceFactory{})
	core.RegisterService(mqtt.ServiceName, &mqtt.MqttFactory{})
	core.RegisterService(http.ServiceName, &http.HttpServiceFactory{})
	core.RegisterService(sub.ServiceName, &sub.SubServiceFactory{})

}
