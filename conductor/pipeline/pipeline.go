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
	Start(r data.Reader, ctx data.Context) error
	PushData(data interface{}, ctx data.Context) error
	Close()
}

func New(c config.Config) Pipeline {
	return &defaultPipeline{
		config:       c,
		transformers: []transformer.Transformer{},
		quitch:       make(chan interface{}),
	}
}

type defaultPipeline struct {
	config       config.Config
	extractor    extractor.Extractor
	transformers []transformer.Transformer
	loader       loader.Loader
	reader       data.Reader
	ctx          data.Context
	quitch       chan interface{}
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

func (p *defaultPipeline) Start(r data.Reader, ctx data.Context) error {
	if p.extractor == nil || p.loader == nil {
		return errors.New("extractor or loader is nil")
	}
	p.reader = r
	p.ctx = ctx
	if r != nil {
		go func(p *defaultPipeline) {
			for {
				select {
				case data := <-p.reader.Data():
					p.PushData(data, p.ctx)
				case <-p.quitch:
					return
				}
			}
		}(p)
	}
	return nil
}

func (p *defaultPipeline) PushData(data interface{}, ctx data.Context) error {
	if p.extractor == nil || p.loader == nil {
		return errors.New("extractor or loader is nil")
	}
	// extract data
	value, err := p.extractor.Extract(data, ctx)
	if err != nil {
		return err
	}
	// transfom data
	for _, transformer := range p.transformers {
		if v, err := transformer.Transform(value, ctx); err != nil {
			return err
		} else {
			value = v
		}
	}
	// load data
	return p.loader.Load(value, ctx)
}

func (p *defaultPipeline) Close() {
	if p.reader != nil {
		p.quitch <- true
	}
	p.extractor.Close()
	for _, trans := range p.transformers {
		trans.Close()
	}
	p.loader.Close()
	close(p.quitch)
}
