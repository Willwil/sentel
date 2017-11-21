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
	SessionId    string   `json:"sessionId"`
	ClientId     string   `json:"clientId"`
	CleanSession bool     `json:"cleanSession"`
	CreatedAt    string   `json:"createdAt"`
	Topics       []string `json:"topics"`
}

type Device struct {
	DeviceId    string    `json:"deviceId"`
	CreatedAt   time.Time `json:"createdAt"`
	LastUpdated time.Time `json:"lastUpdated"`
	Data        []byte    `json:"data"`
}
