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
	// AddRealm add reams to security manager
	AddRealm(realm ...Realm)
	// GetRealm return specified realm by realm name
	GetRealm(realmName string) Realm
	// AddPolicies load authorization policies
	AddPolicies([]AuthorizePolicy)
	// RemovePolicy remove specified authorization policy
	RemovePolicy(AuthorizePolicy)
	// GetPolicy return specified authorization policy
	GetPolicy(path string, ctx RequestContext) (AuthorizePolicy, error)
	// GetAllPolicies return all authorization policy
	GetAllPolicies() []AuthorizePolicy
	// Login authenticate current user and return subject
	Login(AuthenticationToken) (Subject, error)
	// GetSubject return specified subject by authentication token
	GetSubject(token AuthenticationToken) (Subject, error)
	// Authorize check wether the subject request is authorized
	Authorize(token AuthenticationToken, req Request) error
	// SetAdaptor set persistence adaptor
	SetAdaptor(Adaptor)
	// GetAuthorizeInfo return subject's authorization info
	GetAuthorizationInfo(token AuthenticationToken) (AuthorizationInfo, error)
	// AddRole add role with permission into realm
	AddRole(r Role)
	// RemoveRole remove specified role from realm
	RemoveRole(roleName string)
	// GetRole return role's detail information
	GetRole(roleName string) *Role
	// AddRolePermission add permissions to a role
	AddRolePermissions(roleName string, permissons []Permission)
	// RemoveRolePermissions remove permission from role
	RemoveRolePermissions(roleName string, permissions []Permission)
	// GetRolePermission return specfied role's all permissions
	GetRolePermissions(roleName string) []Permission
}
