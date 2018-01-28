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

import "gopkg.in/mgo.v2/bson"

// Client
type Client struct {
	TopicName       string
	ClientId        string `json:"clientId" bson:"ClientId"`
	UserName        string `json:"userName" bson:"UserName"`
	IpAddress       string `json:"ipAddress" bson:"IpAddress"`
	Port            uint16 `json:"port" bson:"Port"`
	CleanSession    bool   `json:"cleanSession" bson:"CleanSession"`
	ProtocolVersion string `json:"protocolVersion" bson:"ProtocolVersion"`
	Keepalive       uint16 `json:"keepalive" bson:"Keepalive"`
	ConnectedAt     string `json:"connectedAt" bson:"ConnectedAt"`
}

func (p *ManagerDB) GetClient(clientId string) (*Client, error) {
	c := p.session.C(collectionClients)
	client := Client{}
	err := c.Find(bson.M{"ClientId": clientId}).One(&client)
	return &client, err
}

func (p *ManagerDB) UpdateClient(client *Client) error {
	c := p.session.C(collectionClients)
	result := Client{}
	if err := c.Find(bson.M{"ClientId": client.ClientId}).One(&result); err == nil {
		// Existed client found
		return c.Update(result, p)
	} else {
		return c.Insert(p)
	}
}
