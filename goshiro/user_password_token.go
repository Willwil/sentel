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

type UserAndPasswordToken struct {
	username string
	password string
}

func NewUserAndPasswordToken(user string, pwd string) *UserAndPasswordToken {
	return &UserAndPasswordToken{username: user, password: pwd}
}

func (p UserAndPasswordToken) GetPrincipal() interface{} {
	return p.username
}
func (p UserAndPasswordToken) GetCrenditals() interface{} {
	return p.password
}
