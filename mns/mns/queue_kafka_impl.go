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
	"errors"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

type kafkaQueue struct {
	queueName string
	config    config.Config
	attribute QueueAttribute
}

func newKafkaQueue(c config.Config, attr QueueAttribute, adaptor Adaptor) (Queue, error) {
	return &kafkaQueue{}, nil
}

func (q *kafkaQueue) SendMessage(msg QueueMessage) error {
	return message.PostMessage(q.config, msg)
}

func (q *kafkaQueue) BatchSendMessage(msgs []QueueMessage) error {
	kmsgs := []message.Message{}
	for _, msg := range msgs {
		kmsgs = append(kmsgs, msg)
	}
	return message.PostMessages(q.config, kmsgs)
}

func (q *kafkaQueue) CreateMessage(topic string) message.Message {
	switch topic {
	case q.queueName:
		return &QueueMessage{TopicName: q.queueName}
	default:
		return nil
	}
}

func queueMessageHandlerFunc(msg message.Message, ctx interface{}) {

}

func (q *kafkaQueue) ReceiveMessage(ws int) (msg QueueMessage, err error) {
	if consumer, err := message.NewConsumer(q.config, "kafkaQueue"); err == nil {
		defer consumer.Close()
		consumer.SetMessageFactory(q)
		msgChan := make(chan QueueMessage)
		consumer.Subscribe(q.queueName, queueMessageHandlerFunc, msgChan)
		consumer.Start()
		timeout := time.NewTimer(time.Second * time.Duration(ws))
		select {
		case <-timeout.C:
			err = errors.New("kafaka queue receive message timeout")
		case v := <-msgChan:
			msg = v
		}
	}
	return
}

type batchMessageContext struct {
	msgChan       chan []QueueMessage
	nbrofMessages int
}

func batchQueueMessageHandlerFunc(msg message.Message, ctx interface{}) {

}

func (q *kafkaQueue) BatchReceiveMessages(ws int, numOfMessages int) (msgs []QueueMessage, err error) {
	if consumer, err := message.NewConsumer(q.config, "kafkaQueue"); err == nil {
		defer consumer.Close()
		consumer.SetMessageFactory(q)
		ctx := batchMessageContext{
			msgChan:       make(chan []QueueMessage),
			nbrofMessages: numOfMessages,
		}
		consumer.Subscribe(q.queueName, queueMessageHandlerFunc, ctx)
		consumer.Start()
		timeout := time.NewTimer(time.Second * time.Duration(ws))
		select {
		case <-timeout.C:
			err = errors.New("kafaka queue receive message timeout")
		case v := <-ctx.msgChan:
			msgs = v
		}
	}
	return
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
