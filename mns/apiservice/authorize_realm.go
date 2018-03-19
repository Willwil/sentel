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
package apiservice

import (
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/goshiro/web"
	"github.com/cloustone/sentel/pkg/registry"
)

type resourceAC struct {
	Path     string
	Resource string
}

type authorizeRealm struct {
	config config.Config
}

func newAuthorizeRealm(c config.Config) (shiro.Realm, error) {
	if _, err := registry.New(c); err != nil {
		return nil, err
	}
	return &authorizeRealm{config: c}, nil
}
func (p *authorizeRealm) GetName() string { return "iot_api_server" }
func (p *authorizeRealm) Supports(token shiro.AuthenticationToken) bool {
	if _, ok := token.(*web.JWTToken); ok {
		return true
	}
	if _, ok := token.(*web.RequestToken); ok {
		return ok
	}
	return false
}

func (p *authorizeRealm) GetPrincipals(token shiro.AuthenticationToken) shiro.PrincipalCollection {
	principals := shiro.NewPrincipalCollection()
	principalName := token.GetPrincipal().(string)
	if r, err := registry.New(p.config); err == nil {
		defer r.Close()
		if tenant, err := r.GetTenant(principalName); err == nil {
			// construct new principals
			principal := shiro.NewPrincipalWithRealm(principalName, p.GetName())
			principal.SetCrenditals(tenant.Password)
			principal.AddRoles(tenant.Roles)
			principals.Add(principal, p.GetName())
		}
	}
	return principals
}