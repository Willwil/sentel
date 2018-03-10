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
	GetPrincipals() PrincipalCollection
	AddPrincipal(principal Principal, realmName string)
	GetRoles() []Role
	AddRole(Role)
	AddRoles([]Role)
}

func NewAuthorizationInfo(principals PrincipalCollection) AuthorizationInfo {
	info := &simpleAuthorizationInfo{
		principals: principals,
		roles:      []Role{},
	}
	return info
}

type simpleAuthorizationInfo struct {
	principals PrincipalCollection
	roles      []Role
}

func (p *simpleAuthorizationInfo) AddRole(r Role)    { p.roles = append(p.roles, r) }
func (p *simpleAuthorizationInfo) AddRoles(r []Role) { p.roles = append(p.roles, r...) }
func (p *simpleAuthorizationInfo) GetRoles() []Role  { return p.roles }
func (p *simpleAuthorizationInfo) AddPrincipal(principal Principal, realmName string) {
	p.principals.Add(principal, realmName)
}
func (p *simpleAuthorizationInfo) GetPrincipals() PrincipalCollection { return p.principals }
