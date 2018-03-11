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

// Adaptor is wraper of policy and roles peristence
type Adaptor interface {
	// GetName return adaptor's name
	GetName() string
	// Role
	// AddRole add role with permission into realm
	AddRole(r Role)
	// RemoveRole remove specified role from realm
	RemoveRole(roleName string)
	// GetRole return role's detail information
	GetRole(roleName string) (Role, error)
	// AddRolePermission add permissions to a role
	AddRolePermissions(roleName string, permissons []Permission)
	// RemoveRolePermissions remove permission from role
	RemoveRolePermissions(roleName string, permissions []Permission)
	// GetRolePermission return specfied role's all permissions
	GetRolePermissions(roleName string) []Permission

	// Resource
	AddResource(Resource)
	RemoveResource(Resource)
	GetResource(resourceName string) *Resource
	GetResources() []Resource
}
