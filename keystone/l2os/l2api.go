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
	"github.com/golang/glog"
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
	var api Api
	var err error

	switch name {
	case "default":
		api, err = newDefaultApi(c)
	case "direct":
		api, err = newDirectApi(c)
	default:
		glog.Errorf("invalid l2api '%s'", name)
	}
	if err != nil || api == nil {
		return &nilApi{}, errors.New("l2api initialize failed, using nilapi")
	}
	return api, nil
}

type nilApi struct{}

func (p *nilApi) CreateAccount(name string) error {
	return errors.New("l2api is null")
}
func (p *nilApi) DestroyAccount(name string) error {
	return errors.New("l2api is null")
}
func (p *nilApi) CreateObject(*Object) error {
	return errors.New("l2api is null")
}
func (p *nilApi) DestroyObject(objid string) error {
	return errors.New("l2api is null")
}
func (p *nilApi) GetObject(objid string) (*Object, error) {
	return nil, errors.New("l2api is null")
}
func (p *nilApi) UpdateObject(*Object) error {
	return errors.New("l2api is null")
}
