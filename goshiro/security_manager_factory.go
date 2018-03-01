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

import "github.com/cloustone/sentel/goshiro/shiro"

type SecurityManagerFactory struct {
	env          shiro.Environment
	realmFactory *RealmFactory
}

func NewSecurityManagerFactory(env shiro.Environment, realmFactory *RealmFactory) SecurityManagerFactory {
	if realmFactory == nil {
		realmFactory = NewRealmFactory(env)
	}
	return SecurityManagerFactory{env: env, realmFactory: realmFactory}
}

// GetInstance return security manager instance using environment
func (p *SecurityManagerFactory) GetInstance() shiro.SecurityManager {
	realms := p.realmFactory.GetRealms()
	securitymgr := shiro.NewSecurityManager(p.env, realms)
	securitymgr.SetAuthenticator(shiro.NewAuthenticator(p.env))
	securitymgr.SetAuthorizer(shiro.NewAuthorizer(p.env))
	securitymgr.SetSessionManager(shiro.NewSessionManager(p.env))
	securitymgr.SetCacheManager(shiro.NewCacheManager(p.env))
	return securitymgr
}
