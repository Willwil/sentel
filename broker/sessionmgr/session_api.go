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
func FindSession(clientId string) (Session, error) {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.findSession(clientId)
}

// RegiserSession register session into session manager
func RegisterSession(s Session) error {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.registerSession(s)
}

// RemoveSession remove session from session manager
func RemoveSession(clientId string) error {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.removeSession(clientId)
}

// DeleteMessageWithValidator delete message in session manager with condition
func DeleteMessageWithValidator(clientId string, validator func(*base.Message) bool) {
	sm := base.GetService(ServiceName).(*sessionManager)
	sm.deleteMessageWithValidator(clientId, validator)
}

// DeleteMessge delete message specified by id from session manager
func DeleteMessage(clientId string, mid uint16, direction uint8) {
	sm := base.GetService(ServiceName).(*sessionManager)
	sm.deleteMessage(clientId, mid, direction)
}

// FindMessage return message already existed in session manager
func FindMessage(clientId string, pid uint16, dir uint8) *base.Message {
	return nil
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
