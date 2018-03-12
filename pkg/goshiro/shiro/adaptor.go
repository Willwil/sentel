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

// Adaptor is wraper of policy and roles peristence
type Adaptor interface {
	GetName() string
	AddPermissions([]Permission) []string
	RemovePermissions(permIds []string)
	GetPermissions(permIds []string) []Permission
	// Principal
	SetPrincipalPermissions(principal Principal, permissions []Permission)
	RemovePrincipal(Principal)
	GetPrincipalPermissions(principal Principal) []Permission
	GetPrincipalsPermissions(principals PrincipalCollection) []Permission
	AddPrincipalPermissions(principal Principal, permissions []Permission)
	RemovePrincipalPermissions(principal Principal, permissions []Permission)
	SetPrincipalRole(principal Principal, roleName ...string) error
	GetPrincipalRoles(principal Principal) []Role
	GetPrincipalRoleNames(principal Principal) []string
	RemovePrincipalRoles(principal Principal)

	// Role
	AddRole(r Role)
	RemoveRole(roleName string)
	GetRole(roleName string) (Role, error)
	AddRolePermissions(roleName string, permissons []string)
	RemoveRolePermissions(roleName string, permissions []string)
	GetRolePermissions(roleName string) []string
}

func NewAdaptor(c config.Config) (Adaptor, error) {
	val, _ := c.StringWithSection("security_manager", "adatpor")
	switch val {
	case "local":
		return NewLocalAdaptor(c)
	case "mongo":
	default:
	}
	return NewLocalAdaptor(c)
}
