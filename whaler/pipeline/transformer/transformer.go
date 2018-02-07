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

package transformer

import (
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/whaler/pipeline/data"
)

type Transformer interface {
	Transform(f *data.DataFrame) (*data.DataFrame, error)
	Close()
}

func New(c config.Config, name string) Transformer {
	switch name {
	default:
		return &noTransformer{}
	}
}

type noTransformer struct{}

func (p *noTransformer) Transform(f *data.DataFrame) (*data.DataFrame, error) {
	return f, nil
}
func (p *noTransformer) Close() {}
