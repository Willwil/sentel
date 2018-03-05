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

package goshiro

import (
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/adaptors"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/golang/glog"
)

func NewSecurityManager(c config.Config, policies []shiro.AuthorizePolicy, realm ...shiro.Realm) shiro.SecurityManager {
	securityMgr, _ := shiro.NewSecurityManager(c)
	adaptor, err := NewAdaptor(c)
	if err != nil {
		glog.Fatal(err)
	}
	securityMgr.SetAdaptor(adaptor)
	securityMgr.LoadPolicies(policies)
	for _, r := range realm {
		securityMgr.AddRealm(r)
	}
	return securityMgr
}

func NewAdaptor(c config.Config) (shiro.Adaptor, error) {
	val, err := c.StringWithSection("security_manager", "adatpor")
	if err != nil {
		val = "simple"
	}
	switch val {
	case "simple":
		return adaptors.NewSimpleAdaptor(c)
	default:
		return adaptors.NewSimpleAdaptor(c)
	}
}
