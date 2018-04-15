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

package service

import (
	"net"
	"os"
	"strings"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

func UpdateServiceConfigs(c config.Config, services ...string) {
	for _, serviceName := range services {
		// in docker environment, check wether environment is set
		servicePort := strings.ToUpper(serviceName) + "_PORT"
		if v := os.Getenv(servicePort); v != "" {
			c.AddConfigItem(serviceName, v)
		} else {
			// check /etc/hosts files
			if _, err := net.LookupHost(serviceName); err == nil {
				c.AddConfigItem(serviceName, serviceName)
			}
		}
		v, _ := c.String(serviceName)
		glog.Infof("##### updating service '%s = %s'\n", serviceName, v)
	}
}
