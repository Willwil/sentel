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

package iothub

import (
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

// RunWithConfigFile create and start iothub with config file
func RunWithConfigFile(fileName string) error {
	glog.Info("Starting iothub ...")

	// Get configuration
	config, err := core.NewWithConfigFile(fileName)
	if err != nil {
		return err
	}

	// Initialize iothub at startup
	glog.Info("Initializing iothub...")
	if err := InitializeIothub(config); err != nil {
		return err
	}

	// Create service manager according to the configuration
	mgr, err := core.NewServiceManager("iothub", config)
	if err != nil {
		return err
	}
	return mgr.Run()
}

// init initialize default configuration and service before startup
func init() {
	core.RegisterConfigGroup(defaultConfigs)
	core.RegisterService("api", &ApiServiceFactory{})
	core.RegisterService("notify", &NotifyServiceFactory{})
}
