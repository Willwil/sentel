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

package base

import "github.com/cloustone/sentel/core"

// SessionInfo
type SessionInfo struct {
	ClientId           string
	CleanSession       bool
	MessageMaxInflight uint64
	MessageInflight    uint64
	MessageInQueue     uint64
	MessageDropped     uint64
	AwaitingRel        uint64
	AwaitingComp       uint64
	AwaitingAck        uint64
	CreatedAt          string
}

type Session interface {
	// Identifier get session identifier
	Identifier() string
	// Info return session information
	Info() *SessionInfo
	// GetService get the service ower for current session
	Service() core.Service
	// Handle indicate service to handle the packet
	Handle() error
	// Destroy will release current session
	Destroy() error
	// RegisterObserver register observer on session
	RegisterObserver(SessionObserver)
}

type SessionObserver interface {
	OnGetMountPoint() string
	// Callback when session is connected
	OnConnect(s Session, userdata interface{}) error
	// Callback when session is disconnected
	OnDisconnect(s Session, userdta interface{}) error
	// Callback before session publish topic
	OnPublish(s Session, userdata interface{}) error
	// Callback when message arrived on session
	OnMessage(s Session, userdata interface{}) error
	// Callback when session received subscribe notification
	OnSubscribe(s Session, userdata interface{}) error
	// Callback when session received unsubscribe notification
	OnUnsubscribe(s Session, userdata interface{}) error
}
