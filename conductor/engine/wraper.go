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
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/conductor/etl"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

type ruleWraper struct {
	rule     *registry.Rule
	etl      *etl.ETL
	curEvent *event.Event
}

func newRuleWraper(c config.Config, r *registry.Rule) (*ruleWraper, error) {
	// extraction configuration at first
	config := config.New()
	options := make(map[string]string)
	options["extractor"] = "event"
	options["transformer"] = "no"
	options["loader"] = string(r.DataTarget.Type)
	options["mongo"] = c.MustString("conductor", "mongo")
	options["kafka"] = c.MustString("conductor", "kafka")
	config.AddConfigSection("etl", options)
	etl, err := etl.New(config)
	if err != nil {
		return nil, err
	}
	return &ruleWraper{etl: etl, rule: r}, nil

}

func (p *ruleWraper) Read() (interface{}, error) {
	return p.curEvent, nil
}

func (p *ruleWraper) executeETL(e *event.Event) error {
	glog.Infof("conductor executing rule '%s' for product '%s'...", p.rule.RuleName, p.rule.ProductId)
	p.curEvent = e
	err := p.etl.Run(p, p.rule)
	p.curEvent = nil
	return err
}
