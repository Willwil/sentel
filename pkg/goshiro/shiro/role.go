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

// Role contains permissions for specified resource
type Role struct {
	Name        string       `json:"name" bson:"Name"`
	Permissions []Permission `json:"permissions" bson:"Permissions"`
}

func (r Role) Implies(action string, resource string) bool {
	for _, permission := range r.Permissions {
		if permission.Implies(action, resource) {
			return true
		}
	}
	return false
}
