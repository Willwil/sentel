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

import (
	"time"

	"github.com/cloustone/sentel/pkg/config"
)

const FormatOfQueueName = "mns-queues-%s-%s"

type QueueAttribute struct {
	AccountId              string    `bson:"AccountId" json:"account_id"`
	QueueName              string    `bson:"QueueName,omitempty" json:"queue_name,omitempty"`
	DelaySeconds           int32     `bson:"DelaySenconds,omitempty" json:"delay_senconds,omitempty"`
	MaxMessageSize         int32     `bson:"MaximumMessageSize,omitempty" json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32     `bson:"MessageRetentionPeriod,omitempty" json:"message_retention_period,omitempty"`
	VisibilityTimeout      int32     `bson:"VisibilityTimeout,omitempty" json:"visibility_timeout,omitempty"`
	PollingWaitSeconds     int32     `bson:"PollingWaitSeconds,omitempty" json:"polling_wait_secods,omitempty"`
	ActiveMessages         int64     `bson:"ActiveMessages,omitempty" json:"active_messages,omitempty"`
	InactiveMessages       int64     `bson:"InactiveMessages,omitempty" json:"inactive_messages,omitempty"`
	DelayMessages          int64     `bson:"DelayMessages,omitempty" json:"delay_messages,omitempty"`
	CreateTime             time.Time `bson:"CreateTime,omitempty" json:"create_time,omitempty"`
	LastModifyTime         time.Time `bson:"LastModifyTime,omitempty" json:"last_modify_time,omitempty"`
}

type Queue interface {
	SendMessage(QueueMessage) error
	BatchSendMessages([]QueueMessage) error
	ReceiveMessage(ws int) (QueueMessage, error)
	BatchReceiveMessages(ws int, numOfMessages int) ([]QueueMessage, error)
	DeleteMessage(handle string) error
	BatchDeleteMessages(handles []string) error
	PeekMessage(waitSeconds int) (QueueMessage, error)
	BatchPeekMessages(ws int, numOfMessages int) ([]QueueMessage, error)
	SetMessageVisibility(handle string, seconds int) error
	Destroy()
}

func NewQueue(c config.Config, attr QueueAttribute) (Queue, error) {
	return newKafkaQueue(c, attr)
}
