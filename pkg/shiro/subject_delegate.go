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

type delegateSubject struct {
	securityMgr            SecurityManager
	principals             PrincipalCollection
	session                Session
	sessionCreationEnabled bool
	authenticated          bool
	host                   string
}

func boolSliceWithCapacity(number int, v bool) []bool {
	result := make([]bool, number)
	for i := 0; i < number; i++ {
		result[i] = v
	}
	return result
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

func (p *delegateSubject) IsPermittedWithPermissions(permissions []Permission) []bool {
	if p.hasPrincipals() {
		return p.authorizer.IsPermittedWithPermissions(p.principals, permissions)
	}
	return boolSliceWithCapacity(len(permissions), false)
}

func (p *delegateSubject) CheckPermission(permission Permission) error {
	if p.hasPrincipals() {
		return p.authorizer.CheckPermission(p.principals, permission)
	}
	return errors.New("no permission")
}

func (p *delegateSubject) CheckPermissions(permissions []Permission) error {
	if p.hasPrincipals() {
		return p.authorizer.CheckPermissions(p.principals, permissions)
	}
	return errors.New("no permission")
}

func (p *delegateSubject) HasRole(id string) bool {
	if p.hasPrincipals() {
		return p.authorizer.HasRole(p.principals, id)
	}
	return false
}

func (p *delegateSubject) HasRoles(ids []string) []bool {
	if p.hasPrincipals() {
		return p.authorizer.HasRoles(p.principals, ids)
	}
	return boolSliceWithCapacity(len(ids), false)
}

func (p *delegateSubject) HasAllRoles(ids []string) bool {
	if p.hasPrincipals() {
		return p.authorizer.HasAllRoles(p.principals, ids)
	}
	return false
}

func (p *delegateSubject) CheckRole(id string) error {
	if p.hasPrincipals() {
		return p.authorizer.CheckRole(p.principals, id)
	}
	return errors.New("no role")
}

func (p *delegateSubject) CheckRoles(ids []string) error {
	if p.hasPrincipals() {
		return p.authorizer.CheckRoles(p.principals, ids)
	}
	return errors.New("no roles")
}

func (p *delegateSubject) Login(token AuthenticationToken) error {
	subject, err := p.securityMgr.Login(p, token)
	if err != nil {
		return err
	}
	p.principals = subject.GetPrincipals()
	p.authenticated = true
	p.host = subject.GetHost()
	p.session = subject.GetSession()
	return nil
}

func (p *delegateSubject) IsAuthenticated() bool { return p.authenticated }

func (p *delegateSubject) GetSession() Session { return p.session }

func (p *delegateSubject) SetSessionCreationEnabled(create bool) { p.sessionCreationEnabled = create }
func (p *delegateSubject) IsSessionCreationEnabled() bool        { return p.sessionCreationEnabled }

func (p *delegateSubject) Save() {}
