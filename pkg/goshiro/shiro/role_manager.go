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

import "errors"

// RoleManager managel all roles in shiro
type RoleManager struct {
	adaptor Adaptor
}

func NewRoleManager(adaptor Adaptor) *RoleManager {
	return &RoleManager{adaptor: adaptor}
}

// AddRole add role with permission into realm
func (r *RoleManager) AddRole(role Role) {}

// RemoveRole remove specified role from realm
func (r *RoleManager) RemoveRole(roleName string) {}

// GetRole return role's detail information
func (r *RoleManager) GetRole(roleName string) (*Role, error) {
	return nil, errors.New("not implemented")
}

// AddRolePermission add permissions to a role
func (r *RoleManager) AddRolePermissions(roleName string, permissons []Permission) {}

// RemoveRolePermissions remove permission from role
func (r *RoleManager) RemoveRolePermissions(roleName string, permissions []Permission) {}

// GetRolePermission return specfied role's all permissions
func (r *RoleManager) GetRolePermissions(roleName string) []Permission {
	return []Permission{}
}
