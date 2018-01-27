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
	"github.com/cloustone/sentel/conductor/pipeline/data"
	"github.com/cloustone/sentel/conductor/pipeline/extractor"
	"github.com/cloustone/sentel/conductor/pipeline/loader"
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
	ppline, err := buildPipeline(c, r)
	if err != nil {
		return nil, err
	}
	p := &ruleWraper{
		ppline: ppline,
		rule:   r,
		datach: make(chan interface{}, 1),
	}
	if usingreader {
		ctx := data.NewContext()
		ctx.Set("rule", r)
		ppline.Start(p, ctx)
	}
	return p, nil
}

func buildPipeline(c config.Config, r *registry.Rule) (pipeline.Pipeline, error) {
	config := config.New("conductor")
	options := make(map[string]string)
	options["mongo"] = c.MustString("mongo")
	options["message_server"] = c.MustString("kafka")
	if r.DataTarget.Type == registry.DataTargetTypeTopic {
		options["topic"] = r.DataTarget.Topic
	}
	config.AddConfigSection("pipeline", options)

	ppline := pipeline.New(config)
	extractor, _ := extractor.New(config, "event")
	loader, _ := loader.New(config, string(r.DataTarget.Type))
	if extractor == nil || loader == nil {
		return nil, errors.New("pipeline initialization failed")
	}
	ppline.AddExtractor(extractor).AddLoader(loader)
	return ppline, nil
}

func (p *ruleWraper) handle(e event.Event) error {
	glog.Infof("conductor executing rule '%s' for product '%s'...", p.rule.RuleName, p.rule.ProductId)
	if usingreader {
		p.datach <- e
		return nil
	} else {
		ctx := data.NewContext()
		ctx.Set("rule", p.rule)
		return p.ppline.PushData(e, ctx)
	}
}

func (p *ruleWraper) Data() chan interface{} { return p.datach }
func (p *ruleWraper) close()                 { p.ppline.Close() }
