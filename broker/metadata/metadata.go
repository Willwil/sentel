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

import "time"

// Session
type Session struct {
	Id                 string
	Username           string
	Password           string
	Keepalive          uint16
	LastMid            uint16
	State              uint8
	CreatedAt          time.Time
	LastMessageInTime  time.Time
	LastMessageOutTime time.Time
	Ping               time.Time
	CleanSession       uint8
	SubscribeCount     uint32
	Protocol           uint8
	RefCount           uint8
}

type Device struct {
	DeviceId    string    `json:"deviceId" bson:"deviceId"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
	Data        []byte    `json:"data" bson:"data"`
}

type Topic struct {
	Name string
}

type MessageDirection int
type MessageState int

const (
	MessageDirectionIn  MessageDirection = 0
	MessageDirectionOut MessageDirection = 1
)

type Message struct {
	ID        uint
	SourceID  string
	Topic     string
	Direction MessageDirection
	State     MessageState
	Qos       uint8
	Retain    bool
	Payload   []uint8
}
