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
	"os"

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/conductor/executor"
	"github.com/cloustone/sentel/conductor/indicator"
	"github.com/cloustone/sentel/conductor/restapi"
	"github.com/golang/glog"
)

var (
	configFile = flag.String("c", "/etc/sentel/conductor.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("conductor is starting...")

	config, _ := createConfig(*configFile)
	mgr, _ := com.NewServiceManager("conductor", config)
	mgr.AddService("executor", &executor.ExecutorServiceFactory{})
	mgr.AddService("indicator", &indicator.IndicatorServiceFactory{})
	mgr.AddService("restapi", &restapi.RestapiServiceFactory{})
	glog.Error(mgr.RunAndWait())
}

func createConfig(fileName string) (com.Config, error) {
	config := com.NewConfig()
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	options := map[string]map[string]string{}
	options["iothub"] = make(map[string]string)
	options["iothub"]["kafka"] = os.Getenv("KAFKA_HOST")
	options["iothub"]["mongo"] = os.Getenv("MONGO_HOST")
	config.AddConfig(options)
	return config, nil
}
