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

package goshiro

import "fmt"

type Environment interface {
	SetValue(key string, val interface{})
	GetValue(key string) (interface{}, error)
}

func NewEnvironment() Environment {
	return &environment{
		values: make(map[string]interface{}),
	}
}

type environment struct {
	values map[string]interface{}
}

func (p *environment) SetValue(key string, val interface{}) {
	p.values[key] = val
}
func (p environment) GetValue(key string) (interface{}, error) {
	if val, found := p.values[key]; found {
		return val, nil
	}
	return nil, fmt.Errorf("'%s' not found in environment", key)
}
