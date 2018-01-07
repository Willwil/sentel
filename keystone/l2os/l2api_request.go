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

type createAccountReq struct {
	name string
	resp chan error
}
type destroyAccountReq struct {
	name string
	resp chan error
}
type createObjectReq struct {
	object *Object
	resp   chan error
}

type destroyObjectReq struct {
	objid string
	resp  chan error
}

type getObjectReq struct {
	objid string
	resp  chan *Object
}

type setObjectAttrReq struct {
	objname string
	key     string
	value   interface{}
	resp    chan error
}

type getObjectAttrReq struct {
	objid string
	key   string
	resp  chan interface{}
}

type updateObjectReq struct {
	object *Object
	resp   chan error
}
