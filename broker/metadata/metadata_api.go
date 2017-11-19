//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package metadata

import "github.com/cloustone/sentel/broker/base"

// GetShadowDeviceStatus return shadow device's status
func GetShadowDeviceStatus(clientId string) (*Device, error) {
	broker := base.GetBroker()
	meta := broker.GetServicesByName(MetadataServiceName)[0].(*MetadataService)
	return meta.getShadowDeviceStatus(clientId)
}

// SyncShadowDeviceStatus synchronize shadow device's status
func SyncShadowDeviceStatus(clientId string, d *Device) error {
	broker := base.GetBroker()
	meta := broker.GetServicesByName(MetadataServiceName)[0].(*MetadataService)
	return meta.syncShadowDeviceStatus(clientId, d)
}
