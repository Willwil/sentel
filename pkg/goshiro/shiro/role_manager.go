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

import "fmt"

// RoleManager managel all roles in shiro, it will do the following things
// 1. Which roles does the principal have
// 2. Which permissions does the role have
type RoleManager interface {
	// AddRole add role with permission into realm
	AddRole(Role)
	// RemoveRole remove specified role from realm
	RemoveRole(roleName string)
	// RemovePrincipalRoles remove principal's roles
	RemovePrincipalRoles(Principal)
	// GetRole return role's detail information
	GetRole(roleName string) (Role, error)
	// AddRolePermission add permissions to a role
	AddRolePermissions(roleName string, permissions []string)
	// RemoveRolePermissions remove permission from role
	RemoveRolePermissions(roleName string, permissions []string)
	// GetRolePermission return specfied role's all permissions
	GetRolePermissions(roleName string) []string
	// GetRolePermission return specfied role's all permissions identifier
	GetRolesPermissions(roleNames []string) []string
	// GetPrincipalRoles return principal's roles
	GetPrincipalRoles(Principal) []Role
	// GetPrincipalRoles return principal's role names
	GetPrincipalRoleNames(Principal) []string
}

func NewRoleManager(adaptor Adaptor, cacheManager CacheManager) RoleManager {
	return &roleManager{adaptor: adaptor, cacheManager: cacheManager}
}

type roleManager struct {
	adaptor      Adaptor
	cacheManager CacheManager
}

// AddRole add role with permission into realm
func (r *roleManager) AddRole(role Role) {
	r.adaptor.AddRole(role)
}

// RemoveRole remove specified role from realm
func (r *roleManager) RemoveRole(roleName string) {
	r.adaptor.RemoveRole(roleName)
}

// GetRole return role's detail information
func (r *roleManager) GetRole(roleName string) (Role, error) {
	if role, err := r.adaptor.GetRole(roleName); err == nil {
		return role, nil
	}
	return Role{}, fmt.Errorf("invalid role name '%s'", roleName)
}

// AddRolePermission add permissions to a role
func (r *roleManager) AddRolePermissions(roleName string, permissions []string) {
	r.adaptor.AddRolePermissions(roleName, permissions)
}

// RemoveRolePermissions remove permission from role
func (r *roleManager) RemoveRolePermissions(roleName string, permissions []string) {
	r.adaptor.RemoveRolePermissions(roleName, permissions)
}

// GetRolePermissions return specfied role's all permissions
func (r *roleManager) GetRolePermissions(roleName string) []string {
	return r.adaptor.GetRolePermissions(roleName)
}

// GetRolesPermissions return specfied role's all permissions
func (r *roleManager) GetRolesPermissions(roleNames []string) []string {
	permissions := []string{}
	for _, roleName := range roleNames {
		perms := r.GetRolePermissions(roleName)
		permissions = append(permissions, perms...)
	}
	return permissions
}

func (r *roleManager) GetPrincipalRoles(principal Principal) []Role {
	return r.adaptor.GetPrincipalRoles(principal)
}

func (r *roleManager) GetPrincipalRoleNames(principal Principal) []string {
	return r.adaptor.GetPrincipalRoleNames(principal)
}

func (r *roleManager) RemovePrincipalRoles(principal Principal) {
	r.adaptor.RemovePrincipalRoles(principal)
}
