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

func newConfig() config.Config {
	c := config.New()
	c.AddConfig(map[string]map[string]string{
		"keystone": {
			"hosts":           "localhost:4147",
			"mongo":           "localhost:27017",
			"connect_timeout": "2",
		},
	})
	return c
}

func Test_CreateAccount(t *testing.T) {
	c := newConfig()
	if err := Initialize(c, "direct"); err != nil {
		t.Errorf("ram initialization failed:%s", err.Error())
		return
	}
	if err := CreateAccount("account1"); err != nil {
		t.Errorf("CreateAccount failed:%s", err.Error())
	}
	if err := DestroyAccount("account1"); err != nil {
		t.Errorf("CreateAccount failed:%s", err.Error())
	}

}

func Test_CreateResource(t *testing.T) {
	c := newConfig()
	if err := Initialize(c, "direct"); err != nil {
		t.Errorf("ram initialization failed:%s", err.Error())
		return
	}
	rid := NewObjectId()
	if _, err := CreateResource("account1", &ResourceCreateOption{
		ObjectId:   rid,
		Name:       "product1",
		Attributes: []string{"devices", "rules"},
	}); err != nil {
		t.Errorf("CreateResource failed:%s", err.Error())
	}
}
