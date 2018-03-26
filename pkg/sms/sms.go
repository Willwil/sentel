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
	"fmt"

	"github.com/cloustone/sentel/pkg/config"
)

type SMS interface {
	AddTo(string)
	SetProject(string)
	AddVariable(key string, val string)
	Send() error
}

func New(c config.Config) (SMS, error) {
	vendor, _ := c.StringWithSection("sms", "vendor")
	switch vendor {
	case "submail":
		return newSubmailSms(c)
	}
	return nil, fmt.Errorf("invalid sms vendor '%s'", vendor)
}
