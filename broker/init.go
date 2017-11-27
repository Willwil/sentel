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
	core.RegisterConfigGroup(defaultConfigs)
	// Get configuration
	config, err := core.NewConfigWithFile(fileName)
	if err != nil {
		return err
	}
	broker.RegisterService(rpc.ServiceName, &rpc.ApiServiceFactory{})
	broker.RegisterService(metric.ServiceName, &metric.MetricServiceFactory{})
	broker.RegisterService(metadata.ServiceName, &metadata.MetadataServiceFactory{})
	broker.RegisterService(quto.ServiceName, &quto.QutoServiceFactory{})
	broker.RegisterService(queue.ServiceName, &queue.QueueServiceFactory{})
	broker.RegisterService(auth.ServiceName, &auth.AuthServiceFactory{})
	broker.RegisterService(mqtt.ServiceName, &mqtt.MqttFactory{})
	broker.RegisterService(http.ServiceName, &http.HttpServiceFactory{})
	broker.RegisterService(sub.ServiceName, &sub.SubServiceFactory{})

	// Create service manager according to the configuration
	broker, err := broker.NewBroker(config)
	if err != nil {
		return err
	}
	return broker.Run()
}

// RunWithConfig create and start broker from loaded configuration
func RunWithConfig(c core.Config) error {
	broker.RegisterService(rpc.ServiceName, &rpc.ApiServiceFactory{})
	broker.RegisterService(metric.ServiceName, &metric.MetricServiceFactory{})
	broker.RegisterService(metadata.ServiceName, &metadata.MetadataServiceFactory{})
	broker.RegisterService(quto.ServiceName, &quto.QutoServiceFactory{})
	broker.RegisterService(queue.ServiceName, &queue.QueueServiceFactory{})
	broker.RegisterService(auth.ServiceName, &auth.AuthServiceFactory{})
	broker.RegisterService(mqtt.ServiceName, &mqtt.MqttFactory{})
	broker.RegisterService(http.ServiceName, &http.HttpServiceFactory{})
	broker.RegisterService(sub.ServiceName, &sub.SubServiceFactory{})

	// Create service manager according to the configuration
	broker, err := broker.NewBroker(c)
	if err != nil {
		return err
	}
	return broker.Run()
}
