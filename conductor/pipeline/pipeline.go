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
	"errors"

	"github.com/cloustone/sentel/conductor/data"
	"github.com/cloustone/sentel/conductor/extractor"
	"github.com/cloustone/sentel/conductor/loader"
	"github.com/cloustone/sentel/conductor/transformer"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type Pipeline interface {
	AddExtractor(extractor.Extractor) Pipeline
	AddTransformer(transformer.Transformer) Pipeline
	AddLoader(loader.Loader)
	Push(r data.Reader, ctx interface{}) error
	Close()
}

func New(c config.Config) Pipeline {
	return &defaultPipeline{
		config:       c,
		transformers: []transformer.Transformer{},
	}
}

type defaultPipeline struct {
	config       config.Config
	extractor    extractor.Extractor
	transformers []transformer.Transformer
	loader       loader.Loader
}

func (p *defaultPipeline) AddExtractor(t extractor.Extractor) Pipeline {
	if p.extractor != nil {
		glog.Error("extractor already exist")
	}
	p.extractor = t
	return p
}

func (p *defaultPipeline) AddTransformer(t transformer.Transformer) Pipeline {
	p.transformers = append(p.transformers, t)
	return p
}
func (p *defaultPipeline) AddLoader(t loader.Loader) {
	if p.loader != nil {
		glog.Error("loader already exist")
	}
	p.loader = t
}

func (p *defaultPipeline) Push(r data.Reader, ctx interface{}) error {
	data, _ := p.extractor.Extract(r, ctx)
	// transfom data
	for _, transformer := range p.transformers {
		data, _ = transformer.Transform(data, ctx)
	}
	if p.loader == nil {
		return errors.New("loader is nill")
	}
	return p.loader.Load(data, ctx)
}

func (p *defaultPipeline) Close() {
	p.extractor.Close()
	for _, trans := range p.transformers {
		trans.Close()
	}
	if p.loader != nil {
		p.loader.Close()
	}
}
