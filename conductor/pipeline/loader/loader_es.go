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

package loader

import (
	"errors"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/conductor/pipeline/data"
	"github.com/cloustone/sentel/pkg/config"
)

type elasticLoader struct {
	config config.Config
}

func newESLoader(c config.Config) (Loader, error) {
	return &elasticLoader{config: c}, nil
}

func (p *elasticLoader) Name() string { return "elastic" }
func (p *elasticLoader) Close()       {}

func (p *elasticLoader) Load(f *data.DataFrame, ctx *data.Context) error {
	topic, ok1 := ctx.Get("topic").(string)
	e, ok2 := ctx.Get("event").(*event.TopicPublishEvent)
	if !ok1 || !ok2 || topic == "" || e == nil {
		return errors.New("invalid topic")
	}
	if _, err := f.Serialize(); err != nil {
		return err
	} else {
		// TODO
		return nil
	}
}
