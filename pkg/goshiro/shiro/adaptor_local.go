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

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
)

// localAdaptor save policy,role and resource objects in local memory in simple context
type localAdaptor struct {
	roles map[string]Role
}

func NewLocalAdaptor(c config.Config) (*localAdaptor, error) {
	return &localAdaptor{
		roles: make(map[string]Role),
	}, nil
}

func (l *localAdaptor) GetName() string { return "local" }
func (l *localAdaptor) AddRole(r Role) {
	if role, found := l.roles[r.Name]; found {
		l.roles[r.Name] = Role{
			Name:          r.Name,
			PermissionIds: append(role.PermissionIds, r.PermissionIds...),
		}
	} else {
		l.roles[r.Name] = r
	}
}

func (l *localAdaptor) RemoveRole(roleName string) {
	if _, found := l.roles[roleName]; found {
		delete(l.roles, roleName)
	}
}

func (l *localAdaptor) GetRole(roleName string) (Role, error) {
	if r, found := l.roles[roleName]; found {
		return r, nil
	}
	return Role{}, errors.New("not implemented")
}

func (l *localAdaptor) AddRolePermissions(roleName string, permissons []string) {
	if r, found := l.roles[roleName]; found {
		l.roles[roleName] = Role{
			Name:          roleName,
			PermissionIds: append(r.PermissionIds, r.PermissionIds...),
		}
	}
}

func (l *localAdaptor) RemoveRolePermissions(roleName string, permissions []string) {
}

func (l *localAdaptor) GetRolePermissions(roleName string) []string {
	permissions := []string{}
	if r, found := l.roles[roleName]; found {
		permissions = append(permissions, r.PermissionIds...)
	}
	return permissions
}

// RemovePrincipal remove principal from security manager
func (l *localAdaptor) RemovePrincipal(Principal) {}

// GetPriincipalsPermissions return principals permissions
func (l *localAdaptor) GetPrincipalsPermissions(principals PrincipalCollection) []Permission {
	return []Permission{}
}

// AddPrincipalPermissions add permissions for principals
func (l *localAdaptor) AddPrincipalPermissions(principal Principal, permissions []Permission) {
}

// RemovePrincipalsPermissions remove principals's permissions
func (l *localAdaptor) RemovePrincipalPermissions(principal Principal, permissions []Permission) {
}
