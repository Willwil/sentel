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

package util

import (
	"encoding/json"
)

const (
	TopicNameTenant  = "tenant"
	TopicNameProduct = "product"
	TopicNameDevice  = "device"
	TopicNameRule    = "rule"
)

const (
	ObjectActionRegister   = "register"
	ObjectActionUnregister = "unregister"
	ObjectActionRetrieve   = "retrieve"
	ObjectActionDelete     = "delete"
	ObjectActionUpdate     = "update"
)

// TopicBase
type TopicBase struct {
	encoded []byte
	err     error
}

func (p *TopicBase) ensureEncoded() {
	if p.encoded == nil && p.err == nil {
		p.encoded, p.err = json.Marshal(p)
	}
}

func (p *TopicBase) Length() int {
	p.ensureEncoded()
	return len(p.encoded)
}

func (p *TopicBase) Encode() ([]byte, error) {
	p.ensureEncoded()
	return p.encoded, p.err
}

// ProductTopic
type ProductTopic struct {
	TopicBase
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Action      string `json:"action"`
	TenantId    string `json:"tenantId"`
}

// DeviceTopic
type DeviceTopic struct {
	TopicBase
	DeviceId     string `json:deviceId"`
	DeviceSecret string `json:deviceKey"`
	Action       string `json:"action"`
	ProductId    string `json:"productId"`
}

// RuleTopic
type RuleTopic struct {
	TopicBase
	RuleName  string `json:"ruleName"`
	ProductId string `json:"productId"`
	RuleId    string `json:"ruleId"`
	Action    string `json:"action"`
	encoded   []byte
	err       error
}

// Tenantopic
type TenantTopic struct {
	TopicBase
	TenantId string `json:"productId"`
	Action   string `json:"action"`
	encoded  []byte
	err      error
}
