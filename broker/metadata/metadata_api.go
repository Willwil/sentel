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
	meta := base.GetService(ServiceName).(*MetadataService)
	return meta.getShadowDeviceStatus(clientId)
}

// SyncShadowDeviceStatus synchronize shadow device's status
func SyncShadowDeviceStatus(clientId string, d *Device) error {
	meta := base.GetService(ServiceName).(*MetadataService)
	return meta.syncShadowDeviceStatus(clientId, d)
}

func Backup(shutdown bool) error {
	return nil
}
func Restore() error {
	return nil
}

// Session
func FindSession(id string) (*Session, error) {
	return nil, nil
}

func DeleteSession(id string) error {
	return nil
}

func UpdateSession(s *Session) error {
	return nil
}
func RegisterSession(s *Session) error {
	return nil
}

// Subscription
func AddSubscription(sessionid string, topic string, qos uint8) error {
	return nil
}

func RetainSubscription(sessionid string, topic string, qos uint8) error {
	return nil
}

func RemoveSubscription(sessionid string, topic string) error {
	return nil
}

// Message Management
func FindMessage(clientid string, mid uint16) bool {
	return false
}
func StoreMessage(clientid string, msg Message) error {
	return nil
}
func DeleteMessageWithValidator(clientid string, validator func(Message) bool) {
}
func DeleteMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func QueueMessage(clientid string, msg Message) error {
	return nil
}
func GetMessageTotalCount(clientid string) int {
	return 0
}
func InsertMessage(clientid string, mid uint16, direction MessageDirection, msg Message) error {
	return nil
}
func ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}
func UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState) {
}

// Client
func GetClients(proto string) []*ClientInfo         { return nil }
func GetClient(proto string, id string) *ClientInfo { return nil }

// Session Info
func GetSessions(proto string, conditions map[string]bool) []*Session { return nil }
func GetSession(proto string, id string) *Session                     { return nil }

// Route Info
func GetRoutes(proto string) []*RouteInfo            { return nil }
func GetRoute(proto string, topic string) *RouteInfo { return nil }

// Topic info
func GetTopics(proto string) []*TopicInfo         { return nil }
func GetTopic(proto string, id string) *TopicInfo { return nil }

// SubscriptionInfo
func GetSubscriptions(proto string) []*SubscriptionInfo         { return nil }
func GetSubscription(proto string, id string) *SubscriptionInfo { return nil }
