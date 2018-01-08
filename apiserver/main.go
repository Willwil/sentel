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

	"github.com/cloustone/sentel/apiserver/console"
	"github.com/cloustone/sentel/apiserver/management"
	"github.com/cloustone/sentel/apiserver/swagger"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"

	"github.com/golang/glog"
)

var (
	configFileName = flag.String("c", "/etc/sentel/apiserver.conf", "config file")
	swaggerFile    = flag.String("s", "../apiserver/swagger/console_swagger.yaml", "swagger file")
)

func main() {
	flag.Parse()
	glog.Info("Starting api server...")

	config, _ := createConfig(*configFileName)
	mgr, err := service.NewServiceManager("apiserver", config)
	if err != nil {
		glog.Fatalf("apiserver create failed: '%s'", err.Error())
	}
	mgr.AddService(console.ServiceFactory{})
	mgr.AddService(management.ServiceFactory{})
	if _, err := config.String("apiserver", "swagger"); err == nil {
		mgr.AddService(swagger.ServiceFactory{})
	}
	glog.Fatal(mgr.RunAndWait())
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New()
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)

	kafka := os.Getenv("KAFKA_HOST")
	mongo := os.Getenv("MONGO_HOST")
	if kafka != "" && mongo != "" {
		options := make(map[string]map[string]string)
		options["apiserver"] = make(map[string]string)
		options["apiserver"]["kafka"] = os.Getenv("KAFKA_HOST")
		options["apiserver"]["mongo"] = os.Getenv("MONGO_HOST")
		config.AddConfig(options)
	}
	if _, err := config.String("apiserver", "swagger"); err == nil {
		config.AddConfig(map[string]map[string]string{
			"swagger": {
				"path": *swaggerFile,
			}})
	}

	return config, nil
}
