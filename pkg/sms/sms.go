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

package sms

type Message struct {
	to        []string
	project   string
	variables map[string]string
}

type Dialer interface {
	DialAndSend(*Message) error
}

func NewMessage() *Message {
	return &Message{
		to:        []string{},
		variables: make(map[string]string),
	}
}

func NewDialer(configs map[string]string) Dialer {
	return newSubmailDialer(configs)
}

func (m *Message) AddTo(to ...string) {
	m.to = append(m.to, to...)
}
func (m *Message) SetProject(project string) {
	m.project = project
}
func (m *Message) AddVariable(key string, val string) {
	m.variables[key] = val
}
