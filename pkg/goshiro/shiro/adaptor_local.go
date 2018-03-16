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
	"fmt"

	"github.com/cloustone/sentel/pkg/config"
)

// localAdaptor save policy,role and resource objects in local memory in simple context
type localAdaptor struct {
	roles                map[string]Role       // all roles
	principals           map[string]Principal  // all principals
	permissions          map[string]Permission // all permissions
	principalRoles       map[string][]Role     // what roles the principal have
	principalPermissions map[string][]string   // what permission the principal have
	permIndex            uint64
}

func NewLocalAdaptor(c config.Config) (*localAdaptor, error) {
	return &localAdaptor{
		roles:                make(map[string]Role),
		principals:           make(map[string]Principal),
		permissions:          make(map[string]Permission),
		principalRoles:       make(map[string][]Role),
		principalPermissions: make(map[string][]string),
		permIndex:            1,
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
	if role, found := l.roles[roleName]; found {
		for _, permission := range permissions {
			for index, permId := range role.PermissionIds {
				if permId == permission {
					role.PermissionIds = append(role.PermissionIds[:index], role.PermissionIds[index:]...)
				}
			}
		}
	}
}

func (l *localAdaptor) GetRolePermissions(roleName string) []string {
	permissions := []string{}
	if r, found := l.roles[roleName]; found {
		permissions = append(permissions, r.PermissionIds...)
	}
	return permissions
}

// RemovePrincipal remove principal from
func (l *localAdaptor) RemovePrincipal(principal Principal) {
	delete(l.principals, principal.Name())
	delete(l.principalRoles, principal.Name())
	delete(l.principalPermissions, principal.Name())
}

// GetPriincipalsPermissions return principals permissions
func (l *localAdaptor) GetPrincipalsPermissions(principals PrincipalCollection) []Permission {
	principal := principals.GetPrimaryPrincipal()
	return l.GetPrincipalPermissions(principal)
}

// GetPrincipalPermissions return principals permissions
func (l *localAdaptor) GetPrincipalPermissions(principal Principal) []Permission {
	if permIds, found := l.principalPermissions[principal.Name()]; found {
		return l.GetPermissions(permIds)
	}
	return []Permission{}
}

func (l *localAdaptor) SetPrincipalPermissions(principal Principal, permissions []Permission) {
	permIds := l.AddPermissions(permissions)
	if len(permIds) > 0 {
		l.principalPermissions[principal.Name()] = permIds
	}
}

// AddPrincipalPermissions add permissions for principals
func (l *localAdaptor) AddPrincipalPermissions(principal Principal, permissions []Permission) {
	name := principal.Name()
	if _, found := l.principalPermissions[name]; !found {
		l.principalPermissions[name] = []string{}
	}
	if permIds := l.AddPermissions(permissions); len(permIds) > 0 {
		l.principalPermissions[name] = append(l.principalPermissions[name], permIds...)
	}
}

// RemovePrincipalsPermissions remove principals's permissions
func (l *localAdaptor) RemovePrincipalPermissions(principal Principal, permissions []Permission) {
	if permIds, found := l.principalPermissions[principal.Name()]; found {
		// get all resources that removed permissions will be applied
		resourceList := []string{}
		for _, permission := range permissions {
			resourceList = append(resourceList, permission.Resources...)
		}
		// remove duplicated resource
		resources := make(map[string]bool)
		for _, resource := range resourceList {
			resources[resource] = true
		}

		// for each resource, remove from original permissions
		permissions = l.GetPermissions(permIds)
		for resource, _ := range resources {
			for _, permission := range permissions {
				for index, name := range permission.Resources {
					if resource == name {
						permission.Resources = append(permission.Resources[:index], permission.Resources[index:]...)
					}
				}
			}
		}
		l.SetPrincipalPermissions(principal, permissions)
	}
}

func (l *localAdaptor) SetPrincipalRole(principal Principal, roleNames ...string) error {
	roles := []Role{}
	for _, roleName := range roleNames {
		if role, err := l.GetRole(roleName); err != nil {
			return err
		} else {
			roles = append(roles, role)
		}
	}
	principalName := principal.Name()
	if _, found := l.principalRoles[principalName]; !found {
		l.principalRoles[principalName] = []Role{}
	}
	l.principalRoles[principalName] = append(l.principalRoles[principalName], roles...)
	return nil
}

func (l *localAdaptor) GetPrincipalRoles(principal Principal) []Role {
	if roles, found := l.principalRoles[principal.Name()]; found {
		return roles
	}
	return []Role{}
}

func (l *localAdaptor) GetPrincipalRoleNames(principal Principal) []string {
	roleNames := []string{}
	if roles, found := l.principalRoles[principal.Name()]; found {
		for _, role := range roles {
			roleNames = append(roleNames, role.Name)
		}
	}
	return roleNames
}

func (l *localAdaptor) RemovePrincipalRoles(principal Principal) {
	delete(l.principalRoles, principal.Name())
}

func (l *localAdaptor) AddPermissions(permissions []Permission) []string {
	permIds := []string{}
	for _, permission := range permissions {
		if permission.PermissionId == "" {
			permission.PermissionId = fmt.Sprintf("%d", l.permIndex)
			l.permIndex += 1
		}
		l.permissions[permission.PermissionId] = permission
		permIds = append(permIds, permission.PermissionId)
	}
	return permIds
}

func (l *localAdaptor) RemovePermissions(permIds []string) {
	for _, permId := range permIds {
		delete(l.permissions, permId)
	}
}

func (l *localAdaptor) GetPermissions(permIds []string) []Permission {
	permissions := []Permission{}
	for _, permId := range permIds {
		if permission, found := l.permissions[permId]; found {
			permissions = append(permissions, permission)
		}
	}
	return permissions
}
