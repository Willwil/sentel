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

type AuthenticationToken interface {
	GetPrincipal() interface{}
	GetCrenditals() interface{}
}

type AuthenticationInfo interface {
	GetPrincipals() PrincipalCollection
	GetCrenditals() interface{}
}

type AuthenticationListener interface{}

// Authenticator is responsible for authenticating accounts in an application
type Authenticator interface {
	Authenticate(AuthenticationToken) (AuthenticationInfo, error)
	SetAuthenticationListeners([]AuthenticationListener)
	GetAuthenticationListeners() []AuthenticationListener
}

type defaultAuthenticator struct {
	listeners map[string]AuthenticationListener
}

func NewAuthenticator(env Environment) Authenticator {
	return nil
}
