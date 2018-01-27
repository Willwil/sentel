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

	"github.com/golang/glog"

	"github.com/cloustone/sentel/iothub/watcher"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
)

var (
	configFileName = flag.String("c", "/etc/sentel/gateway.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("Starting  iothub...")

	config, _ := createConfig(*configFileName)
	mgr, err := service.NewServiceManager("iothub", config)
	if err != nil {
		glog.Fatalf("iothub create failed: '%s'", err.Error())
	}
	mgr.AddService(watcher.ServiceFactory{})
	glog.Fatal(mgr.RunAndWait())
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New("iothub")
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	//config.AddConfigItem("zookeeper", os.Getenv("ZOOKEEPER_HOST"))
	return config, nil
}
