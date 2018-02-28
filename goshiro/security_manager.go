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

package goshiro

type SecurityManager interface {
	Login(subject Subject, token AuthenticationToken) error
	Logout(subject Subject) error
	CreateSubject(ctx subjectContext) (Subject, error)
	GetSubject(token AuthenticationToken) (Subject, error)
	Save(subject Subject)
	SetCacheManager(mgr CacheManager)
	SetSessionManager(mgr SessionManager)
	GetResourceManager() ResourceManager
	SetResourceManager(ResourceManager)
	GetAuthenticator() Authenticator
	SetAuthenticator(Authenticator)
	GetAuthorizer() Authorizer
	SetAuthorizer(Authorizer)
}

type SecurityManagerFactory struct {
	env Environment
}

func newSecurityManager(env Environment) SecurityManager {
	return &defaultSecurityManager{}
}

func NewSecurityManagerFactory(env Environment) SecurityManagerFactory {
	return SecurityManagerFactory{env: env}
}

// GetInstance return security manager instance using environment
func (p *SecurityManagerFactory) GetInstance() SecurityManager {
	securitymgr := newSecurityManager(p.env)
	securitymgr.SetAuthenticator(newAuthenticator(p.env))
	securitymgr.SetAuthorizer(newAuthorizer(p.env))
	securitymgr.SetResourceManager(newResourceManager(p.env))
	securitymgr.SetSessionManager(newSessionManager(p.env))
	securitymgr.SetCacheManager(newCacheManager(p.env))
	return securitymgr
}
