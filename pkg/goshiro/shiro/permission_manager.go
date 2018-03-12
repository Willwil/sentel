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

// PermissionManager manage all permissions with permission identifer and
// which principal have what permissions
type PermissionManager interface {
	// AddPermission add a new permission to manager
	AddPermission(permission Permission) string
	// RemovePermission remove a existed permission from manager
	RemovePermission(permid string)
	// GetPermission return specified permission from manager
	GetPermission(permid string) (Permission, error)
	// RemovePrincipal remove principal from security manager
	RemovePrincipal(Principal)
	// GetPriincipalsPermissions return principals permissions
	GetPrincipalsPermissions(principals PrincipalCollection) []Permission
	// AddPrincipalPermissions add permissions for principals
	AddPrincipalPermissions(principal Principal, permissions []Permission)
	// RemovePrincipalsPermissions remove principals's permissions
	RemovePrincipalPermissions(principal Principal, permissions []Permission)
}

func NewPermissionManager(adaptor Adaptor) PermissionManager {
	return nil
}
