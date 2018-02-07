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

package pipeline

import (
	"github.com/cloustone/sentel/conductor/pipeline/extractor"
	"github.com/cloustone/sentel/conductor/pipeline/loader"
	"github.com/cloustone/sentel/conductor/pipeline/transformer"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type Builder struct {
	config config.Config
}

func NewBuilder() *Builder {
	config := config.New("pipeline")
	return &Builder{config: config}
}

func (p *Builder) AddConfig(key string, val interface{}) {
	p.config.AddConfigItem(key, val)
}

func (p *Builder) Pipeline(name string, extractorName string, transformerNames []string, loaderNames []string) (Pipeline, error) {
	ppline := NewPipeline(p.config, name)
	extractor, err := extractor.New(p.config, extractorName)
	if err != nil {
		return nil, err
	}
	ppline.AddExtractor(extractor)

	for _, name := range transformerNames {
		t := transformer.New(p.config, name)
		ppline.AddTransformer(t)
	}

	for _, name := range loaderNames {
		loader, err := loader.New(p.config, name)
		if err != nil {
			return nil, err
		}
		ppline.AddLoader(loader)
	}
	glog.Infof("pipeline '%s' is created", name)
	return ppline, nil
}
