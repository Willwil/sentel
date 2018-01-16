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
	config      config.Config
	extractor   extractor.Extractor
	transformer transformer.Transformer
	loader      loader.Loader
}

func New(c config.Config) (*ETL, error) {
	extractor, _ := extractor.New(c)
	transformer, _ := transformer.New(c)
	loader, _ := loader.New(c)
	if extractor == nil || transformer == nil || loader == nil {
		return nil, errors.New("etl initialization failed")
	}
	return &ETL{
		config:      c,
		extractor:   extractor,
		transformer: transformer,
		loader:      loader,
	}, nil
}

func (p *ETL) Run(r data.Reader, ctx interface{}) error {
	if data, err := p.extractor.Extract(r, ctx); err != nil {
		if data, err := p.transformer.Transform(data, ctx); err != nil {
			if err := p.loader.Load(data, ctx); err != nil {
				return err
			}
			return err
		}
		return err
	}
	return nil
}
