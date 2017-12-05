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
	Name      string    `bson:"Name"`
	Password  string    `bson:"Password"`
	CreatedAt time.Time `bson:"CreatedAt"`
	UpdatedAt time.Time `bson:"UpdatedAt"`
}

// Product
type Product struct {
	Id           string    `bson:"Id"`
	Name         string    `bson:"Name"`
	Description  string    `bson:"Description"`
	TimeCreated  time.Time `bson:"TimeCreated"`
	TimeModified time.Time `bson:"TimeModified"`
	CategoryId   string    `bson:"CategoryId"`
	ProductKey   string    `bson:"ProductKey"`
}

// Device
type Device struct {
	Id           string    `bson:"Id"`
	Name         string    `bson:"Name"`
	ProductId    string    `bson:"ProductId"`
	ProductKey   string    `bson:"productKey"`
	DeviceStatus string    `bson:"deviceStatus"`
	DeviceSecret string    `bson:"deviceSecret"`
	TimeCreated  time.Time `bson:"timeCreated"`
	TimeModified time.Time `bson:"TimeModified"`
}

// Rule
type Rule struct {
	Id           string    `bson:"Id"`
	Name         string    `bson:"Name"`
	ProductId    string    `bson:"ProductId"`
	TimeCreated  time.Time `bson:"TimeCreated"`
	TimeModified time.Time `bson:"TimeModified"`
	Status       string    `bson:"Status"`
	Method       string    `bson:"Method"`
	Target       string    `bson:"Target"`
}
