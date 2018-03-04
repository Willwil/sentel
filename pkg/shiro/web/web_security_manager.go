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

package web

import (
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/shiro"
)

type WebSecurityManager struct {
}

func NewSecurityManager(c config.Config, realm ...shiro.Realm) (*WebSecurityManager, error) {
	return &WebSecurityManager{}, nil
}

func (w *WebSecurityManager) AddRealms([]shiro.Realm) {}
func (w *WebSecurityManager) Login(shiro.AuthenticationToken) (shiro.Subject, error) {
	return nil, errors.New("not implemented")
}
func (w *WebSecurityManager) Logout(subject shiro.Subject) error {
	return nil
}
func (w *WebSecurityManager) GetSubject(token shiro.AuthenticationToken) (shiro.Subject, error) {
	return nil, errors.New("not implmented")
}
func (w *WebSecurityManager) Save(subject shiro.Subject) {}
func (w *WebSecurityManager) CheckPermission(subject shiro.Subject, resourceName string, action string) error {
	return nil
}

func (w *WebSecurityManager) Authorize(subject shiro.Subject, req shiro.Request) error {
	return nil
}
func NewWebSecurityManager(c config.Config, realm ...shiro.Realm) (shiro.SecurityManager, error) {
	return nil, errors.New("not implemented")
}
