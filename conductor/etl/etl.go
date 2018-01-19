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

package ETL

import (
	"errors"

	"github.com/cloustone/sentel/conductor/data"
	"github.com/cloustone/sentel/conductor/extractor"
	"github.com/cloustone/sentel/conductor/loader"
	"github.com/cloustone/sentel/conductor/transformer"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type ETL interface {
	AddExtractor(extractor.Extractor) ETL
	AddTransformer(transformer.Transformer) ETL
	AddLoader(loader.Loader)
	Run(r data.Reader, ctx interface{}) error
	Close()
}

type etl struct {
	config       config.Config
	extractor    extractor.Extractor
	transformers []transformer.Transformer
	loader       loader.Loader
}

func New(c config.Config) ETL {
	return &etl{
		config:       c,
		transformers: []transformer.Transformer{},
	}
}

func (p *etl) AddExtractor(t extractor.Extractor) ETL {
	if p.extractor != nil {
		glog.Error("extractor already exist")
	}
	p.extractor = t
	return p
}

func (p *etl) AddTransformer(t transformer.Transformer) ETL {
	p.transformers = append(p.transformers, t)
	return p
}
func (p *etl) AddLoader(t loader.Loader) {
	if p.loader != nil {
		glog.Error("loader already exist")
	}
	p.loader = t
}

func (p *etl) Run(r data.Reader, ctx interface{}) error {
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

func (p *etl) Close() {
	p.extractor.Close()
	for _, trans := range p.transformers {
		trans.Close()
	}
	if p.loader != nil {
		p.loader.Close()
	}
}
