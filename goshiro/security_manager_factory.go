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

type SecurityManagerFactory struct {
	env Environment
}

func NewSecurityManagerFactory(env Environment) SecurityManagerFactory {
	return SecurityManagerFactory{env: env}
}

// GetInstance return security manager instance using environment
func (p *SecurityManagerFactory) GetInstance() SecurityManager {
	authenticator := newAuthenticator(p.env)
	authorizer := newAuthorizer(p.env)
	resmgr := NewResourceManager(p.env)
	securitymgr := newSecurityManager(p.env)
	sessionmgr := newSessionManager(p.env)
	cachemgr := newCacheManager(p.env)

	securitymgr.SetAuthenticator(authenticator)
	securitymgr.SetAuthorizer(authorizer)
	securitymgr.SetResourceManager(resmgr)
	securitymgr.SetSessionManager(sessionmgr)
	securitymgr.SetCacheManager(cachemgr)
	return securitymgr
}
