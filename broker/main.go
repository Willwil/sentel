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

	"github.com/cloustone/sentel/broker/auth"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/http"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/mqtt"
	"github.com/cloustone/sentel/broker/queue"
	"github.com/cloustone/sentel/broker/quto"
	"github.com/cloustone/sentel/broker/rpc"
	"github.com/cloustone/sentel/broker/sessionmgr"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
)

var (
	configFile = flag.String("c", "/etc/sentel/broker.conf", "config file")
	tenant     = flag.String("t", "", "tenant id")
	product    = flag.String("p", "", "product id")
	protocol   = flag.String("P", "tcp", "mqtt access protocol, tcp|tls|ws|https")
	listen     = flag.String("l", "localhost:1883", "mqtt broker listen address")
)

func main() {
	flag.Parse()

	glog.Infof("Starting 'broker' server...")
	registerService(event.ServiceName, event.New)
	registerService(queue.ServiceName, queue.New)
	registerService(sessionmgr.ServiceName, sessionmgr.New)
	registerService(auth.ServiceName, auth.New)
	registerService(rpc.ServiceName, rpc.New)
	registerService(metric.ServiceName, metric.New)
	registerService(metadata.ServiceName, metadata.New)
	registerService(quto.ServiceName, quto.New)
	registerService(mqtt.ServiceName, mqtt.New)
	registerService(http.ServiceName, http.New)

	// Create global configuration
	if config, err := createConfig(*configFile); err != nil {
		fmt.Println(err.Error())
		flag.PrintDefaults()
		return
	} else {
		// Create service manager according to the configuration
		broker, err := NewBroker(config)
		if err != nil {
			glog.Fatal(err)
		}
		glog.Fatal(broker.Run())
	}
}

func createConfig(fileName string) (core.Config, error) {
	core.RegisterConfigGroup(defaultConfigs)
	options := map[string]map[string]string{}
	options["broker"] = map[string]string{}
	options["mqtt"] = map[string]string{}

	// tenant and product must be set
	if *tenant == "" || *product == "" {
		*tenant = os.Getenv("BROKER_TENANT")
		*product = os.Getenv("BROKER_PRODUCT")
		if *tenant == "" || *product == "" {
			glog.Errorf("tenant:%s, product:%s", *tenant, *product)
			return nil, errors.New("teant and product must be specified for broker")
		}
	}
	options["broker"]["tenant"] = *tenant
	options["broker"]["product"] = *product

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
	config, _ := core.NewConfigWithFile(fileName)
	config.AddConfigs(options)
	return config, nil
}
