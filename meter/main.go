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

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/meter/api"
	"github.com/cloustone/sentel/meter/collector"
	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "/etc/sentel/meter.conf", "config file")
)

func main() {
	flag.Parse()
	config := com.NewConfig()
	config.AddConfig(defaultConfigs)
	mgr, _ := com.NewServiceManager("meter", config)
	mgr.AddService(&api.ApiServiceFactory{})
	mgr.AddService(&collector.CollectorServiceFactory{})
	glog.Fatal(mgr.RunAndWait())
}
