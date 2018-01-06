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

	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/config"
)

func Test_CreateAccount(t *testing.T) {
	config := config.New()
	Initialize(config, "mock")
	if err := CreateAccount("hello"); err != nil {
		t.Errorf("failed to create account:%s", err.Error())
	}
	if account, err := GetAccount("hello"); err != nil || account.GetAccountName() != "hello" {
		t.Error("account test failed")
	}
}

func Test_DestoryAccount(t *testing.T) {
	config := config.New()
	Initialize(config, "mock")
	if err := CreateAccount("hello"); err != nil {
		t.Errorf("failed to create account:%s", err.Error())
	}
	if err := ram.DestroyAccount("hello"); err != nil {
		t.Error("account destroy test failed")
	}

}
