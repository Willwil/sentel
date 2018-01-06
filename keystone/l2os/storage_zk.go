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
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

type zkStorage struct {
	config config.Config
}

func newZkStorage(c config.Config) (*zkStorage, error) {
	return nil, errors.New("not implemented")
}

func (p *zkStorage) createObject(*Object) error {
	return nil
}
func (p *zkStorage) deleteObject(objid string) error {
	return nil
}
func (p *zkStorage) getObject(objid string) (*Object, error) {
	return nil, errors.New("not implemented")
}

func (p *zkStorage) updateObject(obj *Object) error {
	return errors.New("not implemented")
}

func (p *zkStorage) deinitialize() {
}

func (p *zkStorage) createAccount(name string) error {
	return nil
}
func (p *zkStorage) destroyAccount(name string) error {
	return nil
}
