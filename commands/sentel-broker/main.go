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

	"github.com/cloustone/sentel/broker"
	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "stentel-broker.conf", "config file")
	tenant             = flag.String("t", "", "tenant id")
	product            = flag.String("p", "", "product id")
)

func main() {
	flag.Parse()
	// tenant and product must be set
	if tenant == nil || product == nil || *tenant == "" || *product == "" {
		flag.PrintDefaults()
		return
	}
	err := broker.RunWithConfig(*configFileFullPath, map[string]string{"tenant": *tenant, "product": *product})
	glog.Fatal(err)
}
