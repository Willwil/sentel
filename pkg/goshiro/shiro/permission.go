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

import "strings"

type Permission struct {
	Resource string `json:"resource" bson:"Resource"`
	Actions  string `json:"actions" bson:"Actions"`
}

// NewPermission create permission object from string presentation, such as
// "resource1:read,write,update"
func NewPermission(permission string) Permission {
	result := Permission{}
	items := strings.Split(permission, ":")
	if len(items) > 0 {
		result.Resource = items[0]
		if len(items) > 1 {
			result.Actions = items[1]
		}
	}
	return result
}

func (p Permission) Implies(action string, resource string) bool {
	actions := strings.Split(p.Actions, ",")
	if resource == p.Resource {
		for _, val := range actions {
			if val == action {
				return true
			}
		}
	}
	return false
}
