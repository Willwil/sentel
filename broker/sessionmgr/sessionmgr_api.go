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

package sessionmgr

import "github.com/cloustone/sentel/broker/base"

// FindSesison return session object by clientId if existed
func FindSession(clientId string) (*Session, error) {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.findSession(clientId)
}

// DeleteSession remove session specified by clientId from subdata
func DeleteSession(clientId string) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.deleteSession(clientId)
}

// RegiserSession register session into subdata
func RegisterSession(s *Session) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.registerSession(s)
}

// AddSubscription add a subscription into subdat
func AddSubscription(clientId string, topic string, qos uint8) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.addSubscription(clientId, topic, qos, nil)
}

// RetainSubscription retain the client with topic
func RetainSubscription(clientId string, topic string, qos uint8) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.retainSubscription(clientId, topic, qos)
}

// RmoeveSubscription remove specified topic from subdata
func RemoveSubscription(clientId string, topic string) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.removeSubscription(clientId, topic)
}

// DeleteMessageWithValidator delete message in subdata with confition
func DeleteMessageWithValidator(clientId string, validator func(Message) bool) {
	sub := base.GetService(ServiceName).(*sessionManager)
	sub.deleteMessageWithValidator(clientId, validator)
}

// DeleteMessge delete message specified by idfrom subdata
func DeleteMessage(clientId string, mid uint16, direction MessageDirection) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.deleteMessage(clientId, mid, direction)
}

// QueueMessage save message into subdata
func QueueMessage(clientId string, msg *Message) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.queueMessage(clientId, msg)
}

// InsertMessage insert specified message into subdata
func InsertMessage(clientId string, mid uint16, direction MessageDirection, msg *Message) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.insertMessage(clientId, mid, direction, msg)
}

// ReleaseMessage release message from subdata
func ReleaseMessage(clientId string, mid uint16, direction MessageDirection) error {
	sub := base.GetService(ServiceName).(*sessionManager)
	return sub.releaseMessage(clientId, mid, direction)
}

// Client
func GetClients(proto string) []*ClientInfo         { return nil }
func GetClient(proto string, id string) *ClientInfo { return nil }

// Session Info
func GetSessions(proto string, conditions map[string]bool) []*Session { return nil }
func GetSession(proto string, id string) *Session                     { return nil }

// Topic info
func GetTopics(proto string) []*TopicInfo         { return nil }
func GetTopic(proto string, id string) *TopicInfo { return nil }

// SubscriptionInfo
func GetSubscriptions(proto string) []*SubscriptionInfo         { return nil }
func GetSubscription(proto string, id string) *SubscriptionInfo { return nil }
