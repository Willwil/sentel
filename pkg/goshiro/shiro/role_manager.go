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

// RoleManager managel all roles in shiro, it will do the following things
// 1. Which roles does the principal have
// 2. Which permissions does the role have
type RoleManager interface {
	// AddRole add role with permission into realm
	AddRole(Role)
	// RemoveRole remove specified role from realm
	RemoveRole(roleName string)
	// GetRole return role's detail information
	GetRole(roleName string) *Role
	// AddRolePermission add permissions to a role
	AddRolePermissions(roleName string, permissions []Permission)
	// RemoveRolePermissions remove permission from role
	RemoveRolePermissions(roleName string, permissions []Permission)
	// GetRolePermission return specfied role's all permissions
	GetRolePermissions(roleName string) []Permission
	// GetRolePermission return specfied role's all permissions
	GetRolesPermissions(roleNames []string) []Permission
	// GetPrincipalsRoles return principals roles
	GetPrincipalsRoles(PrincipalCollection) []Role
	// GetPrincipalRoles return principal's roles
	GetPrincipalRoles(Principal) []Role
	// GetPrincipalsRoles return principals role names
	GetPrincipalsRoleNames(PrincipalCollection) []string
	// GetPrincipalRoles return principal's role names
	GetPrincipalRoleNames(Principal) []string
	// GetPrincipalsPermissions return principals's permissions
	GetPrincipalsPermissions(PrincipalCollection) []Permission
}

func NewRoleManager(adaptor Adaptor) RoleManager {
	return &roleManager{adaptor: adaptor}
}

type roleManager struct {
	adaptor Adaptor
}

// AddRole add role with permission into realm
func (r *roleManager) AddRole(role Role) {
	if r.adaptor != nil {
		r.adaptor.AddRole(role)
	}
}

// RemoveRole remove specified role from realm
func (r *roleManager) RemoveRole(roleName string) {
	if r.adaptor != nil {
		r.adaptor.RemoveRole(roleName)
	}
}

// GetRole return role's detail information
func (r *roleManager) GetRole(roleName string) *Role {
	if r.adaptor != nil {
		if role, err := r.adaptor.GetRole(roleName); err == nil {
			return &role
		}
	}
	return nil
}

// AddRolePermission add permissions to a role
func (r *roleManager) AddRolePermissions(roleName string, permissions []Permission) {
	if r.adaptor != nil {
		r.adaptor.AddRolePermissions(roleName, permissions)
	}
}

// RemoveRolePermissions remove permission from role
func (r *roleManager) RemoveRolePermissions(roleName string, permissions []Permission) {
	if r.adaptor != nil {
		r.adaptor.RemoveRolePermissions(roleName, permissions)
	}
}

// GetRolePermissions return specfied role's all permissions
func (r *roleManager) GetRolePermissions(roleName string) []Permission {
	if r.adaptor != nil {
		return r.adaptor.GetRolePermissions(roleName)
	}
	return []Permission{}
}

// GetRolesPermissions return specfied role's all permissions
func (r *roleManager) GetRolesPermissions(roleNames []string) []Permission {
	return []Permission{}
}

func (r *roleManager) GetPrincipalsRoles(principals PrincipalCollection) []Role {
	return []Role{}
}

func (r *roleManager) GetPrincipalRoles(principal Principal) []Role {
	return []Role{}
}

func (r *roleManager) GetPrincipalsRoleNames(principals PrincipalCollection) []string {
	return []string{}
}

func (r *roleManager) GetPrincipalRoleNames(principal Principal) []string {
	return []string{}
}

func (r *roleManager) GetPrincipalsPermissions(principals PrincipalCollection) []Permission {
	return []Permission{}
}
