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

package subject

type Subject interface {
	GetPrincipal() interface{}
	GetPrincipals() []interface{}
	IsPermitted(permission string) bool
	IsPermittedWithPermission(permission Permission) bool
	IsPermittedWithPermissions(permission ...string) []bool
	IsPermittedWithPermissionList(permissions []Permission) []bool
	IsPermittedAll(permissions ...string) bool
	CheckPermission(permission Permission) error
	CheckPermissions(permissions []Permission) error
	CheckPermissionList(permission ...string) error
	HasRole(id string) bool
	HasRoles(ids []string) []bool
	HasAllRoles(ids []string) bool
	CheckRole(id string) error
	CheckRoles(id ...string) error
	CheckRolesWithList(ids []string) error
	Login(token AuthenticationToken) error
	IsAuthenticated() bool
	IsRembered() bool
	GetSession() Session
	GetSessionWithCreation(create bool) Session
	RunAs(principals PricipalCollection) error
	IsRunAs() bool
	GetPreviousPrincipals() PrincipalCollection
	ReleaseRunAs() PrincipalCollection
}
