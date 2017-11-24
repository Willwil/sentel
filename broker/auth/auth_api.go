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
	"errors"

	"github.com/cloustone/sentel/broker/broker"
)

const (
	AuthServiceVersion = "0.1"
	AclRead            = 1
	AclWrite           = 2
)

var (
	ErrorAuthDenied = errors.New("authentication denied")
)

type Options struct {
	ClientId     string
	DeviceName   string
	ProductKey   string
	SecurityMode int
	SignMethod   string
	Timestamp    string
	Password     string
	DeviceSecret string
}

// Keyword declarations used to parse client's request
const (
	SecurityMode = "securemode"
	SignMethod   = "signmethod"
	Timestamp    = "timestamp"
)

// GetVersion return authentication service's version
func GetVersion() string {
	return AuthServiceVersion
}

func Authorize(clientId string, userName string, topic string, access int, opt *Options) error {
	auth := broker.GetService(ServiceName).(*AuthService)
	return auth.authorize(clientId, userName, topic, access, opt)
}

func Authenticate(opt *Options) error {
	auth := broker.GetService(ServiceName).(*AuthService)
	return auth.authenticate(opt)
}
