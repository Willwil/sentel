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
	"fmt"

	"github.com/cloustone/sentel/iotmanager/mgrdb"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

// MetricTopic
type MetricTopic struct {
	TopicName string
	mgrdb.Metric
}

func (p *MetricTopic) Topic() string        { return TopicNameMetric }
func (p *MetricTopic) SetTopic(name string) {}
func (p *MetricTopic) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *MetricTopic) Deserialize(buf []byte, opt message.SerializeOption) error {
	return message.Deserialize(buf, opt, p)
}

func (p *MetricTopic) handleTopic(c config.Config, ctx context) error {
	switch p.Action {
	case ObjectActionUpdate:
		ctx.db.AddMetricHistory(p.Metric)
		ctx.db.UpdateMetric(p.Metric)
	default:
	}
	return fmt.Errorf("invalid topic '%s' action '%s'", p.Topic(), p.Action)
}
