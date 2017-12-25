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
package orm

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type AccessRight uint8

const (
	AccessRightReadOnly AccessRight = 1
	AccessRightWrite    AccessRight = 2
	AccessRightFull     AccessRight = 3
)

type Object struct {
	ObjectName     string
	ObjectId       string
	ParentObjectId string
	Category       string
	CreatedTime    time.Time
	Owner          string
}

func NewObjectId() string {
	return uuid.NewV4().String()
}
func CreateObject(obj *Object) {

}

func AccessObject(objName string, accessId string, right AccessRight) error {
	return nil
}
func ReleaseObject(objName string, accessId string) {}
func AssignObjectRight(objName string, accessId string, right AccessRight) error {
	return nil
}
