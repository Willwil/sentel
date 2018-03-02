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

package shiro

import "errors"

type Resource interface {
	GetName() string
	ResolveName(path string, ctx ResourceContext) string
	GetPermissions() []Permission
}

type Realm interface {
	GetName() string
	AddResource(Resource)
	RemoveResource(Resource)
	GetResource(path string, ctx ResourceContext) Resource
	Supports(token AuthenticationToken) bool
}

func NewRealm(env Environment, realmName string) (Realm, error) {
	return nil, errors.New("not implemented")
}
