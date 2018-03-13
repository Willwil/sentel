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

// defaultSecurityManager is default security manager that manage authorization
// and authentication
type defaultSecurityManager struct {
	realms            map[string]Realm // All backend realms
	roleManager       RoleManager
	permissionManager PermissionManager
}

func newDefaultSecurityManager(c config.Config, adaptor Adaptor, realm ...Realm) (SecurityManager, error) {
	realms := make(map[string]Realm)
	for _, r := range realm {
		realms[r.GetName()] = r
	}
	return &defaultSecurityManager{
		realms:            realms,
		roleManager:       NewRoleManager(adaptor),
		permissionManager: NewPermissionManager(adaptor),
	}, nil
}

func (w *defaultSecurityManager) Load() {
}

func (w *defaultSecurityManager) AddRealm(realm ...Realm) {
	for _, r := range realm {
		w.realms[r.GetName()] = r
	}
}

// GetRealm return specified realm by realm name
func (w *defaultSecurityManager) GetRealm(realmName string) Realm {
	if realm, found := w.realms[realmName]; found {
		return realm
	}
	return nil
}

func (w *defaultSecurityManager) Login(token AuthenticationToken) (Principal, error) {
	for _, r := range w.realms {
		if r.Supports(token) {
			principals := r.GetPrincipals(token)
			if !principals.IsEmpty() {
				pp := principals.GetPrimaryPrincipal()
				if pp.GetCrenditals() == token.GetCrenditals() { // valid principal found
					return pp, nil
				}
			}
		}
	}
	return nil, errors.New("invalid token")
}

// getPrincipals return authentication token's principals
func (w *defaultSecurityManager) getPrincipals(token AuthenticationToken) PrincipalCollection {
	principals := NewPrincipalCollection()
	if token.IsAuthenticated() {
		for _, r := range w.realms {
			if r.Supports(token) {
				pps := r.GetPrincipals(token)
				if !pps.IsEmpty() {
					principals.AddAll(pps)
				}
			}
		}
	}
	return principals
}

func (w *defaultSecurityManager) RemovePrincipal(principal Principal) {
	w.roleManager.RemovePrincipalRoles(principal)
	w.permissionManager.RemovePrincipal(principal)
}

func (w *defaultSecurityManager) GetPrincipalsPermissions(principals PrincipalCollection) []Permission {
	return w.permissionManager.GetPrincipalsPermissions(principals)
}

// AddPrincipalPermissions add permissions for principals
func (w *defaultSecurityManager) AddPrincipalPermissions(principal Principal, permissions []Permission) {
	w.permissionManager.AddPrincipalPermissions(principal, permissions)
}

// RemovePrincipalPermissions remove principals's permissions
func (w *defaultSecurityManager) RemovePrincipalPermissions(principal Principal, permissions []Permission) {
	w.permissionManager.RemovePrincipalPermissions(principal, permissions)
}

// Authorize check wether the subject is authorized with the request
// It will iterate all realms to get principals and roles, comparied with saved
// authorization policies
func (w *defaultSecurityManager) Authorize(token AuthenticationToken, resource string, action string) error {
	// if no resource is requested, the authorization is ok
	if resource == "" {
		return nil
	}
	if token != nil {
		if principals := w.getPrincipals(token); !principals.IsEmpty() {
			permissions := w.GetPrincipalsPermissions(principals)
			for _, permission := range permissions {
				if permission.Implies(resource, action) {
					return nil
				}
			}
		}
	}
	return errors.New("not authorized")
}

// AddRole add role with permission into realm
func (p *defaultSecurityManager) AddRole(r Role) {
	p.roleManager.AddRole(r)
}

// RemoveRole remove specified role from realm
func (p *defaultSecurityManager) RemoveRole(roleName string) {
	p.roleManager.RemoveRole(roleName)
}

// GetRole return role's detail information
func (p *defaultSecurityManager) GetRole(roleName string) (Role, error) {
	return p.roleManager.GetRole(roleName)
}

// AddRolePermission add permissions to a role
func (p *defaultSecurityManager) AddRolePermissions(roleName string, permissions []Permission) {
	permids := []string{}
	for _, permission := range permissions {
		permid := p.permissionManager.AddPermission(permission)
		permids = append(permids, permid)
	}
	p.roleManager.AddRolePermissions(roleName, permids)
}

// RemoveRolePermissions remove permission from role
func (p *defaultSecurityManager) RemoveRolePermissions(roleName string, permissions []Permission) {
	//p.roleManager.RemoveRolePermissions(roleName, permissions)
}

// RemoveRoleAllPermissions remove role's all permissions
func (p *defaultSecurityManager) RemoveRoleAllPermissions(roleName string) {
	permissions := p.roleManager.GetRolePermissions(roleName)
	p.roleManager.RemoveRolePermissions(roleName, permissions)
}

// GetRolePermission return specfied role's all permissions
func (p *defaultSecurityManager) GetRolePermissions(roleName string) []Permission {
	permissions := []Permission{}
	permids := p.roleManager.GetRolePermissions(roleName)
	for _, permid := range permids {
		if permission, err := p.permissionManager.GetPermission(permid); err == nil {
			permissions = append(permissions, permission)
		}
	}
	return permissions
}
