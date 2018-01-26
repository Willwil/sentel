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

const (
	TopicNameTenant  = "sentel-iot-tenant"
	TopicNameProduct = "sentel-iot-product"
	TopicNameDevice  = "sentel-iot-device"
	TopicNameRule    = "sentel-iot-rule"
)

type Message interface {
	Topic() string
	SetTopic(name string)
	Serialize(opt SerializeOption) ([]byte, error)
	Deserialize(buf []byte, opt SerializeOption) error
}

type builtinFactory struct{}

func (p *builtinFactory) CreateMessage(topicName string) Message {
	switch topicName {
	case TopicNameTenant:
		return &Tenant{}
	case TopicNameProduct:
		return &Product{}
	case TopicNameRule:
		return &Rule{}
	default:
		return &Broker{TopicName: topicName}
	}
}
