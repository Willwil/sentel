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

package metadata

import "time"

type Session struct {
	SessionId string    `json:"sessionId" bson:"sessionId"`
	ClientId  string    `json:"clientId" bson:"clientId"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	Topics    []string  `json:"topics" bson:"topics"`
}

type Device struct {
	DeviceId    string    `json:"deviceId" bson:"deviceId"`
	CreatedAt   time.Time `json:"createdAt" bson:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated" bson:"lastUpdated"`
	Data        []byte    `json:"data" bson:"data"`
}
