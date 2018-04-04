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
	accountID = "cloustone"
	queueName = "test"
	topicName = "test"
	configs   = config.M{
		"test": {
			"mongo": "localhost:27017",
			"kafka": "localhost:9092",
		},
	}
	testmgr MnsManager
)

func Test_NewManager(t *testing.T) {
	c := config.New("test")
	c.AddConfig(configs)
	var err error
	//sarama.Logger = log.New(os.Stdout, "[mns] ", log.LstdFlags)
	testmgr, err = NewManager(c)
	if err != nil {
		panic(err)
	}
}

func Test_CreateQueue(t *testing.T) {
	if _, err := testmgr.CreateQueue(accountID, queueName); err != nil {
		t.Error(err)
	}
}

func Test_GetQueueAttribute(t *testing.T) {
	if _, err := testmgr.GetQueueAttribute(accountID, queueName); err != nil {
		t.Error(err)
	}
}

func Test_SetQueueAttribute(t *testing.T) {
	attr := QueueAttribute{}
	if err := testmgr.SetQueueAttribute(accountID, queueName, attr); err != nil {
		t.Error(err)
	}
}

func Test_GetQueues(t *testing.T) {
	queues := testmgr.GetQueues(accountID)
	if len(queues) != 1 {
		t.Errorf("queues length is %d", len(queues))
	}
}

func Test_SendQueueMessage(t *testing.T) {
	msg := QueueMessage{
		MessageBody: []byte("hello, world"),
	}
	if err := testmgr.SendQueueMessage(accountID, queueName, msg); err != nil {
		t.Error(err)
	}
}

func Test_ReceiveQueueMessage(t *testing.T) {
	if msg, err := testmgr.ReceiveQueueMessage(accountID, queueName, 3); err != nil {
		t.Error(err)
	} else {
		if string(msg.MessageBody) != "hello, world" {
			t.Errorf("receive message failed, Message Body is: %s", msg.MessageBody)
		}
	}
}

func Test_BatchSendQueueMessages(t *testing.T) {
	msgs := []QueueMessage{
		{
			MessageBody: []byte("hello, world1"),
		},
		{
			MessageBody: []byte("hello, world2"),
		},
	}
	if err := testmgr.BatchSendQueueMessages(accountID, queueName, msgs); err != nil {
		t.Error(err)
	}
}

func Test_BatchReceiveMessages(t *testing.T) {
	if msgs, err := testmgr.BatchReceiveQueueMessages(accountID, queueName, 3, 2); err != nil {
		t.Error(err)
	} else {
		if len(msgs) != 2 {
			t.Errorf("message count is wrong: %d", len(msgs))
		}
	}
}

func Test_CreateTopic(t *testing.T) {
	if _, err := testmgr.CreateTopic(accountID, topicName); err != nil {
		t.Error(err)
	}
}

func Test_GetTopic(t *testing.T) {
	if _, err := testmgr.GetTopic(accountID, topicName); err != nil {
		t.Error(err)
	}
}

func Test_UpdateTopic(t *testing.T) {
	topic := Topic{
		TopicName: topicName,
		AccountId: accountID,
	}
	if err := testmgr.UpdateTopic(accountID, topicName, topic); err != nil {
		t.Error(err)
	}
}

func Test_ListTopics(t *testing.T) {
	topics := testmgr.ListTopics(accountID)
	if len(topics) != 1 {
		t.Error("topics count is wrong")
	}
}

func Test_DeleteTopic(t *testing.T) {
	if err := testmgr.DeleteTopic(accountID, topicName); err != nil {
		t.Error(err)
	}
}

func Test_DeleteQueue(t *testing.T) {
	if err := testmgr.DeleteQueue(accountID, queueName); err != nil {
		t.Error(err)
	}
}
