//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
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

type WebRequestToken struct{}

func NewWebRequestToken(ctx WebRequestContext) *WebRequestToken {
	return &WebRequestToken{}
}

func (p WebRequestToken) GetPrincipal() interface{} {
	return nil
}
func (p WebRequestToken) GetCrenditals() interface{} {
	return nil
}

func (p WebRequestToken) GetPermission() string {
	return ""
}
