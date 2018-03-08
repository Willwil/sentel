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

// Realm is collection of policy
type Realm interface {
	// GetName return realm's name
	GetName() string
	// Supports check wether the realm support the specified authentication token
	Supports(token AuthenticationToken) bool
	// GetPrincipal return subject's principals
	GetPrincipals(token AuthenticationToken) PrincipalCollection
	// GetAuthorizeInfo return subject's authorization info
	GetAuthorizeInfo(subject Subject) AuthorizationInfo
}

func NewRealm(c config.Config, realmName string) (Realm, error) {
	return nil, errors.New("not implemented")
}
