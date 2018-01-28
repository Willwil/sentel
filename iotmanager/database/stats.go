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

// Stat
type Stats struct {
	TopicName  string
	NodeName   string            `json:"nodeName"`
	Service    string            `json:"service"`
	Action     string            `json:"action"`
	UpdateTime time.Time         `json:"updateTime"`
	Values     map[string]uint64 `json:"values"`
}
