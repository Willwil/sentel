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

import "github.com/cloustone/sentel/broker/base"

// GetVersion return authentication service's version
func GetVersion() string {
	return AuthServiceVersion
}

// CheckAcl check client's access control right
func CheckAcl(clientid string, username string, topic string, access string) error {
	auth := base.GetService(AuthServiceName).(*AuthService)
	return auth.CheckAcl(clientid, username, topic, access)
}

// CheckUserCrenditial check user's name and password
func CheckUserCrenditial(username string, password string) error {
	auth := base.GetService(AuthServiceName).(*AuthService)
	return auth.CheckUserCrenditial(username, password)
}

// GetPskKey return user's psk key
func GetPskKey(hint string, identity string) (string, error) {
	auth := base.GetService(AuthServiceName).(*AuthService)
	return auth.GetPskKey(hint, identity)
}
