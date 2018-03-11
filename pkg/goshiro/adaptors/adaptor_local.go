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

package adaptors

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/golang/glog"
)

// LocalAdaptor save policy,role and resource objects in local memory in simple context
type LocalAdaptor struct {
	roles     map[string]shiro.Role
	resources map[string]shiro.Resource
}

func NewLocalAdaptor(c config.Config) (*LocalAdaptor, error) {
	return &LocalAdaptor{
		roles:     make(map[string]shiro.Role),
		resources: make(map[string]shiro.Resource),
	}, nil
}

func (l *LocalAdaptor) GetName() string { return "local" }
func (l *LocalAdaptor) AddRole(r shiro.Role) {
	if role, found := l.roles[r.Name]; found {
		l.roles[r.Name] = shiro.Role{
			Name:        r.Name,
			Permissions: append(role.Permissions, r.Permissions...),
		}
	} else {
		l.roles[r.Name] = r
	}
}

func (l *LocalAdaptor) RemoveRole(roleName string) {
	if _, found := l.roles[roleName]; found {
		delete(l.roles, roleName)
	}
}

func (l *LocalAdaptor) GetRole(roleName string) (shiro.Role, error) {
	if r, found := l.roles[roleName]; found {
		return r, nil
	}
	return shiro.Role{}, errors.New("not implemented")
}

func (l *LocalAdaptor) AddRolePermissions(roleName string, permissons []shiro.Permission) {
	if r, found := l.roles[roleName]; found {
		l.roles[roleName] = shiro.Role{
			Name:        roleName,
			Permissions: append(r.Permissions, r.Permissions...),
		}
	}
}

func (l *LocalAdaptor) RemoveRolePermissions(roleName string, permissions []shiro.Permission) {
}

func (l *LocalAdaptor) GetRolePermissions(roleName string) []shiro.Permission {
	permissions := []shiro.Permission{}
	if r, found := l.roles[roleName]; found {
		permissions = append(permissions, r.Permissions...)
	}
	return permissions
}

func (l *LocalAdaptor) AddResource(r shiro.Resource) {
	if _, found := l.resources[r.Name]; found {
		glog.Infof("add same resource '%s' twice", r.Name)
	}
	l.resources[r.Name] = r
}

func (l *LocalAdaptor) RemoveResource(r shiro.Resource) {
	if _, found := l.resources[r.Name]; found {
		delete(l.resources, r.Name)
	}
}

func (l *LocalAdaptor) GetResource(resourceName string) *shiro.Resource {
	if r, found := l.resources[resourceName]; found {
		return &r
	}
	return nil
}

func (l *LocalAdaptor) GetResources() []shiro.Resource {
	resources := []shiro.Resource{}
	for _, r := range l.resources {
		resources = append(resources, r)
	}
	return resources
}
