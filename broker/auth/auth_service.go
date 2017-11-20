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

package auth

import (
	"os"
	"sync"

	"github.com/cloustone/sentel/core"
)

const (
	AuthServiceName = "auth"
)

// AuthServiceFactory
type AuthServiceFactory struct{}

// New create coap service factory
func (p *AuthServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	return &AuthService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
	}, nil
}

// Authentication Service
type AuthService struct {
	ServiceBase core.ServiceBase
}

// Name
func (p *AuthService) Name() string {
	return AuthServiceName
}

// Start
func (p *AuthService) Start() error {
	return nil
}

// Stop
func (p *AuthService) Stop() {
}

// CheckAcl check client's access control right
func (p *AuthService) CheckAcl(clientid string, username string, topic string, access string) error {
	return nil
}

// CheckUserCrenditial check user's name and password
func (p *AuthService) CheckUserCrenditial(username string, password string) error {
	return nil
}

// GetPskKey return user's psk key
func (p *AuthService) GetPskKey(hint string, identity string) (string, error) {
	return "", nil
}
