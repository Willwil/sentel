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

package sessionmgr

import "github.com/cloustone/sentel/broker/queue"

type virtualSession struct {
	clientId   string
	persistent bool
	queue      queue.Queue
	brokerId   string
}

// Id return identifier of session
func (p *virtualSession) Id() string { return p.clientId }

// IsPersistent return wether the session is persistent
func (p *virtualSession) IsPersistent() bool { return p.persistent }

// IsValidd return true if the session state valid
func (p *virtualSession) IsValid() bool { return true }

// Info return session's information
func (p *virtualSession) Info() *SessionInfo { return &SessionInfo{} }

func (p *virtualSession) BrokerId() string { return p.brokerId }
func newVirtualSession(brokerId, clientId string, persistent bool) (Session, error) {
	q, err := queue.NewQueue(clientId, persistent)
	if err != nil {
		return nil, err
	}
	return &virtualSession{
		brokerId:   brokerId,
		clientId:   clientId,
		persistent: persistent,
		queue:      q,
	}, nil
}
