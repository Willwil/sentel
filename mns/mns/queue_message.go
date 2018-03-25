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

	"github.com/cloustone/sentel/pkg/message"
)

type QueueMessage struct {
	TopicName        string                 `json:"-" bson:"TopicName"`
	MessageId        string                 `json:"message_id" bson:"MessageId"`
	ReceiptHandle    string                 `json:"receipt_handle" bson:"ReceiptHandle"`
	MessageBody      []byte                 `json:"message_body" bson:"MessageBody"`
	MessageBodyMd5   string                 `json:"message_body_md5" bson:"MessageBodyMD5"`
	MessageTag       string                 `json:"message_tag,omitempty" bson:"MessageTag,omitempty"`
	Priority         uint8                  `json:"priority,omitempty" bson:"Priority, omitempty"`
	FirstEnqueueTime time.Time              `json:"first_enqueue_time,omitempty" bson:"FirstEnqueueTime,omitempty"`
	Attributes       map[string]interface{} `json:"attributes,omitempty" bson:"Attributes,omitempty"`
}

func (m QueueMessage) Topic() string        { return m.TopicName }
func (m QueueMessage) SetTopic(name string) { m.TopicName = name }
func (m QueueMessage) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(m, opt)
}

func (m QueueMessage) Deserialize(buf []byte, opt message.SerializeOption) error {
	return message.Deserialize(buf, opt, &m)
}
