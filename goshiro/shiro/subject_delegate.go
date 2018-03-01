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

type delegateSubject struct {
	securityMgr            SecurityManager
	principals             PrincipalCollection
	session                Session
	sessionCreationEnabled bool
	authenticated          bool
	host                   string
	authenticator          Authenticator
	authorizer             Authorizer
}

func (p *delegateSubject) hasPrincipals() bool { return p.principals.IsEmpty() != true }
func (p *delegateSubject) GetHost() string     { return p.host }
func (p *delegateSubject) GetPrincipal() Principal {
	return p.principals.GetPrimaryPrincipal()
}
func (p *delegateSubject) GetPrincipals() PrincipalCollection {
	return p.principals
}
func (p *delegateSubject) IsPermitted(permission string) bool {
	return p.hasPrincipals() && p.authorizer.IsPermitted(p.principals, permission)
}
func (p *delegateSubject) IsPermittedWithPermission(permission Permission) bool {
	return p.hasPrincipals() && p.authorizer.IsPermittedWithPermission(p.principals, permission)
}
func (p *delegateSubject) IsPermittedWithPermissions(permission ...string) []bool {
	result := []bool{}
	return result
}
func (p *delegateSubject) IsPermittedWithPermissionList(permissions []Permission) []bool {
	result := []bool{}
	return result
}
func (p *delegateSubject) IsPermittedAll(permissions ...string) bool {
	return true
}
func (p *delegateSubject) CheckPermission(permission Permission) error {
	return nil
}
func (p *delegateSubject) CheckPermissions(permissions []Permission) error {
	return nil
}
func (p *delegateSubject) CheckPermissionList(permission ...string) error {
	return nil
}
func (p *delegateSubject) HasRole(id string) bool {
	return true
}
func (p *delegateSubject) HasRoles(ids []string) []bool {
	result := []bool{}
	return result
}
func (p *delegateSubject) HasAllRoles(ids []string) bool {
	return true
}
func (p *delegateSubject) CheckRole(id string) error {
	return nil
}
func (p *delegateSubject) CheckRoles(id ...string) error {
	return nil
}
func (p *delegateSubject) CheckRolesWithList(ids []string) error {
	return nil
}
func (p *delegateSubject) Login(token AuthenticationToken) error {
	return nil
}
func (p *delegateSubject) IsAuthenticated() bool {
	return true
}
func (p *delegateSubject) GetSession() Session {
	return nil
}
func (p *delegateSubject) GetSessionWithCreation(create bool) Session {
	return nil
}
func (p *delegateSubject) Save() {}
