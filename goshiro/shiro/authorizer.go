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

type Authorizer interface {
	IsPermitted(principals PrincipalCollection, permission string) bool
	IsPermittedWithPermission(principals PrincipalCollection, permission Permission) bool
	IsPermittedWithPermissions(principals PrincipalCollection, permission ...string) []bool
	IsPermittedWithPermissionList(principals PrincipalCollection, permissions []Permission) []bool
	IsPermittedAll(principals PrincipalCollection, permissions ...string) bool
	CheckPermission(principals PrincipalCollection, permission Permission) error
	CheckPermissions(principals PrincipalCollection, permissions []Permission) error
	CheckPermissionList(principals PrincipalCollection, permission ...string) error
	HasRole(pricipals PrincipalCollection, id string) bool
	HasRoles(principals PrincipalCollection, ids []string) []bool
	HasAllRoles(principals PrincipalCollection, ids []string) bool
	CheckRole(pricipals PrincipalCollection, id string) error
	CheckRoles(principals PrincipalCollection, id ...string) error
	CheckRolesWithList(principals PrincipalCollection, ids []string) error
}

func newAuthorizer(env Environment) Authorizer {
	return nil
}
