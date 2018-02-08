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

package collector

import (
	"github.com/cloustone/sentel/iotmanager/mgrdb"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

// Stat
type StatsTopic struct {
	TopicName string
	mgrdb.Stats
	Action string `json:"action"`
}

func (p *StatsTopic) Topic() string        { return TopicNameStats }
func (p *StatsTopic) SetTopic(name string) {}
func (p *StatsTopic) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *StatsTopic) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *StatsTopic) handleTopic(c config.Config, ctx context) error {
	switch p.Action {
	case ObjectActionUpdate:
		ctx.db.AddStatsHistory(p.Stats)
		ctx.db.UpdateStats(p.Stats)
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil
}
