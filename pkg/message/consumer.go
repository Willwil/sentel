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

type MessageHandlerFunc func(msg Message, ctx interface{})

type MessageFactory interface {
	CreateMessage(topic string) Message
}

type Consumer interface {
	Subscribe(topic string, handler MessageHandlerFunc, ctx interface{}) error
	Start() error
	Close()
	SetMessageFactory(factory MessageFactory)
}

func NewConsumer(khosts string, clientId string) (Consumer, error) {
	return newKafkaConsumer(khosts, clientId)
}
