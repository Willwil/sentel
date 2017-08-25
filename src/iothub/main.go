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
	"iothub/base"
	"iothub/util/config"

	"github.com/golang/glog"
)

var (
	configFileFullPath = flag.String("c", "/etc/sentel/iothub/iothub.conf", "config file")
	logFileFullPath    = flag.String("l", "/var/log/sentel/iothub.log", "log file")
)

func main() {
	flag.Parse()
	// Get configuration
	config, err := config.NewWithConfigFile(*configFileFullPath)
	if err != nil {
		flag.PrintDefaults()
		return
	}

	mgr, err := base.NewServiceManager(config)
	if err != nil {
		glog.Error("Failed to create ServiceManager")
		return
	}
	glog.Error(mgr.Start())
}
