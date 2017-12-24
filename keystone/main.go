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

	com "github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/iothub/restapi"

	"github.com/golang/glog"
)

var (
	configFileName = flag.String("c", "/etc/sentel/keystone.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("Starting keystone server...")

	config, _ := createConfig(*configFileName)
	mgr, err := com.NewServiceManager("keystone", config)
	if err != nil {
		glog.Fatalf("keystone create failed: '%s'", err.Error())
	}
	mgr.AddService(restapi.ServiceFactory{})
	glog.Fatal(mgr.RunAndWait())
}

func createConfig(fileName string) (com.Config, error) {
	config := com.NewConfig()
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	options := map[string]map[string]string{}
	options["keystone"] = map[string]string{}
	options["keystone"]["mongo"] = os.Getenv("MONGO_HOST")
	config.AddConfig(options)
	return config, nil
}
