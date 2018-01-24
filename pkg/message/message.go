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
package message

type MessageHandlerFunc func(topic string, msg []byte, ctx interface{})

type Consumer interface {
	Subscribe(topic string, handler MessageHandlerFunc, ctx interface{}) error
	Start() error
	Close()
}

func NewConsumer(khosts string, clientId string) (Consumer, error) {
	return newKafkaConsumer(khosts, clientId)
}

type Producer interface {
	SendMessage(topic string, val []byte) error
	Close()
}

func NewProducer(khosts string, clientId string, sync bool) (Producer, error) {
	return newKafkaProducer(khosts, clientId, sync)
}

func PostMessage(khosts string, topic string, value []byte) error {
	if producer, err := NewProducer(khosts, "", false); err != nil {
		return err
	} else {
		defer producer.Close()
		return producer.SendMessage(topic, value)
	}
}

func SendMessage(khosts string, topic string, value []byte) error {
	if producer, err := NewProducer(khosts, "", true); err != nil {
		return err
	} else {
		defer producer.Close()
		return producer.SendMessage(topic, value)
	}
}
