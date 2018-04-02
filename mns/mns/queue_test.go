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
	"testing"

	"github.com/cloustone/sentel/pkg/config"
)

var (
	queueAttribute = QueueAttribute{
		AccountId: "cloustone",
		QueueName: "test",
	}
	defaultQueueConfigs = config.M{
		"test": {
			"mongo": "localhost:27817",
			"kafka": "localhost:9092",
		},
	}
	queueConfig  config.Config
	defaultQueue Queue
)

func newQueue(t *testing.T) Queue {
	queueConfigs := config.New("test")
	queueConfigs.AddConfig(defaultQueueConfigs)
	queue, err := NewQueue(queueConfigs, queueAttribute)
	if err != nil {
		t.Error(err)
		panic("nil queue")
	}
	return queue
}

func TestQueue_SendMessage(t *testing.T) {
	msg := QueueMessage{
		MessageBody: []byte("hello, world"),
		MessageTag:  "test",
	}

	queue := newQueue(t)
	queue.Destroy()
	if err := queue.SendMessage(msg); err != nil {
		t.Error(err)
	}
}

func TestQueue_BatchSendMessage(t *testing.T) {
	msgs := []QueueMessage{
		{
			MessageBody: []byte("hello, world1"),
			MessageTag:  "test",
		},
		{
			MessageBody: []byte("hello, world2"),
			MessageTag:  "test",
		},
		{
			MessageBody: []byte("hello, world3"),
			MessageTag:  "test",
		},
	}
	queue := newQueue(t)
	defer queue.Destroy()
	if err := queue.BatchSendMessage(msgs); err != nil {
		t.Error(err)
	}
}

func TestQueue_ReceiveMessage(t *testing.T) {
	queue := newQueue(t)
	defer queue.Destroy()
	if _, err := queue.ReceiveMessage(3); err != nil {
		t.Error(err)
	}
}

func TestQueue_BatchReceiveMessage(t *testing.T) {
	queue := newQueue(t)
	defer queue.Destroy()
	if _, err := queue.BatchReceiveMessages(3, 2); err != nil {
		t.Error(err)
	}
}

func TestQueue_DeleteMessage(t *testing.T) {
}

func TestQueue_BatchDeleteMessage(t *testing.T) {
}

func TestQueue_PeekMessage(t *testing.T) {
}

func TestQueue_BatchPeekMessage(t *testing.T) {
}

func TestQueue_SetMessageVisibility(t *testing.T) {
}

func TestQueue_Destroy(t *testing.T) {
}
