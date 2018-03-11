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

type Subject interface{}

type SecurityManager interface {
	// AddRealm add reams to security manager
	AddRealm(realm ...Realm)
	// GetRealm return specified realm by realm name
	GetRealm(realmName string) Realm
	// Login authenticate current user and return subject
	Login(AuthenticationToken) (Subject, error)
	// Authorize check wether the subject request is authorized
	Authorize(token AuthenticationToken, resource, action string) error
	GetPrincipalsPermissions(principals PrincipalCollection) []Permission
}
