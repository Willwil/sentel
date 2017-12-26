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
	TenantId    string    `bson:"TenantId" json:"tenantId"`
	ProductId   string    `bson:"ProductId" json:"productId"`
	ProductName string    `bson:"ProductId" json:"productId"`
	Description string    `bson:"Description" json:"description"`
	TimeCreated time.Time `bson:"TimeCreated"`
	TimeUpdated time.Time `bson:"TimeModified"`
	Category    string    `bson:"CategoryId" json:"category"`
}

// Device
type Device struct {
	ProductId    string    `bson:"ProductKey" json:"productKey"`
	DeviceId     string    `bson:"DeviceId" json:"deviceId"`
	DeviceName   string    `bson:"DeviceName" json:"deviceName"`
	DeviceStatus string    `bson:"DeviceStatus"`
	DeviceSecret string    `bson:"DeviceSecret" json:"deviceSecret"`
	TimeCreated  time.Time `bson:"TimeCreated" json:"timeCreated"`
	TimeUpdated  time.Time `bson:"TimeModified" json:"timeUpdated`
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
	ProductId   string          `json:"productKey" bson:"ProductKey"`
	RuleName    string          `json:"ruleName" bson:"RuleName"`
	DataFormat  string          `json:"format" bson:"Format"`
	Description string          `json:"desc" bson:"Desc"`
	DataProcess RuleDataProcess `json:"dataprocess" bson"Dataprocess"`
	DataTarget  RuleDataTarget  `json:"datatarget", bson:"Datatarget"`
	Status      string          `json:status, bson:"Status"`
	TimeCreated time.Time       `json:timeCreated, bson:"TimeCreated"`
	TimeUpdated time.Time       `json:timeUpdated, bson:"TimeUpdated"`
}
