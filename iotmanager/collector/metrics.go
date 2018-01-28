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
	"context"

	db "github.com/cloustone/sentel/iotmanager/database"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

// Metric
type Metric db.Metric

func (p *Metric) Topic() string        { return TopicNameMetric }
func (p *Metric) SetTopic(name string) {}
func (p *Metric) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *Metric) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *Metric) handleTopic(c config.Config, ctx context.Context) error {
	dbc, err := db.NewManagerDB(c)
	if err != nil {
		return err
	}
	defer dbc.Close()

	switch p.Action {
	case ObjectActionUpdate:
		// dbc.UpdateMetric(p)
		// TODO:save history data
	case ObjectActionDelete:
	case ObjectActionRegister:
	default:
	}
	return nil

}
