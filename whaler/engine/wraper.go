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
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/whaler/pipeline"
	"github.com/cloustone/sentel/whaler/pipeline/extractor"
	"github.com/golang/glog"
)

const usingreader = false

type ruleWraper struct {
	rule   *registry.Rule    // Underlay rule object
	ppline pipeline.Pipeline // Pipeline for the rule
	datach chan interface{}  // Asynchrous data channel
}

func newRuleWraper(c config.Config, r *registry.Rule) (*ruleWraper, error) {
	// construct pipeline builder and add configuration
	builder := pipeline.NewBuilder()
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
	ppline, err := builder.Pipeline(r.RuleName, extractor.EventExtractor, []string{}, []string{r.DataTarget.Type})
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
	glog.Infof("executing rule '%s' for product '%s'...", p.rule.RuleName, p.rule.ProductId)
	if usingreader {
		p.datach <- e
		return nil
	} else {
		if pe, ok := e.(*event.TopicPublishEvent); ok {
			return p.ppline.PushData(pe)
		}
		return errors.New("invalid event type")
	}
}

func (p *ruleWraper) Data() chan interface{} { return p.datach }
func (p *ruleWraper) close()                 { p.ppline.Close() }
