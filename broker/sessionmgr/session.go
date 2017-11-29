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

type Session interface {
	// Id return identifier of session
	Id() string

	// BrokerId return session's client id
	BrokerId() string
	// IsPersistent return wether the session is persistent
	IsPersistent() bool
	// IsValidd return true if the session state valid
	IsValid() bool

	// Info return session's information
	Info() *SessionInfo
}

// SessionInfo
type SessionInfo struct {
	ClientId           string
	CleanSession       uint8
	MessageMaxInflight uint64
	MessageInflight    uint64
	MessageInQueue     uint64
	MessageDropped     uint64
	AwaitingRel        uint64
	AwaitingComp       uint64
	AwaitingAck        uint64
	CreatedAt          string
}

type Message struct {
	Topic     string
	Id        uint16
	Direction uint8
	State     uint8
	Qos       uint8
	Retain    bool
	Payload   []uint8
}

// ClientInfo
type ClientInfo struct {
	UserName     string
	CleanSession bool
	PeerName     string
	ConnectTime  string
}

// TopicInfo
type TopicInfo struct {
	Topic     string
	Attribute string
}

// SubscriptionInfo
type SubscriptionInfo struct {
	ClientId  string
	Topic     string
	Attribute string
}
