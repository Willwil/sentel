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

const (
	ActionRead   = "read"
	ActionWrite  = "write"
	ActionCreate = "create"
	ActionRemove = "remove"
)

type Permission struct {
	Resources []string `json:"resources" bson:"Resources"`
	Actions   string   `json:"actions" bson:"Actions"`
}

// NewPermission create permission object from string presentation, such as
// "resource1:read,write,update"
func NewPermission(resources []string, action string) Permission {
	return Permission{Resources: resources, Actions: action}
}

func (p Permission) Implies(resource string, action string) bool {
	actions := strings.Split(p.Actions, ",")
	for _, res := range p.Resources {
		if resource == res {
			for _, val := range actions {
				if val == action {
					return true
				}
			}
		}
	}
	return false
}
