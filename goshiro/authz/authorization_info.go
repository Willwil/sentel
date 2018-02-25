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

package authz

type AuthorizationInfo interface {
	GetRoles() []string
	GetStringPermissions() []string
	GetObjectPermissions() []interface{}
}

func NewAuthorizationInfo(roles []string) AuthorizationInfo {
	info := &simpleAuthorizationInfo{
		roles:             make([]string),
		stringPermissions: make([]string),
		objectPermissions: make([]interface{}),
	}
	info.AddRoles(roles)
	return info
}

type simpleAuthorizationInfo struct {
	roles             []string
	stringPermissions []string
	objectPermissions []interface{}
}

func (p *simpleAuthorizationInfo) GetRoles() []string {
	return p.roles
}
func (p *simpleAuthorizationInfo) GetStringPermissions() []string {
	return p.stringPermissions
}
func (p *simpleAuthorizationInfo) GetObjectPermissions() []interface{} {
	return p.objectPermissions
}

func (p *simpleAuthorizationInfo) SetRoles(roles []string) {
	p.roles = roles
}

func (p *simpleAuthorizationInfo) AddRole(role string) {
	p.roles = append(p.roles, role)
}
func (p *simpleAuthorizationInfo) AddRoles(roles []string) {
	p.roles = append(p.roles, roles...)
}
