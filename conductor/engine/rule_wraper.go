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
	"github.com/cloustone/sentel/conductor/data"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/golang/glog"
)

type ruleWraper struct {
	rule *registry.Rule
}

func (p *ruleWraper) execute(c config.Config, e *event.Event) error {
	glog.Infof("conductor executing rule '%s' for product '%s'...", p.rule.RuleName, p.rule.ProductId)
	dataProcessor, _ := data.NewProcessor(c, p.rule)
	endpoint, err := data.NewEndpoint(c, p.rule)
	if err != nil {
		return err
	}
	if result, err := dataProcessor.Execute(e); err != nil {
		return err
	} else {
		return endpoint.Write(result)
	}
}