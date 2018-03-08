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

type AuthorizationInfo interface {
	GetRoles() []Role
	GetStringPermissions() []string
	GetObjectPermissions() []interface{}
}

func NewAuthorizationInfo(roles []Role) AuthorizationInfo {
	info := &simpleAuthorizationInfo{
		roles:             []Role{},
		stringPermissions: []string{},
		objectPermissions: []interface{}{},
	}
	info.AddRoles(roles)
	return info
}

type simpleAuthorizationInfo struct {
	roles             []Role
	stringPermissions []string
	objectPermissions []interface{}
}

func (p *simpleAuthorizationInfo) GetRoles() []Role {
	return p.roles
}
func (p *simpleAuthorizationInfo) GetStringPermissions() []string {
	return p.stringPermissions
}
func (p *simpleAuthorizationInfo) GetObjectPermissions() []interface{} {
	return p.objectPermissions
}

func (p *simpleAuthorizationInfo) SetRoles(roles []Role) {
	p.roles = roles
}

func (p *simpleAuthorizationInfo) AddRole(role Role) {
	p.roles = append(p.roles, role)
}
func (p *simpleAuthorizationInfo) AddRoles(roles []Role) {
	p.roles = append(p.roles, roles...)
}

func (p *simpleAuthorizationInfo) GetAuthorizationInfo() AuthorizationInfo {
	return nil
}
