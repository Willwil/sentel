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

type defaultApi struct {
	config config.Config
}

func newDefaultApi(c config.Config) (*defaultApi, error) {
	return &defaultApi{config: c}, nil
}

func (p *defaultApi) CreateObject(obj *Object) error {
	req := createObjectReq{
		object: obj,
		resp:   make(chan error),
	}
	getService().sendRequest(req)
	return <-req.resp
}
func (p *defaultApi) DestroyObject(objid string) error {
	req := destroyObjectReq{
		objid: objid,
		resp:  make(chan error),
	}
	getService().sendRequest(req)
	return <-req.resp
}

// GetObject retieve object by object identifier
func (p *defaultApi) GetObject(objid string) (*Object, error) {
	req := getObjectReq{
		objid: objid,
		resp:  make(chan *Object),
	}
	getService().sendRequest(req)
	obj := <-req.resp
	if obj == nil {
		return nil, fmt.Errorf("object '%s' no exist", objid)
	}
	return obj, nil
}

func (p *defaultApi) UpdateObject(obj *Object) error {
	req := updateObjectReq{
		object: obj,
		resp:   make(chan error),
	}
	getService().sendRequest(req)
	return <-req.resp
}

func (p *defaultApi) CreateAccount(name string) error {
	req := createAccountReq{
		name: name,
		resp: make(chan error),
	}
	getService().sendRequest(req)
	return <-req.resp
}

func (p *defaultApi) DestroyAccount(name string) error {
	req := destroyAccountReq{
		name: name,
		resp: make(chan error),
	}
	getService().sendRequest(req)
	return <-req.resp
}
