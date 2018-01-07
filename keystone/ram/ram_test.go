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
package ram

import (
	"testing"

	"github.com/cloustone/sentel/pkg/config"
)

var (
	resourceId = NewObjectId()
)

func Test_Initialize(t *testing.T) {
	c := config.New()
	c.AddConfig(map[string]map[string]string{
		"keystone": {
			"hosts":           "localhost:4147",
			"mongo":           "localhost:27017",
			"connect_timeout": "2",
		},
	})
	if err := Initialize(c, "direct"); err != nil {
		t.Error(err)
		return
	}
}

func Test_CreateAccount(t *testing.T) {
	if err := CreateAccount("account1"); err != nil {
		t.Error(err)
	}
}

func Test_CreateResource(t *testing.T) {
	if _, err := CreateResource("account1", &ResourceCreateOption{
		ObjectId:   resourceId,
		Name:       "product1",
		Attributes: []string{"devices", "rules"},
	}); err != nil {
		t.Error(err)
	}
	if _, err := CreateResource("account1", &ResourceCreateOption{
		ObjectId:   resourceId,
		Name:       "product1",
		Attributes: []string{"devices", "rules"},
	}); err == nil {
		t.Error(err)
	}
}

func Test_AddResourceGrantee(t *testing.T) {
	if err := AddResourceGrantee(resourceId, "client1", RightRead); err != nil {
		t.Error(err)
	}
	if err := AddResourceGrantee(resourceId+"/devices", "client1", RightWrite); err != nil {
		t.Error(err)
	}
}

func Test_Authorize(t *testing.T) {
	if err := Authorize(resourceId, "client1", ActionRead); err != nil {
		t.Error(err)
	}
	if err := Authorize(resourceId+"/devices", "client1", ActionRead); err != nil {
		t.Error(err)
	}
}

func Test_DestroyResource(t *testing.T) {
	if err := DestroyResource(resourceId); err != nil {
		t.Error(err.Error())
	}
}

func Test_DestroyAccount(t *testing.T) {
	if err := DestroyAccount("account1"); err != nil {
		t.Error(err)
	}
}
