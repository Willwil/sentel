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
	"github.com/cloustone/sentel/broker/api"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/mqtt"
	"github.com/cloustone/sentel/core"

	"github.com/golang/glog"
)

// RunWithConfigFile create and start broker
func RunWithConfigFile(fileName string) error {
	glog.Info("Starting mqtt broker...")

	// Get configuration
	config, err := core.NewWithConfigFile(fileName)
	if err != nil {
		return err
	}
	// Create service manager according to the configuration
	mgr, err := core.NewServiceManager("broker", config)
	if err != nil {
		glog.Fatal("Failed to launch ServiceManager")
		return err
	}
	return mgr.Run()
}

// init initialize default configurations and services before startup
func init() {
	for group, values := range defaultConfigs {
		core.RegisterConfig(group, values)
	}
	core.RegisterService("mqtt:tcp", &mqtt.MqttFactory{})
	core.RegisterService("api", &api.ApiServiceFactory{})
	core.RegisterService("metric", &metric.MetricServiceFactory{})
}
