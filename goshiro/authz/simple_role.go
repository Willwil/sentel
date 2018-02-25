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

package authz

type Role interface {
	GetName() string
	SetName(name string)
	GetPermissions() []Permission
	SetPermissions(permissions []Permission)
	AddPermission(permission Permission)
	AddPermissions(permission []Permission)
	IsPermitted(permission Permission)
	IsEqual(obj interface{}) bool
}

func NewRole() Role {
	return nil
}

func NewRoleWithName(name string) Role {
	return nil
}

type simpleRole struct {
	name        string
	permissions []string
}
