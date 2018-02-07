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
		ppline.Start(p)
	}
	return p, nil
}

func buildPipeline(c config.Config, r *registry.Rule) (pipeline.Pipeline, error) {
	// make a new configuration and add default common configuration
	config := config.New("pipeline")
	config.AddConfigItem("productId", r.ProductId)
	config.AddConfigItem("ruleName", r.RuleName)
	config.AddConfigItem("dataprocess", r.DataProcess)
	config.AddConfigItem("datatarget", r.DataTarget)

	// add datatarget specified settings
	switch r.DataTarget.Type {
	case registry.DataTargetTypeTopic:
		config.AddConfigItem("mongo", c.MustString("mongo"))
		config.AddConfigItem("message_server", c.MustString("kafka"))
		config.AddConfigItem("topic", r.DataTarget.Topic)
	default:
	}

	// conduct pipeline
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
		return p.ppline.PushData(e)
	}
}

func (p *ruleWraper) Data() chan interface{} { return p.datach }
func (p *ruleWraper) close()                 { p.ppline.Close() }
