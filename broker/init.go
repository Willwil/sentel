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

package main

import (
	"github.com/cloustone/sentel/broker/auth"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/http"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/mqtt"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/broker/quto"
	"github.com/cloustone/sentel/broker/rpc"
	"github.com/cloustone/sentel/broker/sessionmgr"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

// RunWithConfigFile create and start broker
func RunWithConfig(fileName string, opts map[string]map[string]string) error {
	glog.Infof("Starting 'broker' server...")
	// Get configuration
	config, _ := core.NewConfigWithFile(fileName)
	config.AddConfigs(opts)
	registerService(event.ServiceName, event.New)
	registerService(queue.ServiceName, queue.New)
	registerService(sessionmgr.ServiceName, sessionmgr.New)
	registerService(auth.ServiceName, auth.New)
	registerService(rpc.ServiceName, rpc.New)
	registerService(metric.ServiceName, metric.New)
	registerService(metadata.ServiceName, metadata.New)
	registerService(quto.ServiceName, quto.New)
	registerService(mqtt.ServiceName, mqtt.New)
	registerService(http.ServiceName, http.New)

	// Create service manager according to the configuration
	broker, err := NewBroker(config)
	if err != nil {
		return err
	}
	return broker.Run()
}

func init() {
	core.RegisterConfigGroup(defaultConfigs)
}
