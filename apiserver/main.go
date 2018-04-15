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
	"fmt"
	"os"

	"github.com/cloustone/sentel/apiserver/console"
	"github.com/cloustone/sentel/apiserver/management"
	"github.com/cloustone/sentel/apiserver/swagger"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"

	"github.com/golang/glog"
)

var (
	configFileName = flag.String("c", "/etc/sentel/apiserver.conf", "config file")
	swaggerFile    = flag.String("s", "apiserver/swagger/console_swagger.yaml", "swagger file")
)

func main() {
	flag.Parse()
	glog.Info("Starting api server...")

	config, _ := createConfig(*configFileName)
	mgr, _ := service.NewServiceManager("apiserver", config)
	mgr.AddService(console.ServiceFactory{})
	mgr.AddService(management.ServiceFactory{})
	if _, err := config.String("swagger"); err == nil {
		mgr.AddService(swagger.ServiceFactory{})
	}
	if err := mgr.RunAndWait(); err != nil {
		glog.Fatal(err)
	}
	glog.Info("apiserver is normally terminated")
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New("apiserver")
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)

	for _, e := range os.Environ() {
		fmt.Println(e)
	}
	kafka := os.Getenv("KAFKA_PORT")
	mongo := os.Getenv("MONGO_PORT")
	if kafka != "" && mongo != "" {
		config.AddConfigItem("kafka", kafka)
		config.AddConfigItem("mongo", mongo)
	}
	if _, err := config.String("swagger"); err == nil {
		config.AddConfigItemWithSection("swagger", "path", *swaggerFile)
	}

	return config, nil
}
