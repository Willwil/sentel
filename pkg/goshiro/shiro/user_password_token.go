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

type UserAndPasswordToken struct {
	Username      string
	Password      string
	Authenticated bool
}

func (p UserAndPasswordToken) GetPrincipal() interface{} {
	return p.Username
}
func (p UserAndPasswordToken) GetCrenditals() interface{} {
	return p.Password
}

func (p UserAndPasswordToken) IsAuthenticated() bool {
	return p.Authenticated
}
