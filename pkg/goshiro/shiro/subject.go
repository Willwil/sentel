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

type Subject interface {
	GetPrincipal() Principal
	GetPrincipals() PrincipalCollection
	HasRole(id string) bool
	HasRoles(ids []string) []bool
	HasAllRoles(ids []string) bool
	IsAuthenticated() bool
	GetSession() Session
	SetSessionCreationEnabled(create bool)
	IsSessionCreationEnabled() bool
	GetHost() string
}

func NewSubject(ctx SubjectContext) Subject {
	return nil
}

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

func (p *delegateSubject) HasRole(id string) bool {
	return false
}

func (p *delegateSubject) HasRoles(ids []string) []bool {
	return boolSliceWithCapacity(len(ids), false)
}

func (p *delegateSubject) HasAllRoles(ids []string) bool {
	return false
}

func (p *delegateSubject) IsAuthenticated() bool { return p.authenticated }

func (p *delegateSubject) GetSession() Session            { return p.session }
func (p *delegateSubject) IsSessionCreationEnabled() bool { return false }

func (p *delegateSubject) SetSessionCreationEnabled(create bool) {}
