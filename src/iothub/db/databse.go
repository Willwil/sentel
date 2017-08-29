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

package db

import (
	"fmt"
	"iothub/util/config"
	"time"

	"github.com/golang/glog"
)

// Session database
type Session struct {
	Id                 string
	Username           string
	Password           string
	Keepalive          uint16
	LastMid            uint16
	State              uint16
	LastMessageInTime  time.Time
	LastMessageOutTime time.Time
	Ping               time.Time
	CleanSession       uint16
	SubscribeCount     uint32
	Protocol           uint8
	RefCount           uint8
}

type MessageDirection int

const (
	MessageDirectionIn  MessageDirection = 0
	MessageDirectionOut MessageDirection = 1
)

type Device struct{}
type Topic struct{}
type MessageState int
type Message struct {
	Id        uint
	Direction MessageDirection
	State     MessageState
	Qos       uint8
	Retain    bool
	Payload   []uint8
}

type Context interface{}

type Database interface {
	Open() error
	Close()
	Backup(shutdown bool) error
	Restore() error

	// Session
	FindSession(c Context, id string) (*Session, error)
	DeleteSession(c Context, id string) error
	UpdateSession(c Context, s *Session) error

	// Device
	AddDevice(c Context, d Device) error
	DeleteDevice(c Context, id string) error
	UpdateDevice(c Context, d Device) error
	GetDeviceState(c Context, id string) (int, error)
	SetDeviceState(c Context, state int) error

	// Topic
	TopicExist(c Context, t Topic) (bool, error)
	AddTopic(c Context, t Topic) error
	DeleteTopic(c Context, id string) error
	UpdateTopic(c Context, t Topic) error
	AddTopicSubscriber(c Context, t Topic, clientid string) error
	RemoveTopicSubscriber(c Context, t Topic, clientid string) error
	GetTopicSubscribers(c Context, t Topic) ([]string, error)

	// Message Management
	FindMessage(clientid string, mid uint) (bool, error)
	StoreMessage(clientid string, msg Message) error
	QueueMessage(clientid string, msg Message) error
	GetMessageTotalCount(clientid string) int
	DeleteMessage(clientid string, mid int, direction MessageDirection) error
	InsertMessage(clientid string, mid int, direction MessageDirection, msg Message) error
	ReleaseMessage(clientid string, mid int, direction MessageDirection) error
	UpdateMessage(clientid string, mid int, direction MessageDirection, state MessageState)
}

type databaseFactory interface {
	New(c config.Config) (Database, error)
}

var _allDatabase = make(map[string]databaseFactory)

func registerDatabase(name string, d databaseFactory) {
	if _allDatabase[name] != nil {
		glog.Fatalf("Database %s already registered", name)
		return
	}
	_allDatabase[name] = d
}

func NewDatabase(c config.Config) (Database, error) {
	repo, err := c.String("database", "repository")
	if err != nil {
		glog.Error("Database configuration has no repository")
		return nil, err
	}
	if _allDatabase[repo] == nil {
		return nil, fmt.Errorf("Database %s is not registered", repo)
	}
	return _allDatabase[repo].New(c)
}