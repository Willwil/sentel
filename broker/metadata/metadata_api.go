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
	meta := base.GetService(ServiceName).(*metadataService)
	return meta.getShadowDeviceStatus(clientId)
}

// SyncShadowDeviceStatus synchronize shadow device's status
func SyncShadowDeviceStatus(clientId string, d *Device) error {
	meta := base.GetService(ServiceName).(*metadataService)
	return meta.syncShadowDeviceStatus(clientId, d)
}

// DeleteMessageWithValidator delete message in session manager with condition
func DeleteMessageWithValidator(clientId string, validator func(*base.Message) bool) {
	meta := base.GetService(ServiceName).(*metadataService)
	meta.deleteMessageWithValidator(clientId, validator)
}

// DeleteMessge delete message specified by id from session manager
func DeleteMessage(clientId string, pid uint16, direction uint8) {
	meta := base.GetService(ServiceName).(*metadataService)
	meta.deleteMessage(clientId, pid, direction)
}

// FindMessage return message already existed in session manager
func FindMessage(clientId string, pid uint16, dir uint8) *base.Message {
	return nil
}

// AddMessage save message in session's store queue
func AddMessage(clientId string, msg *base.Message) {
}
