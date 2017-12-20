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
	configFile     = flag.String("c", "/etc/sentel/conductor.conf", "config file")
	defaultConfigs = map[string]map[string]string{
		"conductor": {
			"loglevel":        "debug",
			"kafka":           "localhost:9092",
			"mongo":           "localhost:27017",
			"connect_timeout": "5",
		},
	}
)

func main() {
	flag.Parse()
	glog.Info("conductor is starting...")

	com.RegisterServiceWithConfig("executor", &executor.ExecutorServiceFactory{}, executor.Configs)
	com.RegisterServiceWithConfig("indicator", &indicator.IndicatorServiceFactory{}, indicator.Configs)
	com.RegisterServiceWithConfig("restapi", &restapi.RestapiServiceFactory{}, restapi.Configs)
	config, _ := createConfig(*configFile)
	// Create service manager according to the configuration
	mgr, err := com.NewServiceManager("conductor", config)
	if err != nil {
		glog.Fatal(err)
	}
	glog.Error(mgr.RunAndWait())
}

func createConfig(fileName string) (com.Config, error) {
	options := map[string]map[string]string{}
	options["iothub"] = make(map[string]string)
	options["iothub"]["kafka"] = os.Getenv("KAFKA_HOST")
	options["iothub"]["mongo"] = os.Getenv("MONGO_HOST")
	// Get configuration
	com.RegisterConfigGroup(defaultConfigs)
	config, _ := com.NewConfigWithFile(fileName)
	config.AddConfigs(options)
	return config, nil
}
