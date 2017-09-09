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

package base

type ClientInfo struct {
	UserName     string
	CleanSession bool
	PeerName     string
	ConnectTime  string
}

type ProtocolService interface {
	// Stats and Metrics
	GetStats() *Stats
	GetMetrics() *Metrics
	// Client
	GetClients() []*ClientInfo
	GetClient(id string) *ClientInfo
	KickoffClient(id string) error
	// Session
	GetSessions() []*SessionInfo
	GetSession(id string) *SessionInfo
}
