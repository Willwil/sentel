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

import "fmt"

type defaultPolicyManager struct {
	policies []AuthorizePolicy
	adaptor  Adaptor
}

// AddPolicies load authorization policies
func (p *defaultPolicyManager) AddPolicies(policies []AuthorizePolicy) {
	p.policies = append(p.policies, policies...)
	if p.adaptor != nil {
		for _, policy := range policies {
			p.adaptor.AddPolicy(policy)
		}
	}
}

// RemovePolicy remove specified authorization policy
func (p *defaultPolicyManager) RemovePolicy(policy AuthorizePolicy) {
	for index, pp := range p.policies {
		if pp.Equal(policy) {
			p.policies = append(p.policies[:index], p.policies[index:]...)
			if p.adaptor != nil {
				p.adaptor.RemovePolicy(policy)
			}
			return
		}
	}
}

// GetPolicy return specified authorization policy
func (p *defaultPolicyManager) GetPolicy(path string) (AuthorizePolicy, error) {
	for _, policy := range p.policies {
		if policy.Path == path {
			return policy, nil
		}
	}
	return AuthorizePolicy{}, fmt.Errorf("invalid policy path '%s'", path)
}

// GetAllPolicies return all authorization policy
func (p *defaultPolicyManager) GetAllPolicies() []AuthorizePolicy {
	return p.policies
}

// Reload load all policies from database
func (p *defaultPolicyManager) Reload() {
	if p.adaptor != nil {
		p.policies = p.adaptor.GetAllPolicies()
	}
}
