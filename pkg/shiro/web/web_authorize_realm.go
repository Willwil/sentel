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
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/shiro"
)

type ApiAuthorizePolicy struct {
	Url          string `json:"url" bson:"Url"`
	Method       string `json:"method" bson:"Method"`
	Resource     string `json:"resource" bson:"AuthorizePolicy"`
	AllowedRoles string `json:"allowedRoles" bson:"AllowedRoles"`
}

func (r *ApiAuthorizePolicy) GetName() string        { return r.Url }
func (r *ApiAuthorizePolicy) GetPermissions() string { return r.AllowedRoles }
func (r *ApiAuthorizePolicy) ResolveName(path string, ctx shiro.RequestContext) (string, error) {
	return "", nil
}

type WebAuthorizeRealm struct {
	realmName string
	policies  []*ApiAuthorizePolicy
}

func NewWebAuthorizeRealm(c config.Config, realmName string) *WebAuthorizeRealm {
	return &WebAuthorizeRealm{
		realmName: realmName,
		policies:  []*ApiAuthorizePolicy{},
	}
}
func (w *WebAuthorizeRealm) LoadPolicies(policies []ApiAuthorizePolicy) error {
	return nil
}
func (w *WebAuthorizeRealm) GetName() string {
	return w.realmName
}

func (w *WebAuthorizeRealm) GetPolicy(path string, ctx shiro.RequestContext) shiro.Policy {
	return nil
}

func (w *WebAuthorizeRealm) AddPolicy(policy shiro.Policy)    {}
func (w *WebAuthorizeRealm) RemovePolicy(policy shiro.Policy) {}
func (w *WebAuthorizeRealm) GetAllPolicies() []shiro.Policy   { return []shiro.Policy{} }
func (w *WebAuthorizeRealm) Supports(token shiro.AuthenticationToken) bool {
	return false
}

func (w *WebAuthorizeRealm) GetAuthorizationInfo() shiro.AuthorizationInfo {
	return nil
}
