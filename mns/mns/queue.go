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
	QueueName              string `bson:"QueueName,omitempty" json:"queue_name,omitempty"`
	DelaySeconds           int32  `bson:"DelaySenconds,omitempty" json:"delay_senconds,omitempty"`
	MaxMessageSize         int32  `bson:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32  `bson:"MessageRetentionPeriod,omitempty" json:"message_retention_period,omitempty"`
	VisibilityTimeout      int32  `bson:"VisibilityTimeout,omitempty" json:"visibility_timeout,omitempty"`
	PollingWaitSeconds     int32  `bson:"PollingWaitSeconds,omitempty" json:"polling_wait_secods,omitempty"`
	ActiveMessages         int64  `bson:"ActiveMessages,omitempty" json:"active_messages,omitempty"`
	InactiveMessages       int64  `bson:"InactiveMessages,omitempty" json:"inactive_messages,omitempty"`
	DelayMessages          int64  `bson:"DelayMessages,omitempty" json:"delay_messages,omitempty"`
	CreateTime             int64  `bson:"CreateTime,omitempty" json:"create_time,omitempty"`
	LastModifyTime         int64  `bson:"LastModifyTime,omitempty" json:"last_modify_time,omitempty"`
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
