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
	"errors"
	"iothub/util/config"

	"github.com/golang/glog"
)

type localDatabase struct {
	config   config.Config
	sessions map[string]Session
}

// Open local database
func (l *localDatabase) Open() error {
	glog.Info("local database Open")

	return nil
}

// Close local database
func (l *localDatabase) Close() {

}

// Backup serialize local database
func (l *localDatabase) Backup(shutdown bool) error {
	return nil
}

// Restore recover data from serialization
func (l *localDatabase) Restore() error {
	return nil
}

// FindSession find session by id
func (l *localDatabase) FindSession(c Context, id string) (*Session, error) {
	v, ok := l.sessions[id]
	if !ok {
		return nil, errors.New("Session id does not exist")
	}

	return &v, nil
}

// DeleteSession delete session by id
func (l *localDatabase) DeleteSession(c Context, id string) error {
	_, ok := l.sessions[id]
	if !ok {
		return errors.New("Session id does not exist")
	}

	delete(l.sessions, id)

	return nil
}

// UpdateSession update session
func (l *localDatabase) UpdateSession(c Context, s *Session) error {
	_, ok := l.sessions[s.Id]
	if !ok {
		return errors.New("Session id does not exist")
	}

	l.sessions[s.Id] = *s
	return nil
}

// RegisterSession register new session
func (l *localDatabase) RegisterSession(c Context, s Session) error {
	if _, ok := l.sessions[s.Id]; ok {
		return errors.New("Session id already exists")
	} else {
		l.sessions[s.Id] = s
		return nil
	}
}

// Device
func (l *localDatabase) AddDevice(c Context, d Device) error {
	return nil
}

func (l *localDatabase) DeleteDevice(c Context, id string) error {
	return nil
}

func (l *localDatabase) UpdateDevice(c Context, d Device) error {
	return nil
}

func (l *localDatabase) GetDeviceState(c Context, id string) (int, error) {
	return 0, nil
}

func (l *localDatabase) SetDeviceState(c Context, state int) error {
	return nil
}

// Topic
func (l *localDatabase) TopicExist(c Context, t Topic) (bool, error) {
	return false, nil
}

func (l *localDatabase) AddTopic(c Context, t Topic) error {
	return nil
}

func (l *localDatabase) DeleteTopic(c Context, id string) error {
	return nil
}

func (l *localDatabase) UpdateTopic(c Context, t Topic) error {
	return nil
}

func (l *localDatabase) AddSubscriber(c Context, t Topic, clientid string) error {
	return nil
}

func (l *localDatabase) RemoveSubscriber(c Context, t Topic, clientid string) error {
	return nil
}

func (l *localDatabase) GetTopicSubscribers(c Context, t Topic) ([]string, error) {
	return nil, nil
}

// Subscription
func (l *localDatabase) AddSubscription(c Context, sub string, qos uint8) error {
	return nil
}

func (l *localDatabase) RetainSubscription(c Context, sub string, qos uint8) error {
	return nil
}

// Message Management
func (l *localDatabase) FindMessage(clientid string, mid uint16) (bool, error) {
	return false, nil
}

func (l *localDatabase) StoreMessage(clientid string, msg Message) error {
	return nil
}

func (l *localDatabase) DeleteMessageWithValidator(clientid string, validator func(msg Message) bool) {

}

func (l *localDatabase) DeleteMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (l *localDatabase) QueueMessage(clientid string, msg Message) error {
	return nil
}

func (l *localDatabase) GetMessageTotalCount(clientid string) int {
	return 0
}

func (l *localDatabase) InsertMessage(clientid string, mid uint16, direction MessageDirection, msg Message) error {
	return nil
}

func (l *localDatabase) ReleaseMessage(clientid string, mid uint16, direction MessageDirection) error {
	return nil
}

func (l *localDatabase) UpdateMessage(clientid string, mid uint16, direction MessageDirection, state MessageState) {

}

// localDatabaseFactory
type localDatabaseFactory struct{}

func (l *localDatabaseFactory) New(c config.Config) (Database, error) {
	d := &localDatabase{
		config:   c,
		sessions: make(map[string]Session),
	}
	return d, nil
}
