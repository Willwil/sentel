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

// FindSesison return session object by clientID if existed
func FindSession(clientID string) (Session, error) {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.findSession(clientID)
}

// RegiserSession register session into session manager
func RegisterSession(s Session) error {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.registerSession(s)
}

// RemoveSession remove session from session manager
func RemoveSession(clientID string) error {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.removeSession(clientID)
}

// DeleteMessageWithValidator delete message in session manager with condition
func DeleteMessageWithValidator(clientID string, validator func(*base.Message) bool) {
	sm := base.GetService(ServiceName).(*sessionManager)
	sm.deleteMessageWithValidator(clientID, validator)
}

// DeleteMessge delete message specified by id from session manager
func DeleteMessage(clientID string, pid uint16, direction uint8) {
	sm := base.GetService(ServiceName).(*sessionManager)
	sm.deleteMessage(clientID, pid, direction)
}

// FindMessage return message already existed in session manager
func FindMessage(clientID string, pid uint16, dir uint8) *base.Message {
	return nil
}

// GetTopics return all topic in the broker
func GetTopics() []*Topic {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getTopics()
}

// GetClientTopicss return client's subscribed topics
func GetClientTopics(clientID string) []*Topic {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getClientTopics(clientID)
}

// GetTopic return specified topic subscription info
func GetTopicSubscription(topic string) []*Subscription {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getTopicSubscription(topic)
}

// GetSubscriptions return all subscriptions in the broker
func GetSubscriptions() []*Subscription {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getSubscriptions()
}

// GetClientSubscriptions return client's subscription
func GetClientSubscriptions(clientID string) []*Subscription {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getClientSubscriptions(clientID)
}

//  GetSession return all sessions in the broker
func GetSessions() []*Session {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getSessions()
}

// GetSessionInfo return client's session detail information
func GetSessionInfo(clientID string) *Session {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getSessionInfo(clientID)
}

// GetClients return all clients in broker cluster
func GetClients() []*Client {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getClients()
}

// GetClient return client's detail information
func GetClient(clientID string) *Client {
	sm := base.GetService(ServiceName).(*sessionManager)
	return sm.getClient(clientID)
}
