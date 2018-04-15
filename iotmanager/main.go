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

	"github.com/cloustone/sentel/iotmanager/collector"
	"github.com/cloustone/sentel/iotmanager/conductor"
	"github.com/cloustone/sentel/iotmanager/restapi"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

var (
	configFileName = flag.String("c", "/etc/sentel/iotmanager.conf", "config file")
)

func main() {
	flag.Parse()

	glog.Info("Initializing iotmanager...")
	c := config.New("iotmanager")
	c.AddConfig(defaultConfigs)
	c.AddConfigFile(*configFileName)
	service.UpdateServiceConfigs(c, "mongo", "kafka", "zookeeper")

	mgr, _ := service.NewServiceManager("iotmanager", c)
	mgr.AddService(conductor.ServiceFactory{})
	mgr.AddService(restapi.ServiceFactory{})
	mgr.AddService(collector.ServiceFactory{})
	glog.Fatal(mgr.RunAndWait())
}
