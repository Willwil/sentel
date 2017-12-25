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

package db

import "time"

// Tenant
type Tenant struct {
	TenantId  string    `bson:"Name"`
	Password  string    `bson:"Password"`
	CreatedAt time.Time `bson:"CreatedAt"`
	UpdatedAt time.Time `bson:"UpdatedAt"`
}

// Product
type Product struct {
	TenantId     string    `bson:"TenantId" json:"tenantId"`
	ProductId    string    `bson:"ProductId"`
	Description  string    `bson:"Description"`
	TimeCreated  time.Time `bson:"TimeCreated"`
	TimeModified time.Time `bson:"TimeModified"`
	CategoryId   string    `bson:"CategoryId"`
	ProductKey   string    `bson:"ProductKey"`
}

// Device
type Device struct {
	DeviceId     string    `bson:"DeviceId"`
	ProductId    string    `bson:"ProductId"`
	ProductKey   string    `bson:"ProductKey"`
	DeviceStatus string    `bson:"DeviceStatus"`
	DeviceSecret string    `bson:"DeviceSecret"`
	TimeCreated  time.Time `bson:"TimeCreated"`
	TimeModified time.Time `bson:"TimeModified"`
}

// Rule
type DataTargetType string

const (
	DataTargetTypeTopic          = "topic"
	DataTargetTypeOuterDatabase  = "outerDatabase"
	DataTargetTypeInnerDatabase  = "innerDatabase"
	DataTargetTypeMessageService = "message"
)

type RuleDataProcess struct { // select keyword from /productid/topic with condition
	Topic     string   `json:"topic" bson:"topic"`
	Condition string   `json:"condition" bson:"condiction"`
	Fields    []string `json:"fields" bson:"fields"`
}
type RuleDataTarget struct {
	Type         DataTargetType `json:"type"`     // Transfer type
	Topic        string         `json:"topic"`    // Transfer data to another topic
	DatabaseHost string         `json:"dbhost"`   // Database host
	DatabaseName string         `json:"database"` // Transfer data to database
	Username     string         `json:"username"` // Database's user name
	Password     string         `json:"password"` // Database's password
}

const (
	RuleStatusIdle    = "idle"
	RuleStatusStarted = "started"
	RuleStatusStoped  = "stoped"
)

type Rule struct {
	TenantId    string          `json:"tenantId" bson:"tenantId"`
	ProductKey  string          `json:"productKey" bson:"productKey"`
	RuleName    string          `json:"ruleName" bson:"ruleName"`
	DataFormat  string          `json:"format" bson:"format"`
	Description string          `json:"desc" bson:"desc"`
	DataProcess RuleDataProcess `json:"dataprocess" bson"dataprocess"`
	DataTarget  RuleDataTarget  `json:"datatarget", bson:"datatarget"`
	Status      string          `json:status, bson:"status"`
	TimeCreated time.Time       `json:timeCreated, bson:"timeCreated"`
	TimeUpdated time.Time       `json:timeUpdated, bson:"timeUpdated"`
}
