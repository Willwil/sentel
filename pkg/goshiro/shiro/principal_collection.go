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

type PrincipalCollection interface {
	Add(principal Principal, realmName string)
	AddAll(principals PrincipalCollection)
	GetPrimaryPrincipal() Principal
	GetRealmPrincipals(realmName string) []Principal
	Remove(principalName string, realmName string)
	RemoveRealmPrincipals(realmName string)
	GetRealmNames() []string
	Clear()
	IsEmpty() bool
}

func NewPrincipalCollection() PrincipalCollection {
	return &principalCollection{
		principals:    make(map[string][]Principal),
		principalList: []Principal{},
		index:         0,
	}
}

type principalCollection struct {
	principals    map[string][]Principal
	principalList []Principal
	index         int
}

func (p *principalCollection) Add(principal Principal, realmName string) {
	if _, found := p.principals[realmName]; !found {
		p.principals[realmName] = []Principal{}
	}
	p.principals[realmName] = append(p.principals[realmName], principal)
	p.principalList = append(p.principalList, principal)
}

func (p *principalCollection) AddAll(r PrincipalCollection) {
	realmNames := r.GetRealmNames()
	for _, realmName := range realmNames {
		principals := r.GetRealmPrincipals(realmName)
		for _, principal := range principals {
			p.Add(principal, realmName)
		}
	}
}

func (p *principalCollection) Remove(principalName string, realmName string) {
	principals := p.GetRealmPrincipals(realmName)
	for index, principal := range principals {
		if principal.Name() == principalName {
			// remove from principal maps
			principals = append(principals[:index], principals[index:]...)
			// remove from principal list
			for index, principal := range p.principalList {
				if principal.Name() == principalName {
					p.principalList = append(p.principalList[:index], p.principalList[index:]...)
				}
			}
		}
	}
}

func (p *principalCollection) RemoveRealmPrincipals(realmName string) {
	if principals, found := p.principals[realmName]; found {
		delete(p.principals, realmName)
		for _, principal := range principals {
			for index, pp := range p.principalList {
				if principal == pp {
					p.principalList = append(p.principalList[:index], p.principalList[index:]...)
				}
			}
		}
	}
}

func (p *principalCollection) Clear() {
	for realm, _ := range p.principals {
		delete(p.principals, realm)
	}
	p.principalList = []Principal{}
	p.index = 0
}

func (p *principalCollection) GetPrimaryPrincipal() Principal {
	if p.index < len(p.principalList) {
		principal := p.principalList[p.index]
		p.index += 1
		return principal
	}
	return nil
}

func (p *principalCollection) GetRealmPrincipals(realmName string) []Principal {
	realms := []Principal{}
	if r, found := p.principals[realmName]; found {
		realms = append(realms, r...)
	}
	return realms
}

func (p *principalCollection) GetRealmNames() []string {
	realmNames := []string{}
	for realmName, _ := range p.principals {
		realmNames = append(realmNames, realmName)
	}
	return realmNames
}
func (p *principalCollection) IsEmpty() bool {
	return len(p.principals) == 0
}
