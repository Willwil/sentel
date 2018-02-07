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

	"github.com/cloustone/sentel/conductor/pipeline/data"
	"github.com/cloustone/sentel/conductor/pipeline/extractor"
	"github.com/cloustone/sentel/conductor/pipeline/loader"
	"github.com/cloustone/sentel/conductor/pipeline/transformer"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type Pipeline interface {
	AddExtractor(extractor.Extractor) Pipeline
	AddTransformer(transformer.Transformer) Pipeline
	AddLoader(loader.Loader) Pipeline
	Start(r data.Reader) error
	PushData(data interface{}) error
	Close()
}

func New(c config.Config) Pipeline {
	return &defaultPipeline{
		config:       c,
		transformers: []transformer.Transformer{},
		loaders:      []loader.Loader{},
		quitch:       make(chan interface{}),
	}
}

type defaultPipeline struct {
	config       config.Config
	extractor    extractor.Extractor
	transformers []transformer.Transformer
	loaders      []loader.Loader
	reader       data.Reader
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

func (p *defaultPipeline) AddLoader(t loader.Loader) Pipeline {
	p.loaders = append(p.loaders, t)
	return p
}

func (p *defaultPipeline) Start(r data.Reader) error {
	if p.extractor == nil || len(p.loaders) == 0 {
		return errors.New("extractor or loader is nil")
	}
	p.reader = r
	if r != nil {
		go func(p *defaultPipeline) {
			for {
				select {
				case data := <-p.reader.Data():
					p.PushData(data)
				case <-p.quitch:
					return
				}
			}
		}(p)
	}
	return nil
}

func (p *defaultPipeline) PushData(data interface{}) error {
	if p.extractor == nil || len(p.loaders) == 0 {
		return errors.New("extractor or loader is nil")
	}
	// extract data
	if frame, err := p.extractor.Extract(data); err == nil {
		// transfom data
		for _, transformer := range p.transformers {
			if f, err := transformer.Transform(frame); err != nil {
				return err
			} else {
				frame = f
			}
		}
		// load data
		for _, loader := range p.loaders {
			loader.Load(frame)
		}
		return nil
	}
	return errors.New("extract data failed")
}

func (p *defaultPipeline) Close() {
	if p.reader != nil {
		p.quitch <- true
	}
	p.extractor.Close()
	for _, trans := range p.transformers {
		trans.Close()
	}
	for _, loader := range p.loaders {
		loader.Close()
	}
	close(p.quitch)
}
