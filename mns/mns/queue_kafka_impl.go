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
	"fmt"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

type kafkaQueue struct {
	queueName string
	config    config.Config
	attribute QueueAttribute
}

func newKafkaQueue(c config.Config, attr QueueAttribute) (Queue, error) {
	return &kafkaQueue{
		queueName: fmt.Sprintf(FormatOfQueueName, attr.AccountId, attr.QueueName),
		config:    c,
		attribute: attr,
	}, nil
}

func (q *kafkaQueue) SendMessage(msg QueueMessage) error {
	msg.TopicName = q.queueName
	return message.SendMessage(q.config, &msg)
}

func (q *kafkaQueue) BatchSendMessages(msgs []QueueMessage) error {
	kmsgs := []message.Message{}
	for _, msg := range msgs {
		msg.TopicName = q.queueName
		kmsgs = append(kmsgs, &msg)
	}
	return message.SendMessages(q.config, kmsgs)
}

func (q *kafkaQueue) CreateMessage(topic string) message.Message {
	switch topic {
	case q.queueName:
		return &QueueMessage{TopicName: q.queueName}
	default:
		return nil
	}
}

func (q *kafkaQueue) ReceiveMessage(ws int) (msg QueueMessage, err error) {
	var consumer message.Consumer
	if consumer, err = message.NewConsumer(q.config, "kafkaQueue"); err == nil {
		consumer.SetMessageFactory(q)
		msgChan := make(chan *QueueMessage)
		consumer.Subscribe(q.queueName,
			func(msg message.Message, ctx interface{}) {
				msgChan := ctx.(chan *QueueMessage)
				if _, ok := <-msgChan; ok {
					msgChan <- msg.(*QueueMessage)
				}
			}, msgChan)
		consumer.Start()
		timeout := time.NewTimer(time.Second * time.Duration(ws))
		select {
		case <-timeout.C:
			err = errors.New("kafaka queue receive message timeout")
		case v := <-msgChan:
			msg = *v
		}
		consumer.Close()
	}
	return
}

type batchMessageContext struct {
	fullChan      chan bool
	nbrofMessages int
	msgs          []*QueueMessage
}

func batchQueueMessageHandlerFunc(msg message.Message, ctx interface{}) {
	batch := ctx.(batchMessageContext)
	if len(batch.msgs) < batch.nbrofMessages {
		batch.msgs = append(batch.msgs, msg.(*QueueMessage))
		return
	}
	if _, ok := <-batch.fullChan; ok {
		batch.fullChan <- true
	}
}

func (q *kafkaQueue) BatchReceiveMessages(ws int, numOfMessages int) (msgs []*QueueMessage, err error) {
	if consumer, err := message.NewConsumer(q.config, "kafkaQueue"); err == nil {
		consumer.SetMessageFactory(q)
		ctx := batchMessageContext{
			fullChan:      make(chan bool),
			nbrofMessages: numOfMessages,
			msgs:          []*QueueMessage{},
		}
		consumer.Subscribe(q.queueName, batchQueueMessageHandlerFunc, ctx)
		consumer.Start()
		timeout := time.NewTimer(time.Second * time.Duration(ws))
		select {
		case <-timeout.C:
			err = errors.New("kafaka queue receive message timeout")
		case <-ctx.fullChan:
			msgs = ctx.msgs
		}
		consumer.Close()
		close(ctx.fullChan)
	}
	return
}

func (q *kafkaQueue) DeleteMessage(handle string) error          { return ErrNotImplemented }
func (q *kafkaQueue) BatchDeleteMessages(handles []string) error { return ErrNotImplemented }

func (q *kafkaQueue) PeekMessage(waitSeconds int) (QueueMessage, error) {
	return QueueMessage{}, ErrNotImplemented
}

func (q *kafkaQueue) BatchPeekMessages(wailtSeconds int, numOfMessages int) ([]QueueMessage, error) {
	return []QueueMessage{}, ErrNotImplemented
}

func (q *kafkaQueue) SetMessageVisibility(handle string, seconds int) error { return ErrNotImplemented }
func (q *kafkaQueue) Destroy()                                              {}
