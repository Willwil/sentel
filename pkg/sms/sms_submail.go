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

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/sms/submail"
)

type submailSms struct {
	configs      map[string]string
	messagexsend *submail.MessageXSend
}

func newSubmailSms(c config.Config) (SMS, error) {
	appid, _ := c.StringWithSection("sms", "appid")
	appkey, _ := c.StringWithSection("sms", "appkey")
	signtype, _ := c.StringWithSection("sms", "signtype")
	if appid == "" || appkey == "" || signtype == "" {
		return nil, errors.New("invalid submail configs setting")
	}
	configs := make(map[string]string)
	configs["appid"] = appid
	configs["appkey"] = appkey
	configs["signtype"] = signtype
	return &submailSms{
		configs:      configs,
		messagexsend: submail.CreateMessageXSend(),
	}, nil
}

func (s *submailSms) AddTo(to string) {
	submail.MessageXSendAddTo(s.messagexsend, to)
}

func (s *submailSms) SetProject(project string) {
	submail.MessageXSendSetProject(s.messagexsend, project)
}

func (s *submailSms) AddVariable(key string, val string) {
	submail.MessageXSendAddVar(s.messagexsend, key, val)
}
func (s *submailSms) Send() error {
	submail.MessageXSendRun(submail.MessageXSendBuildRequest(s.messagexsend), s.configs)
	return nil
}
