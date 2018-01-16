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

package etl

import (
	"errors"

	"github.com/cloustone/sentel/conductor/etl/data"
	"github.com/cloustone/sentel/conductor/etl/extractor"
	"github.com/cloustone/sentel/conductor/etl/loader"
	"github.com/cloustone/sentel/conductor/etl/transformer"
	"github.com/cloustone/sentel/pkg/config"
)

type ETL struct {
	config       config.Config
	extractor    extractor.Extractor
	transformers []transformer.Transformer
	loader       loader.Loader
}

func New(c config.Config) (*ETL, error) {
	extractor, _ := extractor.New(c)
	transformers := transformer.New(c)
	loader, _ := loader.New(c)
	if extractor == nil || len(transformers) == 0 || loader == nil {
		return nil, errors.New("etl initialization failed")
	}
	return &ETL{
		config:       c,
		extractor:    extractor,
		transformers: transformers,
		loader:       loader,
	}, nil
}

func (p *ETL) Run(r data.Reader, ctx interface{}) error {
	if data, err := p.extractor.Extract(r, ctx); err != nil {
		for _, trans := range p.transformers {
			result, err := trans.Transform(data, ctx)
			if err != nil {
				return err
			}
			data = result
		}
		if err := p.loader.Load(data, ctx); err != nil {
			return err
		}
	}
	return nil
}

func (p *ETL) Close() {
	p.extractor.Close()
	for _, trans := range p.transformers {
		trans.Close()
	}
	p.loader.Close()
}
