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

import "github.com/cloustone/sentel/pkg/config"

type directApi struct {
	config  config.Config
	storage objectStorage
}

func newDirectApi(c config.Config) (*directApi, error) {
	s, err := newStorage(c)
	if err != nil {
		return nil, err
	}
	return &directApi{config: c, storage: s}, nil
}

func (p *directApi) CreateObject(obj *Object) error {
	return p.storage.createObject(obj)
}
func (p *directApi) DestroyObject(objid string) error {
	return p.storage.deleteObject(objid)
}

// GetObject retieve object by object identifier
func (p *directApi) GetObject(objid string) (*Object, error) {
	return p.storage.getObject(objid)
}

func (p *directApi) UpdateObject(obj *Object) error {
	return p.storage.updateObject(obj)
}

func (p *directApi) CreateAccount(name string) error {
	return p.storage.createAccount(name)

}
func (p *directApi) DestroyAccount(name string) error {
	return p.storage.destroyAccount(name)
}
