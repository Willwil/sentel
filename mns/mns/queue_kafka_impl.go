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

import "github.com/cloustone/sentel/pkg/config"

type kafkaQueue struct {
	config    config.Config
	attribute QueueAttribute
}

func newKafkaQueue(c config.Config, attr QueueAttribute, adaptor Adaptor) (Queue, error) {
	return &kafkaQueue{}, nil
}

func (q *kafkaQueue) SendMessage(QueueMessage) error {
	return nil
}

func (q *kafkaQueue) BatchSendMessage([]QueueMessage) error {
	return nil
}

func (q *kafkaQueue) ReceiveMessage(waitSeconds int) (QueueMessage, error) {
	return QueueMessage{}, nil
}

func (q *kafkaQueue) BatchReceiveMessages(wailtSeconds int, numOfMessages int) ([]QueueMessage, error) {
	return []QueueMessage{}, nil
}

func (q *kafkaQueue) DeleteMessage(handle string) error {
	return nil
}

func (q *kafkaQueue) BatchDeleteMessages(handles []string) error {
	return nil
}

func (q *kafkaQueue) PeekMessage(waitSeconds int) (QueueMessage, error) {
	return QueueMessage{}, nil
}

func (q *kafkaQueue) BatchPeekMessages(wailtSeconds int, numOfMessages int) ([]QueueMessage, error) {
	return []QueueMessage{}, nil
}

func (q *kafkaQueue) SetMessageVisibility(handle string, seconds int) error {
	return nil
}

func (q *kafkaQueue) Destroy() {
}
