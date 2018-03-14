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

// TopicSubscribeEvent
type TopicSubscribeEvent struct {
	BrokerID   string `json:"brokerId"`   // Broker identifier where event come from
	Type       uint32 `json:"type"`       // Event type
	ClientID   string `json:"clientId"`   // Client identifier where event come from
	Persistent bool   `json:"persistent"` // Whether the session is persistent
	Topic      string `json:"topic"`      // Topic
	Qos        uint8  `json:"qos"`
	Data       []byte `json:"data"` // Topic data
	Retain     bool   `json:"retain"`
}

func (p *TopicSubscribeEvent) SetBrokerId(brokerID string) {
	p.BrokerID = brokerID
}

func (p *TopicSubscribeEvent) SetType(eventType uint32) {
	p.Type = eventType
}

func (p *TopicSubscribeEvent) SetClientId(clientID string) {
	p.ClientID = clientID
}
func (p *TopicSubscribeEvent) GetBrokerId() string        { return p.BrokerID }
func (p *TopicSubscribeEvent) GetType() uint32            { return TopicSubscribe }
func (p *TopicSubscribeEvent) GetClientId() string        { return p.ClientID }
func (p *TopicSubscribeEvent) Serialize() ([]byte, error) { return nil, nil }

// TopicUnsubscribeEvent
type TopicUnsubscribeEvent struct {
	BrokerID   string `json:"brokerId"`   // Broker identifier where event come from
	Type       uint32 `json:"type"`       // Event type
	ClientID   string `json:"clientId"`   // Client identifier where event come from
	Persistent bool   `json:"persistent"` // Whether the session is persistent
	Topic      string `json:"topic"`      // Topic
	Data       []byte `json:"data"`       // Topic data
}

func (p *TopicUnsubscribeEvent) SetBrokerId(brokerId string) {
	p.BrokerID = brokerId
}

func (p *TopicUnsubscribeEvent) SetType(eventType uint32) {
	p.Type = eventType
}

func (p *TopicUnsubscribeEvent) SetClientId(clientId string) {
	p.ClientID = clientId
}
func (p *TopicUnsubscribeEvent) GetBrokerId() string        { return p.BrokerID }
func (p *TopicUnsubscribeEvent) GetType() uint32            { return TopicUnsubscribe }
func (p *TopicUnsubscribeEvent) GetClientId() string        { return p.ClientID }
func (p *TopicUnsubscribeEvent) Serialize() ([]byte, error) { return nil, nil }

// TopicPublishEvent
type TopicPublishEvent struct {
	BrokerID   string `json:"brokerId"`   // Broker identifier where event come from
	Type       uint32 `json:"type"`       // Event type
	ClientID   string `json:"clientId"`   // Client identifier where event come from
	ProductID  string `json:"productId"`  // Product identifier where event come from
	Persistent bool   `json:"persistent"` // Whether the session is persistent
	ID         uint16 `json:"id"`         // Message Id
	Topic      string `json:"topic"`      // Topic
	Payload    []byte `json:"data"`       // Topic data
	Qos        uint8  `json:"qos"`
	Direction  uint8  `json:"direction"`
	Retain     bool   `json:"retain"`
	Dup        bool   `json:"dup"`
}

func (p *TopicPublishEvent) SetBrokerId(brokerId string)   { p.BrokerID = brokerId }
func (p *TopicPublishEvent) SetType(eventType uint32)      { p.Type = eventType }
func (p *TopicPublishEvent) SetClientId(clientId string)   { p.ClientID = clientId }
func (p *TopicPublishEvent) SetProductId(productId string) { p.ProductID = productId }
func (p *TopicPublishEvent) GetBrokerId() string           { return p.BrokerID }
func (p *TopicPublishEvent) GetType() uint32               { return TopicPublish }
func (p *TopicPublishEvent) GetClientId() string           { return p.ClientID }
func (p *TopicPublishEvent) GetProductId() string          { return p.ProductID }
func (p *TopicPublishEvent) Serialize() ([]byte, error)    { return nil, nil }
