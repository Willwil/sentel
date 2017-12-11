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

package meter

import (
	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/meter/api"
	"github.com/cloustone/sentel/meter/collector"
)

// RunWithConfigFile create and start meter service
func RunWithConfigFile(fileName string) error {
	return core.RunWithConfigFile("meter", fileName)
}

// init initialize default configurations and services
func init() {
	core.RegisterConfigGroup(defaultConfigs)
	core.RegisterServiceWithConfig("api", &api.ApiServiceFactory{}, api.Configs)
	core.RegisterServiceWithConfig("collector", &collector.CollectorServiceFactory{}, collector.Configs)
}
