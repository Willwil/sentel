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

	"github.com/cloustone/sentel/broker"
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
	options, err := parseCliOptions()
	if err != nil {
		fmt.Println(err.Error())
		flag.PrintDefaults()
		return
	}
	err = broker.RunWithConfig(*configFile, options)
	glog.Fatal(err)
}

func parseCliOptions() (map[string]map[string]string, error) {
	cliOptions := map[string]map[string]string{}
	// tenant and product must be set
	if *tenant == "" || *product == "" {
		return nil, errors.New("teant and product must be specified for broker")
	}
	cliOptions["broker"] = map[string]string{"tenant": *tenant}
	cliOptions["broker"] = map[string]string{"product": *product}
	cliOptions["mqtt"] = map[string]string{}
	// Mqtt protocol
	if *protocol != "" {
		switch *protocol {
		case "tcp", "ws", "tls", "https":
			cliOptions["broker"]["protocol"] = *protocol
			if *listen != "" {
				cliOptions["mqtt"][*protocol] = *listen
			}
		default:
			return nil, errors.New("unknown mqtt access protocol '%s', *protocol")
		}
	}
	return cliOptions, nil
}
