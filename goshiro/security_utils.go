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

package goshiro

import (
	"github.com/cloustone/sentel/goshiro/auth"
	"github.com/cloustone/sentel/pkg/config"
)

var (
	securityManager auth.SecurityManager
)

func InitializeSecurityManager(c config.Config) error {
	mgr, err := auth.NewSecurityManager(c)
	if err != nil {
		return err
	}
	securityManager = mgr
	return nil
}

func GetSubject(token auth.AuthenticationToken) (auth.Subject, error) {
	return securityManager.GetSubject(token)
}

func NewSubject() (auth.Subject, error) {
	return securityManager.NewSubject()
}

func GetResourceName(uri string, ctx auth.ResourceContext) (string, error) {
	return securityManager.GetResourceName(uri, ctx)
}

func GetSecurityManager() auth.SecurityManager {
	return securityManager
}
