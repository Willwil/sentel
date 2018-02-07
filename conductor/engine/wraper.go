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

package engine

import (
	"errors"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/conductor/pipeline"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

const usingreader = false

type ruleWraper struct {
	rule   *registry.Rule
	ppline pipeline.Pipeline
	datach chan interface{}
}

func newRuleWraper(c config.Config, r *registry.Rule) (*ruleWraper, error) {
	// make a new configuration and add default common configuration
	builder := pipeline.NewBuilder(c)
	builder.AddConfig("productId", r.ProductId)
	builder.AddConfig("ruleName", r.RuleName)
	builder.AddConfig("dataprocess", r.DataProcess)
	builder.AddConfig("datatarget", r.DataTarget)

	// add datatarget specified settings
	switch r.DataTarget.Type {
	case registry.DataTargetTypeTopic:
		builder.AddConfig("mongo", c.MustString("mongo"))
		builder.AddConfig("message_server", c.MustString("kafka"))
		builder.AddConfig("topic", r.DataTarget.Topic)
	default:
	}
	// build pipeline
	ppline, err := builder.Pipeline("event", []string{}, []string{r.DataTarget.Type})
	if err != nil {
		return nil, err
	}
	p := &ruleWraper{
		ppline: ppline,
		rule:   r,
		datach: make(chan interface{}, 1),
	}
	if usingreader {
		ppline.Start(p)
	}
	return p, nil
}

func (p *ruleWraper) handle(e event.Event) error {
	glog.Infof("conductor executing rule '%s' for product '%s'...", p.rule.RuleName, p.rule.ProductId)
	if usingreader {
		p.datach <- e
		return nil
	} else {
		if pe, ok := e.(*event.TopicPublishEvent); ok {
			return p.ppline.PushData(pe)
		}
		return errors.New("iinvalid event type")
	}
}

func (p *ruleWraper) Data() chan interface{} { return p.datach }
func (p *ruleWraper) close()                 { p.ppline.Close() }
