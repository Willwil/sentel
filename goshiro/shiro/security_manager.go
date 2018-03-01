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

type SecurityManager interface {
	Login(subject Subject, token AuthenticationToken) error
	Logout(subject Subject) error
	CreateSubject(ctx subjectContext) (Subject, error)
	GetSubject(token AuthenticationToken) (Subject, error)
	Save(subject Subject)
	SetCacheManager(mgr CacheManager)
	SetSessionManager(mgr SessionManager)
	GetAuthenticator() Authenticator
	SetAuthenticator(Authenticator)
	GetAuthorizer() Authorizer
	SetAuthorizer(Authorizer)
	GetRealm(realmName string) Realm
}

func NewSecurityManager(env Environment, realms []Realm) SecurityManager {
	return &defaultSecurityManager{env: env, realms: realms}
}