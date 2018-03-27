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

import "submail"

type submailDialer struct {
	configs      map[string]string
	messagexsend *submail.MessageXSend
}

func newSubmailDialer(configs map[string]string) Dialer {
	return &submailDialer{
		configs:      configs,
		messagexsend: submail.CreateMessageXSend(),
	}
}

func (s *submailDialer) DialAndSend(m *Message) error {
	for _, to := range m.to {
		submail.MessageXSendAddTo(s.messagexsend, to)
	}
	submail.MessageXSendSetProject(s.messagexsend, m.project)
	for k, v := range m.variables {
		submail.MessageXSendAddVar(s.messagexsend, k, v)
	}
	submail.MessageXSendRun(submail.MessageXSendBuildRequest(s.messagexsend), s.configs)
	return nil
}
