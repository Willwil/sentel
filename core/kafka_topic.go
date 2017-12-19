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
package core

import "encoding/json"

const (
	TopicNameTenant  = "sentel-iot-tenant"
	TopicNameProduct = "sentel-iot-product"
	TopicNameDevice  = "sentel-iot-device"
	TopicNameRule    = "sentel-iot-rule"
)

const (
	ObjectActionRegister   = "register"
	ObjectActionUnregister = "unregister"
	ObjectActionRetrieve   = "retrieve"
	ObjectActionDelete     = "delete"
	ObjectActionUpdate     = "update"
	ObjectActionStart      = "start"
	ObjectActionStop       = "stop"
)

type TopicBase struct {
	Action  string `json:"action"`
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

// Tenantopic
type TenantTopic struct {
	TopicBase
	TenantId string `json:"productId"`
	Action   string `json:"action"`
}

// ProductTopic
type ProductTopic struct {
	TopicBase
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Action      string `json:"action"`
	TenantId    string `json:"tenantId"`
	Replicas    int32  `json:"replicas"`
}

// DeviceTopic
type DeviceTopic struct {
	TopicBase
	DeviceId     string `json:deviceId"`
	DeviceSecret string `json:deviceKey"`
	Action       string `json:"action"`
	ProductId    string `json:"productId"`
}

// RuleTopic declaration
const (
	RuleActionCreate = "create"
	RuleActionRemove = "remove"
	RuleActionUpdate = "update"
	RuleActionStart  = "start"
	RuleActionStop   = "stop"
)

type DataTargetType string

const (
	DataTargetTypeTopic          = "topic"
	DataTargetTypeOuterDatabase  = "outerDatabase"
	DataTargetTypeInnerDatabase  = "innerDatabase"
	DataTargetTypeMessageService = "message"
)

type RuleTopic struct {
	TopicBase
	RuleName    string   `json:"ruleName"`
	DataFormat  string   `json:"dataFormat"`
	Description string   `json:"description"`
	ProductId   string   `json:"productId"`
	DataProcess struct { // select keyword from /productid/topic with condition
		Keyword   string `json:"keyword"`
		Topic     string `json:"topic"`
		Condition string `json:"condition"`
		Sql       string `json:"sql"`
	}
	DataTarget struct {
		Type         DataTargetType `json:"type"`     // Transfer type
		Topic        string         `json:"topic"`    // Transfer data to another topic
		DatabaseHost string         `json:"dbhost"`   // Database host
		DatabaseName string         `json:"database"` // Transfer data to database
		Username     string         `json:"username"` // Database's user name
		Password     string         `json:"password"` // Database's password

	}
	Status     string `json:"status"`
	RuleAction string `json:"action"`
}
