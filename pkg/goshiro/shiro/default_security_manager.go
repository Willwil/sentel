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

import "errors"

type DefaultSecurityManager struct {
}

func (d *DefaultSecurityManager) AddRealms([]Realm)    {}
func (d *DefaultSecurityManager) SetAdaptor(a Adaptor) {}
func (d *DefaultSecurityManager) Login(AuthenticationToken) (Subject, error) {
	return nil, errors.New("not implmented")
}
func (d *DefaultSecurityManager) Logout(subject Subject) error {
	return nil
}
func (d *DefaultSecurityManager) GetSubject(token AuthenticationToken) (Subject, error) {
	return nil, errors.New("not implemented")
}
func (d *DefaultSecurityManager) Authorize(subject Subject, req Request) error {
	return errors.New("not implemented")
}
