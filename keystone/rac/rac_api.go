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
package rac

type AccessRight uint16

const (
	AccessRightRead  AccessRight = 1
	AccessRightWrite AccessRight = 2
	AccessRightFull  AccessRight = 3
)

func CreateResource(category string, name string, ower string, key string) {}
func AccessResource(category string, name string, key string, right AccessRight) error {
	return nil
}
func ReleaseResource(category string, name string, key string) {}
