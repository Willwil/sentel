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

package registry

import "time"

// Tenant
type Tenant struct {
	TenantId  string    `bson:"TenantId"`
	Password  string    `bson:"Password"`
	CreatedAt time.Time `bson:"CreatedAt"`
	UpdatedAt time.Time `bson:"UpdatedAt"`
}

// Product
type Product struct {
	TenantId    string    `bson:"TenantId" json:"TenantId"`
	ProductId   string    `bson:"ProductId" json:"ProductId"`
	ProductName string    `bson:"ProductName" json:"ProductName"`
	Description string    `bson:"Description" json:"Description"`
	TimeCreated time.Time `bson:"TimeCreated"`
	TimeUpdated time.Time `bson:"TimeModified"`
	Category    string    `bson:"CategoryId" json:"Category"`
	TopicFlavor string    `bson:"TopicFlavor" json:"TopicFlavor"`
}

// Device
type Device struct {
	ProductId    string    `bson:"ProductId" json:"ProductId"`
	DeviceId     string    `bson:"DeviceId" json:"DeviceId"`
	DeviceName   string    `bson:"DeviceName" json:"DeviceName"`
	DeviceStatus string    `bson:"DeviceStatus"`
	DeviceSecret string    `bson:"DeviceSecret" json:"DeviceSecret"`
	TimeCreated  time.Time `bson:"TimeCreated" json:"TimeCreated"`
	TimeUpdated  time.Time `bson:"TimeUpdated" json:"TimeUpdated"`
}

// Rule
type DataTargetType string

const (
	DataTargetTypeTopic = "topic"
	DataTargetTypeES    = "es"
	DataTargetTypeSQLDB = "sqldb"
	DataTargetTypeCAS   = "cassandra"
)

type RuleDataProcess struct { // select keyword from /productid/topic with condition
	Topic     string   `json:"topic" bson:"Topic"`
	Condition string   `json:"condition" bson:"Condition"`
	Fields    []string `json:"fields" bson:"Fields"`
}
type RuleDataTarget struct {
	Type         DataTargetType `json:"type"`       // Transfer type
	Topic        string         `json:"topic"`      // Transfer data to another topic
	DatabaseHost string         `json:"dbhost"`     // Database host
	DatabaseName string         `json:"dbname"`     // Transfer data to database
	Collection   string         `json:"collection”` // Data Collection
	Username     string         `json:"username"`   // Database's user name
	Password     string         `json:"password"`   // Database's password
}

const (
	RuleStatusIdle    = "idle"
	RuleStatusStarted = "started"
	RuleStatusStoped  = "stoped"
)

type Rule struct {
	ProductId   string          `json:"productId" bson:"ProductId"`
	RuleName    string          `json:"ruleName" bson:"RuleName"`
	DataFormat  string          `json:"format" bson:"Format"`
	Description string          `json:"desc" bson:"Desc"`
	DataProcess RuleDataProcess `json:"dataprocess" bson:"Dataprocess"`
	DataTarget  RuleDataTarget  `json:"datatarget" bson:"Datatarget"`
	Status      string          `json:"status" bson:"Status"`
	TimeCreated time.Time       `json:"timeCreated" bson:"TimeCreated"`
	TimeUpdated time.Time       `json:"timeUpdated" bson:"TimeUpdated"`
}

type TopicFlavor struct {
	FlavorName string           `json:"flavorName" bson:"FlavorName"`
	Builtin    bool             `json:"bultin" bson:"Builtin"`
	TenantId   string           `json:"tenantId" bson:"TenantId"`
	Topics     map[string]uint8 `json:"topic" bson:"Topics"`
}
