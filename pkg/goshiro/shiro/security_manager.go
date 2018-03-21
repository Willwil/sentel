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

import "github.com/cloustone/sentel/pkg/config"

type SecurityManager interface {
	// Load load security metadata from backend if execption occurred
	Load()
	// AddRealm add reams to security manager
	AddRealm(realm ...Realm)
	// GetRealm return specified realm by realm name
	GetRealm(realmName string) Realm
	// Login authenticate current user and return subject
	Login(AuthenticationToken) (Principal, error)
	// Authorize check wether the subject request is authorized
	Authorize(principals PrincipalCollection, resource, action string) error
	// AuthorizeWithToken check wether the subject request is authorized
	AuthorizeWithToken(token AuthenticationToken, resource, action string) error
	// RemovePrincipal remove principal from security manager
	RemovePrincipal(Principal)
	// GetPriincipalsPermissions return principals permissions
	GetPrincipalsPermissions(principals PrincipalCollection) []Permission
	// AddPrincipalPermissions add permissions for principals
	AddPrincipalPermissions(principal Principal, permissions []Permission)
	// RemovePrincipalsPermissions remove principals's permissions
	RemovePrincipalPermissions(principal Principal, permissions []Permission)
	// AddRole add role with permission into realm
	AddRole(r Role)
	// RemoveRole remove specified role from realm
	RemoveRole(roleName string)
	// GetRole return role's detail information
	GetRole(roleName string) (Role, error)
	// AddRolePermission add permissions to a role
	AddRolePermissions(roleName string, permissions []Permission)
	// RemoveRoleAllPermission remove role's all permissions
	RemoveRoleAllPermissions(roleName string)
	// RemoveRolePermissions remove permission from role
	RemoveRolePermissions(roleName string, permissions []Permission)
	// GetRolePermission return specfied role's all permissions
	GetRolePermissions(roleName string) []Permission
}

func NewSecurityManager(c config.Config, realm ...Realm) (SecurityManager, error) {
	adaptor, err := NewAdaptor(c)
	if err != nil {
		return nil, err
	}
	return newDefaultSecurityManager(c, adaptor, realm...)
}
