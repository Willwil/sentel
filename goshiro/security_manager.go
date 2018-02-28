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

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

type securityManager interface {
	login(subject Subject, token AuthenticationToken) error
	logout(subject Subject) error
	createSubject(ctx SubjectContext) (Subject, error)
	getSubject(token AuthenticationToken) (Subject, error)
	save(subject Subject)
	setCacheManager(mgr CacheManager)
	setSessionManager(mgr SessionManager)
	getResourceManager() ResourceManager
	getAuthenticator() Authenticator
	setAuthenticator(Authenticator)
	getAuthorizer() Authorizer
	setAuthorizer(Authorizer)
}

func newSecurityManager(config.Config) (securityManager, error) {
	return nil, errors.New("no implemented")
}
