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
	"fmt"

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

func NewPipeline(c config.Config, name string) Pipeline {
	return &defaultPipeline{
		config:       c,
		name:         name,
		transformers: []transformer.Transformer{},
		loaders:      []loader.Loader{},
		quitch:       make(chan interface{}),
	}
}

type defaultPipeline struct {
	config       config.Config             // Pipeline configuration
	name         string                    // Pipeline name
	extractor    extractor.Extractor       // Pipeline only have one extractor
	transformers []transformer.Transformer // Many transformers
	loaders      []loader.Loader           // Many laoders
	reader       data.Reader               // Data reader if pipeline is asynchrous
	quitch       chan interface{}          // Quit channel if pipeline is asynchrouns
}

func (p *defaultPipeline) AddExtractor(t extractor.Extractor) Pipeline {
	if p.extractor != nil {
		glog.Errorf("pipeline '%s' extractor already exist", p.name)
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
		return fmt.Errorf("pipeline '%s' extractor or loader is nil", p.name)
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
		return fmt.Errorf("pipeline '%s' extractor or loader is nil", p.name)
	}
	// extract data
	glog.Infof("pipeline '%s' is processing data", p.name)
	if frame, err := p.extractor.Extract(data); err == nil {
		// transfom data
		for _, transformer := range p.transformers {
			if f, err := transformer.Transform(frame); err != nil {
				return fmt.Errorf("pipeline '%s' transform data failed, %s", p.name, err.Error())
			} else {
				frame = f
			}
		}
		// load data
		for _, loader := range p.loaders {
			if err := loader.Load(frame); err != nil {
				return fmt.Errorf("pipeline '%s' loader '%s' failed, %s", p.name, loader.Name(), err.Error())
			}
		}
		return nil
	}
	return fmt.Errorf("pipeline '%s' extract data failed", p.name)
}

func (p *defaultPipeline) Close() {
	glog.Infof("pipeline '%s' closed")
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
