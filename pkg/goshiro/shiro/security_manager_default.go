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
	realms          []Realm // All backend realms
	roleManager     RoleManager
	policyManager   PolicyManager
	resourceManager ResourceManager
}

func NewDefaultSecurityManager(c config.Config, adaptor Adaptor, realm ...Realm) (SecurityManager, error) {
	return &defaultSecurityManager{
		realms:          realm,
		roleManager:     NewRoleManager(adaptor),
		policyManager:   NewPolicyManager(adaptor),
		resourceManager: NewResourceManager(adaptor),
	}, nil
}

func (w *defaultSecurityManager) AddRealm(r ...Realm) {
	w.realms = append(w.realms, r...)
}

// GetRealm return specified realm by realm name
func (w *defaultSecurityManager) GetRealm(realmName string) Realm {
	for _, r := range w.realms {
		if r.GetName() == realmName {
			return r
		}
	}
	return nil
}

func (w *defaultSecurityManager) AddPolicies(policies []AuthorizePolicy) {
	w.policyManager.AddPolicies(policies)
}

func (w *defaultSecurityManager) RemovePolicy(policy AuthorizePolicy) {
	w.policyManager.RemovePolicy(policy)
}

func (w *defaultSecurityManager) GetPolicy(path string) (AuthorizePolicy, error) {
	return w.policyManager.GetPolicy(path)
}

func (w *defaultSecurityManager) GetAllPolicies() []AuthorizePolicy {
	return w.policyManager.GetAllPolicies()
}

func (w *defaultSecurityManager) Login(token AuthenticationToken) (Subject, error) {
	for _, r := range w.realms {
		if r.Supports(token) {
			principals := r.GetPrincipals(token)
			if !principals.IsEmpty() {
				pp := principals.GetPrimaryPrincipal()
				if pp.GetCrenditals() == token.GetCrenditals() { // valid principal found
					return nil, nil
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

func anySliceElementInSlice(slice1, slice2 []string) bool {
	for _, elem1 := range slice1 {
		for _, elem2 := range slice2 {
			if elem1 == elem2 {
				return true
			}
		}
	}
	return false
}

func (w *defaultSecurityManager) getPrincipalsPermissions(principals PrincipalCollection) []Permission {
	return []Permission{}
}

// Authorize check wether the subject is authorized with the request
// It will iterate all realms to get principals and roles, comparied with saved
// authorization policies
func (w *defaultSecurityManager) Authorize(token AuthenticationToken, resource, action string) error {
	if token != nil {
		// get principals and roles
		if principals := w.getPrincipals(token); !principals.IsEmpty() {
			permissions := w.getPrincipalsPermissions(principals)
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
func (p *defaultSecurityManager) AddRole(r Role) { p.roleManager.AddRole(r) }

// RemoveRole remove specified role from realm
func (p *defaultSecurityManager) RemoveRole(roleName string) {
	p.roleManager.RemoveRole(roleName)
}

// GetRole return role's detail information
func (p *defaultSecurityManager) GetRole(roleName string) *Role {
	return p.roleManager.GetRole(roleName)
}

// AddRolePermission add permissions to a role
func (p *defaultSecurityManager) AddRolePermissions(roleName string, permissions []Permission) {
	p.roleManager.AddRolePermissions(roleName, permissions)
}

// RemoveRolePermissions remove permission from role
func (p *defaultSecurityManager) RemoveRolePermissions(roleName string, permissions []Permission) {
	p.roleManager.RemoveRolePermissions(roleName, permissions)
}

// GetRolePermission return specfied role's all permissions
func (p *defaultSecurityManager) GetRolePermissions(roleName string) []Permission {
	return p.roleManager.GetRolePermissions(roleName)
}
