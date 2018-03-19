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

package mns

type QueueAttribute struct {
	DelaySeconds           uint32 `json:"delaySeconds" bson:"DelaySeconds"`
	MaximumMessageSize     uint32 `json:"maximumMessageSize" bson:"MaximumMessageSize"`
	MessageRetentionPeriod uint32 `json:"messageRetentionPeriod" bson:"MessageRetentionPeriod"`
	VisibilityTimeout      uint32 `json:"visibilityTimeout" bson:"VisibilityTimeout"`
	PollingWaitSeconds     uint8  `json:"pollingWaitSeconds" bson:"PollingWaitSeconds"`
	LoggingEnabled         bool   `json:"loggingEnabled" bson:"LoggingEnabled"`
}

type Queue interface {
	Name() string
	SetAttribute(QueueAttribute)
	GetAttribute() QueueAttribute
	SendMessage(Message) error
	BatchSendMessage([]Message) error
	ReceiveMessage(waitSeconds int) (Message, error)
	BatchReceiveMessages(wailtSeconds int, numOfMessages int) ([]Message, error)
	DeleteMessage(handle string) error
	PeekMessage(waitSeconds int) (Message, error)
	BatchPeekMessages(wailtSeconds int, numOfMessages int) ([]Message, error)
	SetMessageVisibility(handle string, seconds int)
	Destroy()
}
