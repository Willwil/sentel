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
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/mqtt"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/broker/quto"
	"github.com/cloustone/sentel/broker/rpc"
	"github.com/cloustone/sentel/core"
)

// RunWithConfigFile create and start broker
func RunWithConfigFile(fileName string) error {
	return core.RunWithConfigFile("broker", fileName)
}

// init initialize default configurations and services before startup
func init() {
	core.RegisterConfigGroup(defaultConfigs)
	core.RegisterService("rpc", &rpc.ApiServiceFactory{})
	core.RegisterService("metric", &metric.MetricServiceFactory{})
	core.RegisterService("metadata", &metadata.MetadataServiceFactory{})
	core.RegisterService("quto", &quto.QutoServiceFactory{})
	core.RegisterService("event", &event.EventServiceFactory{})
	core.RegisterService("queue", &queue.QueueServiceFactory{})
	core.RegisterService("auth", &auth.AuthServiceFactory{})
	core.RegisterService("mqtt:tcp", &mqtt.MqttFactory{Protocol: mqtt.MqttProtocolTcp})
	core.RegisterService("mqtt:ws", &mqtt.MqttFactory{Protocol: mqtt.MqttProtocolWs})
	core.RegisterService("mqtt:tls", &mqtt.MqttFactory{Protocol: mqtt.MqttProtocolTls})

}
