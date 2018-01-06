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

package l2

import (
	"fmt"

	"github.com/cloustone/sentel/pkg/config"
)

type Api interface {
	CreateAccount(name string) error
	DestroyAccount(name string) error
	CreateObject(*Object) error
	DestroyObject(objid string) error
	GetObject(objid string) (*Object, error)
	UpdateObject(*Object) error
}

func NewApi(name string, c config.Config) (Api, error) {
	switch name {
	case "default":
		return &defaultApi{config: c}, nil
	case "direct":
		return &directApi{config: c}, nil
	default:
		return nil, fmt.Errorf("invalid api '%s'", name)
	}
}
