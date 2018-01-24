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

const (
	ActionCreate = "register"
	ActionRemove = "unregister"
	ActionUpdate = "update"
	ActionStart  = "start"
	ActionStop   = "stop"
)

// Tenantopic
type TenantTopic struct {
	TenantId string `json:"productId"`
	Action   string `json:"action"`
}

// ProductTopic
type ProductTopic struct {
	ProductId string `json:"productId"`
	Action    string `json:"action"`
	TenantId  string `json:"tenantId"`
	Replicas  int32  `json:"replicas"`
}

// DeviceTopic
type DeviceTopic struct {
	DeviceId     string `json:"deviceId"`
	DeviceSecret string `json:"deviceKey"`
	Action       string `json:"action"`
	ProductId    string `json:"productId"`
}

type RuleTopic struct {
	RuleName  string `json:"ruleName"`
	ProductId string `json:"productId"`
	Action    string `json:"action"`
}
