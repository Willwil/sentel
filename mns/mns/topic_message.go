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

type TopicMessage struct {
	TopicName      string                 `json:"-" bson:"TopicName"`
	MessageId      string                 `json:"message_id" bson:"MessageId"`
	MessageBody    string                 `json:"message_body" bson:"MessageBody"`
	MessageBodyMd5 string                 `json:"message_body_md5" bson:"MessageBodyMD5"`
	MessageTag     string                 `json:"message_tag,omitempty" bson:"MessageTag,omitempty"`
	Attributes     map[string]interface{} `json:"attributes,omitempty" bson:"Attributes,omitempty"`
	PublishTime    time.Time              `json:"publish_time",omitempty" bson:"PublishTime,omitempty"`
}

func (t TopicMessage) Topic() string        { return t.TopicName }
func (t TopicMessage) SetTopic(name string) { t.TopicName = name }
func (t TopicMessage) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(t, opt)
}

func (t TopicMessage) Deserialize(buf []byte, opt message.SerializeOption) error {
	return message.Deserialize(buf, opt, &t)
}
