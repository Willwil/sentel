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

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

type DefaultSecurityManager struct {
}

func NewSecurityManager(c config.Config) (*DefaultSecurityManager, error) {
	return &DefaultSecurityManager{}, nil
}

func (w *DefaultSecurityManager) AddRealm(r Realm) {}
func (w *DefaultSecurityManager) Login(AuthenticationToken) (Subject, error) {
	return nil, errors.New("not implemented")
}
func (w *DefaultSecurityManager) Logout(subject Subject) error {
	return nil
}
func (w *DefaultSecurityManager) GetSubject(token AuthenticationToken) (Subject, error) {
	return nil, errors.New("not implmented")
}
func (w *DefaultSecurityManager) Save(subject Subject) {}
func (w *DefaultSecurityManager) CheckPermission(subject Subject, resourceName string, action string) error {
	return nil
}

func (w *DefaultSecurityManager) Authorize(subject Subject, req Request) error {
	return nil
}

func (w *DefaultSecurityManager) SetAdaptor(a Adaptor)                 {}
func (w *DefaultSecurityManager) SetAuthenticator(authc Authenticator) {}

func (w *DefaultSecurityManager) AddPolicies([]AuthorizePolicy) {}
func (w *DefaultSecurityManager) RemovePolicy(AuthorizePolicy)  {}
func (w *DefaultSecurityManager) GetPolicy(path string, ctx RequestContext) AuthorizePolicy {
	return AuthorizePolicy{}
}
func (w *DefaultSecurityManager) GetAllPolicies() []AuthorizePolicy { return nil }
func (w *DefaultSecurityManager) LoadPolicies([]AuthorizePolicy)    {}

func NewDefaultSecurityManager(c config.Config) (SecurityManager, error) {
	return nil, errors.New("not implemented")
}
