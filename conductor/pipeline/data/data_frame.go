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

package data

import "encoding/json"

type DataFrame struct {
	items map[string]interface{}
}

func NewDataFrame() *DataFrame {
	return &DataFrame{
		items: make(map[string]interface{}),
	}
}

func (p *DataFrame) AddField(k string, v interface{}) {
	p.items[k] = v
}

func (p *DataFrame) GetField(k string) interface{} {
	if v, found := p.items[k]; found {
		return v
	}
	return nil
}

func (p *DataFrame) Serialize() ([]byte, error) {
	return json.Marshal(p.items)
}
