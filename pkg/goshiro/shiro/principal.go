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

type Principal interface {
	Name() string
	GetRealmName() string
	GetCrenditals() string
}

type simplePrincipal struct {
	principalName string
	crenditals    string
	realmName     string
}

func NewPrincipal(crendital string, realmName string) Principal {
	return &simplePrincipal{}
}

func (p *simplePrincipal) Name() string          { return p.principalName }
func (p *simplePrincipal) GetRealmName() string  { return p.realmName }
func (p *simplePrincipal) GetCrenditals() string { return p.crenditals }
