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
	"strings"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

type Transformer interface {
	Transform(data map[string]interface{}, ctx interface{}) (map[string]interface{}, error)
}

func New(c config.Config) []Transformer {
	transformers := []Transformer{}
	if names, err := c.String("etl", "transformers"); err == nil {
		for _, name := range strings.Split(names, ",") {
			switch name {
			case "no":
				transformers = append(transformers, &noTransformer{})
			default:
				glog.Errorf("invalid transformer '%s'", name)
				break
			}
		}
	}
	return transformers
}

type noTransformer struct{}

func (p *noTransformer) Transform(data map[string]interface{}, ctx interface{}) (map[string]interface{}, error) {
	return data, nil
}
