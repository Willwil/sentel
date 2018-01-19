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
	"github.com/cloustone/sentel/conductor/ETL"
	"github.com/cloustone/sentel/conductor/extractor"
	"github.com/cloustone/sentel/conductor/loader"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

type ruleWraper struct {
	rule *registry.Rule
	etl  ETL.ETL
	// curEvent *event.Event
	datach chan *event.Event
}

func newRuleWraper(c config.Config, r *registry.Rule) (*ruleWraper, error) {
	// extraction configuration at first
	config := config.New()
	options := make(map[string]string)
	options["mongo"] = c.MustString("conductor", "mongo")
	options["message_server"] = c.MustString("conductor", "kafka")
	if r.DataTarget.Type == registry.DataTargetTypeTopic {
		options["topic"] = r.DataTarget.Topic
	}
	config.AddConfigSection("etl", options)

	// build ETL
	etl := ETL.New(config)
	extractor, _ := extractor.New(config, "event")
	loader, _ := loader.New(config, string(r.DataTarget.Type))
	if extractor == nil || loader == nil {
		return nil, errors.New("etl initialization failed")
	}
	etl.AddExtractor(extractor).AddLoader(loader)

	return &ruleWraper{
		etl:    etl,
		rule:   r,
		datach: make(chan *event.Event),
	}, nil

}

func (p *ruleWraper) Read() (interface{}, error) {
	//return p.curEvent, nil
	return <-p.datach, nil
}

func (p *ruleWraper) executeETL(e *event.Event) error {
	glog.Infof("conductor executing rule '%s' for product '%s'...", p.rule.RuleName, p.rule.ProductId)
	//p.curEvent = e
	p.datach <- e
	err := p.etl.Run(p, p.rule)
	//p.curEvent = nil
	return err
}

func (p *ruleWraper) close() {
	p.etl.Close()
	close(p.datach)
}
