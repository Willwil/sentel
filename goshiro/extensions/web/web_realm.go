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
)

type ApiDeclaration struct {
	Url          string `json:"url" bson:"Url"`
	Method       string `json:"method" bson:"Method"`
	Resource     string `json:"resource" bson:"Declaration"`
	AllowedRoles string `json:"allowedRoles" bson:"AllowedRoles"`
}

func (r *ApiDeclaration) ResolveDeclaration(principals shiro.PrincipalCollection, ctx WebRequestContext) error {
	return nil
}

type WebAuthorizeRealm struct {
	realmName    string
	declarations []ApiDeclaration
}

func NewWebAuthorizeRealm(env shiro.Environment, realmName string) *WebAuthorizeRealm {
	return &WebAuthorizeRealm{
		realmName:    realmName,
		declarations: []ApiDeclaration{},
	}
}
func (r *WebAuthorizeRealm) LoadDeclarations(resources []ApiDeclaration) error {
	return nil
}
func (r *WebAuthorizeRealm) LoadDeclarationsWithJsonFile(fname string) error {
	return nil
}
func (r *WebAuthorizeRealm) LoadDeclarationsWithYamlFile(fname string) error {
	return nil
}
func (r *WebAuthorizeRealm) GetDeclarationWithUri(uri string) *ApiDeclaration {
	return nil
}
func (r *WebAuthorizeRealm) GetDeclarationName(uri string, ctx WebRequestContext) (string, error) {
	return "", nil
}

func (r *WebAuthorizeRealm) GetName() string {
	return r.realmName
}
func (r *WebAuthorizeRealm) Supports(token shiro.AuthenticationToken) bool {
	return true
}
func (r *WebAuthorizeRealm) GetAuthenticationInfo(token shiro.AuthenticationToken) (shiro.AuthenticationInfo, error) {
	return nil, errors.New("no implmented")
}
func (r *WebAuthorizeRealm) GetAuthorizationInfo(principals shiro.PrincipalCollection) (shiro.AuthorizationInfo, error) {
	return nil, errors.New("no implmented")
}
func (r *WebAuthorizeRealm) SetCacheEnable(bool) {}
