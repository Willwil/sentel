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

// roleManager managel all roles in shiro
type RoleManager interface {
	AddRole(Role)
	RemoveRole(roleName string)
	GetRole(roleName string) *Role
	AddRolePermissions(roleName string, permissions []Permission)
	RemoveRolePermissions(roleName string, permissions []Permission)
	GetRolePermissions(roleName string) []Permission
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

// GetRolePermission return specfied role's all permissions
func (r *roleManager) GetRolePermissions(roleName string) []Permission {
	if r.adaptor != nil {
		return r.adaptor.GetRolePermissions(roleName)
	}
	return []Permission{}
}
