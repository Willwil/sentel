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
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/http"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/mqtt"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/broker/quto"
	"github.com/cloustone/sentel/broker/rpc"
	"github.com/cloustone/sentel/broker/sessionmgr"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

var (
	configFile = flag.String("c", "/etc/sentel/broker.conf", "config file")
	tenant     = flag.String("t", "", "tenant id")
	protocol   = flag.String("P", "tcp", "mqtt access protocol, tcp|tls|ws|https")
	listen     = flag.String("l", "localhost:1883", "mqtt broker listen address")
)

func main() {
	flag.Parse()

	glog.Infof("Starting 'broker' server...")
	config, err := createConfig(*configFile)
	if err != nil {
		fmt.Println(err.Error())
		flag.PrintDefaults()
		return
	}
	// Create service manager according to the configuration
	broker, _ := NewBroker(config)
	broker.AddService(event.ServiceFactory{})
	broker.AddService(queue.ServiceFactory{})
	broker.AddService(sessionmgr.ServiceFactory{})
	//broker.AddService(auth.ServiceFactory{})
	broker.AddService(rpc.ServiceFactory{})
	broker.AddService(metric.ServiceFactory{})
	broker.AddService(metadata.ServiceFactory{})
	broker.AddService(quto.ServiceFactory{})
	broker.AddService(mqtt.ServiceFactory{})
	broker.AddService(http.ServiceFactory{})

	glog.Fatal(broker.Run())
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New()
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	options := map[string]map[string]string{}
	options["broker"] = map[string]string{}
	options["broker"]["kafka"] = os.Getenv("KAFKA_HOST")
	options["broker"]["mongo"] = os.Getenv("MONGO_HOST")
	options["mqtt"] = map[string]string{}

	// tenant and product must be set
	if *tenant == "" {
		*tenant = os.Getenv("BROKER_TENANT")
		if *tenant == "" {
			return nil, errors.New("teant must be specified for broker")
		}
	}
	options["broker"]["tenant"] = *tenant

	// Mqtt protocol
	if *protocol != "" {
		switch *protocol {
		case "tcp", "ws", "tls", "https":
			options["broker"]["protocol"] = *protocol
			if *listen != "" {
				options["mqtt"][*protocol] = *listen
			}
		default:
			return nil, errors.New("unknown mqtt access protocol '%s', *protocol")
		}
	}
	config.AddConfig(options)
	return config, nil
}
