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

// Session
type Session struct {
	TopicName          string
	Action             string `json:"action"`
	ClientId           string `json:"clientId"`
	CleanSession       bool   `json:"cleanSession"`
	MessageMaxInflight uint64 `json:"messageMaxInflight"`
	MessageInflight    uint64 `json:"messageInflight"`
	MessageInQueue     uint64 `json:"messageInQueue"`
	MessageDropped     uint64 `json:"messageDropped"`
	AwaitingRel        uint64 `json:"awaitingRel"`
	AwaitingComp       uint64 `json:"awaitingComp"`
	AwaitingAck        uint64 `json:"awaitingAck"`
	CreatedAt          string `json:"createdAt"`
}

func (p *ManagerDB) GetSession(clientId string) (Session, error) {
	c := p.session.C(collectionSessions)
	session := Session{}
	err := c.Find(bson.M{"ClientId": clientId}).One(&session)
	return session, err
}

func (p *ManagerDB) UpdateSession(s Session) error {
	c := p.session.C(collectionSessions)
	session := Session{}
	if err := c.Find(bson.M{"ClientId": s.ClientId}).One(&session); err != nil {
		return c.Insert(s)
	} else {
		return c.Update(session, s)
	}
}

func (p *ManagerDB) GetNodeSessions(nodeId string) []Session {
	c := p.session.C(collectionSessions)
	sessions := []Session{}
	c.Find(bson.M{"NodeId": nodeId}).Iter().All(&sessions)
	return sessions
}

func (p *ManagerDB) GetSessionWithNode(nodeId string, clientId string) (Session, error) {
	c := p.session.C(collectionSessions)
	session := Session{}
	err := c.Find(bson.M{"NodeId": nodeId, "ClientId": clientId}).One(&session)
	return session, err
}
