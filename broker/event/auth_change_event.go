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

package event

import "github.com/cloustone/sentel/pkg/registry"

type AuthChangeEvent struct {
	BrokerID    string               `json:"brokerId"`   // Broker identifier where event come from
	Type        uint32               `json:"type"`       // Event type
	ClientID    string               `json:"clientID"`   // Client identifier where event come from
	Persistent  bool                 `json:"persistent"` // Whether the session is persistent
	ProductID   string               `json:"productId"`
	TopicFlavor registry.TopicFlavor `json:"topicFlavor"`
}

func (p *AuthChangeEvent) SetBrokerId(brokerId string) { p.BrokerID = brokerId }
func (p *AuthChangeEvent) SetType(eventType uint32)    { p.Type = eventType }
func (p *AuthChangeEvent) SetClientId(clientID string) { p.ClientID = clientID }
func (p *AuthChangeEvent) GetBrokerId() string         { return p.BrokerID }
func (p *AuthChangeEvent) GetType() uint32             { return AuthChange }
func (p *AuthChangeEvent) GetClientId() string         { return p.ClientID }
func (p *AuthChangeEvent) Serialize() ([]byte, error)  { return nil, nil }
