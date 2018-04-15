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

	// "net/http"
	// _ "net/http/pprof"

	"github.com/cloustone/sentel/broker/auth"
	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/commands"
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
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"
)

var (
	configFile = flag.String("c", "/etc/sentel/broker.conf", "config file")
	tenant     = flag.String("t", "", "tenant id")
	protocol   = flag.String("P", "tcp", "mqtt access protocol, tcp|tls|ws|https")
	listen     = flag.String("l", "localhost:1883", "mqtt broker listen address")
	daemon     = flag.Bool("d", false, "broker daemon mode")
)

func main() {
	flag.Parse()

	// go func() {
	// 	log.Println(http.ListenAndServe("localhost:6060", nil))
	// }()

	config, err := createConfig(*configFile, *daemon)
	if err != nil {
		fmt.Println(err.Error())
		flag.PrintDefaults()
		return
	}
	if *daemon == true {
		glog.Infof("Starting 'broker' server...")
		glog.Fatal(runBrokerDaemon(config))
	} else {
		runBrokerClient(config)
	}
}

func createConfig(fileName string, daemon bool) (config.Config, error) {
	config := config.New("broker")
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	service.UpdateServiceConfigs(config, "mongo", "kafka")
	if daemon {
		// tenant and product must be set
		if *tenant == "" {
			*tenant = os.Getenv("BROKER_TENANT")
			if *tenant == "" {
				return nil, errors.New("teant must be specified for broker")
			}
		}
		config.AddConfigItem("tenant", *tenant)

		// Mqtt protocol
		if *protocol != "" {
			switch *protocol {
			case "tcp", "ws", "tls", "https":
				config.AddConfigItem("protocol", *protocol)
				if *listen != "" {
					config.AddConfigItemWithSection("mqtt", *protocol, *listen)
				}
			default:
				return nil, errors.New("unknown mqtt access protocol '%s', *protocol")
			}
		}
	}
	return config, nil
}

func runBrokerClient(c config.Config) {
	if err := commands.Run(c); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func runBrokerDaemon(c config.Config) error {
	base.SetBrokerStartupInfo(
		&base.BrokerStartupInfo{
			Id:     uuid.NewV4().String(),
			Tenant: c.MustString("tenant"),
		})
	mgr, _ := service.NewServiceManager("broker", c)
	mgr.AddService(event.ServiceFactory{})
	mgr.AddService(queue.ServiceFactory{})
	mgr.AddService(sessionmgr.ServiceFactory{})
	mgr.AddService(auth.ServiceFactory{})
	mgr.AddService(rpc.ServiceFactory{})
	mgr.AddService(metric.ServiceFactory{})
	mgr.AddService(metadata.ServiceFactory{})
	mgr.AddService(quto.ServiceFactory{})
	mgr.AddService(mqtt.ServiceFactory{})
	mgr.AddService(http.ServiceFactory{})
	return mgr.RunAndWait()
}
