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
	BrokerId    string               `json:"brokerId"`   // Broker identifier where event come from
	Type        uint32               `json:"type"`       // Event type
	ClientId    string               `json:"clientId"`   // Client identifier where event come from
	Persistent  bool                 `json:"persistent"` // Whether the session is persistent
	ProductId   string               `json:"productId"`
	TopicFlavor registry.TopicFlavor `json:"topicFlavor"`
}

func (p *AuthChangeEvent) SetBrokerId(brokerId string) { p.BrokerId = brokerId }
func (p *AuthChangeEvent) SetType(eventType uint32)    { p.Type = eventType }
func (p *AuthChangeEvent) SetClientId(clientId string) { p.ClientId = clientId }
func (p *AuthChangeEvent) GetBrokerId() string         { return p.BrokerId }
func (p *AuthChangeEvent) GetType() uint32             { return AuthChange }
func (p *AuthChangeEvent) GetClientId() string         { return p.ClientId }
func (p *AuthChangeEvent) Serialize() ([]byte, error)  { return nil, nil }
