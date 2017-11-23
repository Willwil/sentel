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
	QutoId   string `json:"qutoId" bson:"qutoId"`
	BrokerId string `json:"brokerId" bson:"brokerId"`
	Key      string `json:"key" bson:"key"`
	Value    uint64 `json:"value" bson:"value"`
}

// Quto items decleation
const (
	MaxConnections  = "max_connections"
	CurConnections  = "cur_connections"
	MaxPublications = "max_publications"
	CurPublications = "cur_publications"
)
