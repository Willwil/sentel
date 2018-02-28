//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package web

import (
	"errors"

	"github.com/cloustone/sentel/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/config"
)

type ResourceRealm struct {
	realmName string
	resources []Resource
	config    config.Config
}

func NewResourceRealm(c config.Config, realmName string) *ResourceRealm {
	return &ResourceRealm{
		config:    c,
		realmName: realmName,
		resources: []Resource{},
	}
}
func (r *ResourceRealm) LoadResources(resources []Resource) error {
	return nil
}
func (r *ResourceRealm) LoadResourceWithJsonFile(fname string) error {
	return nil
}
func (r *ResourceRealm) LoadResourceWithYamlFile(fname string) error {
	return nil
}
func (r *ResourceRealm) GetResourceWithUri(uri string) *Resource {
	return nil
}
func (r *ResourceRealm) GetResourceName(uri string, ctx ResourceContext) (string, error) {
	return "", nil
}

func (r *ResourceRealm) GetName() string {
	return r.realmName
}
func (r *ResourceRealm) Supports(token shiro.AuthenticationToken) bool {
	return true
}
func (r *ResourceRealm) GetAuthenticationInfo(token shiro.AuthenticationToken) (shiro.AuthenticationInfo, error) {
	return nil, errors.New("no implmented")
}
func (r *ResourceRealm) GetAuthorizationInfo(principals shiro.PrincipalCollection) (shiro.AuthorizationInfo, error) {
	return nil, errors.New("no implmented")
}
func (r *ResourceRealm) SetCacheEnable(bool) {}
