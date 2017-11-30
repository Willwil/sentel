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

package base

import "time"

type Message struct {
	Topic     string    `json:"topic" bson:"topic"`
	PacketId  uint16    `json:"packetId" bson:"packetId"`
	Direction uint8     `json:"diretion" bson:"direction"`
	State     uint8     `json:"state" bson:"state"`
	Qos       uint8     `json:"qos" bson:"qos"`
	Dup       uint8     `json:"dup" bson:"dup"`
	Retain    bool      `json:"retain" bson:"retain"`
	Payload   []uint8   `json:"payload" bson:"payload"`
	Time      time.Time `json:"time" bson:"time"`
}
