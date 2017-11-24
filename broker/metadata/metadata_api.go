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

import "github.com/cloustone/sentel/broker/broker"

// GetShadowDeviceStatus return shadow device's status
func GetShadowDeviceStatus(clientId string) (*Device, error) {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.getShadowDeviceStatus(clientId)
}

// SyncShadowDeviceStatus synchronize shadow device's status
func SyncShadowDeviceStatus(clientId string, d *Device) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.syncShadowDeviceStatus(clientId, d)
}

// FindSesison return session object by clientId if existed
func FindSession(clientId string) (*Session, error) {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.findSession(clientId)
}

// DeleteSession remove session specified by clientId from metadata
func DeleteSession(clientId string) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.deleteSession(clientId)
}

// RegiserSession register session into metadata
func RegisterSession(s *Session) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.registerSession(s)
}

// AddSubscription add a subscription into metadat
func AddSubscription(clientId string, topic string, qos uint8) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.addSubscription(clientId, topic, qos)
}

// RetainSubscription retain the client with topic
func RetainSubscription(clientId string, topic string, qos uint8) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.retainSubscription(clientId, topic, qos)
}

// RmoeveSubscription remove specified topic from metadata
func RemoveSubscription(clientId string, topic string) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.removeSubscription(clientId, topic)
}

// DeleteMessageWithValidator delete message in metadata with confition
func DeleteMessageWithValidator(clientId string, validator func(Message) bool) {
	meta := broker.GetService(ServiceName).(*MetadataService)
	meta.deleteMessageWithValidator(clientId, validator)
}

// DeleteMessge delete message specified by idfrom metadata
func DeleteMessage(clientId string, mid uint16, direction MessageDirection) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.deleteMessage(clientId, mid, direction)
}

// QueueMessage save message into metadata
func QueueMessage(clientId string, msg *Message) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.queueMessage(clientId, msg)
}

// InsertMessage insert specified message into metadata
func InsertMessage(clientId string, mid uint16, direction MessageDirection, msg *Message) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.insertMessage(clientId, mid, direction, msg)
}

// ReleaseMessage release message from metadata
func ReleaseMessage(clientId string, mid uint16, direction MessageDirection) error {
	meta := broker.GetService(ServiceName).(*MetadataService)
	return meta.releaseMessage(clientId, mid, direction)
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
