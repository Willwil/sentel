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
	"flag"

	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/iothub/hub"
	"github.com/golang/glog"
)

var (
	configFileName = flag.String("c", "/etc/sentel/iothub.conf", "config file")
)

func main() {
	flag.Parse()

	glog.Info("Starting iothub ...")
	core.RegisterConfigGroup(defaultConfigs)
	core.RegisterService("api", &hub.ApiServiceFactory{})
	core.RegisterService("notify", &hub.NotifyServiceFactory{})

	// Get configuration
	config, err := core.NewConfigWithFile(*configFileName)
	if err != nil {
		glog.Fatal(err)
	}
	// Initialize iothub at startup
	glog.Info("Initializing iothub...")
	if err := hub.InitializeIothub(config); err != nil {
		glog.Fatal(err)
	}

	// Create service manager according to the configuration
	mgr, err := core.NewServiceManager("iothub", config)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Error(mgr.RunAndWait())
}
