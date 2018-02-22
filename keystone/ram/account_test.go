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

func TestAccount_initialize(t *testing.T) {
	c := config.New("keystone")
	c.AddConfig(config.M{
		"keystone": {
			"hosts":           "localhost:4147",
			"mongo":           "localhost:27017",
			"connect_timeout": 2,
		},
	})
	if err := Initialize(c, "direct"); err != nil {
		t.Error(err)
	}
}

func TestRam_CreateAccount(t *testing.T) {
	if err := CreateAccount("account1"); err != nil {
		t.Error(err)
	}
}

func TestRam_DestroyAccount(t *testing.T) {
	if err := DestroyAccount("account1"); err != nil {
		t.Error(err)
	}
}
