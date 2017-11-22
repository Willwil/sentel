//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package quto

// Quto Object declaration
type Quto struct {
	QutoId   string           `json:"qutoId" bson:"qutoId"`
	Target   string           `json:"target" bson:"target"`
	TargetId string           `json:"targetId" bson:"targetId"`
	Values   map[string]int64 `json:"values" bson:"values"`
}

// Quto Target declaration
const (
	Tenant  = "tenant"
	Product = "product"
	Device  = "device"
)

// Quto items decleation
const (
	MaxTenants      = "max_tenants"
	MaxProducts     = "max_products"
	CurProducts     = "cur_products"
	MaxConnections  = "max_connections"
	CurConnections  = "cur_connections"
	MaxPublications = "max_publications"
	CurPublications = "cur_publications"
)
