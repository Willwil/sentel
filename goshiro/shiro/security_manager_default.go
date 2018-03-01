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

type defaultSecurityManager struct {
	env           Environment
	realms        []Realm
	authenticator Authenticator
	authorizer    Authorizer
	sessionMgr    SessionManager
	cacheMgr      CacheManager
}

func (p *defaultSecurityManager) Login(subject Subject, token AuthenticationToken) (Subject, error) {
	return nil, errors.New("not implemented")
}
func (p *defaultSecurityManager) Logout(subject Subject) error {
	return nil
}
func (p *defaultSecurityManager) CreateSubject(ctx SubjectContext) (Subject, error) {
	return nil, nil
}
func (p *defaultSecurityManager) GetSubject(token AuthenticationToken) (Subject, error) {
	return nil, nil
}
func (p *defaultSecurityManager) Save(subject Subject)                 {}
func (p *defaultSecurityManager) SetCacheManager(mgr CacheManager)     { p.cacheMgr = mgr }
func (p *defaultSecurityManager) SetSessionManager(mgr SessionManager) { p.sessionMgr = mgr }
func (p *defaultSecurityManager) GetAuthenticator() Authenticator      { return p.authenticator }
func (p *defaultSecurityManager) SetAuthenticator(authc Authenticator) { p.authenticator = authc }
func (p *defaultSecurityManager) GetAuthorizer() Authorizer            { return p.authorizer }
func (p *defaultSecurityManager) SetAuthorizer(authz Authorizer)       { p.authorizer = authz }
func (p *defaultSecurityManager) GetRealm(realmName string) Realm      { return nil }
