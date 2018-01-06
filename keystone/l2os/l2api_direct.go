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
	config config.Config
}

func (p *directApi) CreateObject(obj *Object) error {
	s := getStorage()
	return s.createObject(obj)
}
func (p *directApi) DestroyObject(objid string) error {
	s := getStorage()
	return s.deleteObject(objid)
}

// GetObject retieve object by object identifier
func (p *directApi) GetObject(objid string) (*Object, error) {
	s := getStorage()
	return s.getObject(objid)
}

func (p *directApi) UpdateObject(obj *Object) error {
	s := getStorage()
	return s.updateObject(obj)
}

func (p *directApi) CreateAccount(name string) error {
	s := getStorage()
	return s.createAccount(name)

}
func (p *directApi) DestroyAccount(name string) error {
	s := getStorage()
	return s.destroyAccount(name)
}
