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

	"github.com/cloustone/sentel/mns/apiservice"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

var (
	configFile = flag.String("c", "/etc/sentel/mns.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("message notify service is starting...")

	config, _ := createConfig(*configFile)
	mgr, _ := service.NewServiceManager("mns", config)
	mgr.AddService(apiservice.ServiceFactory{})
	glog.Error(mgr.RunAndWait())
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New("mns")
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	k := os.Getenv("KAFKA_HOST")
	m := os.Getenv("MONGO_HOST")
	if k != "" && m != "" {
		config.AddConfigItem("kafka", k)
		config.AddConfigItem("mongo", m)
	}
	return config, nil
}
