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
	"errors"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/goshiro/web"
	"github.com/golang/glog"
)

func NewSecurityManager(c config.Config, realm ...shiro.Realm) shiro.SecurityManager {
	var securityMgr shiro.SecurityManager
	if mgrName, err := c.StringWithSection("security_manager", "manager"); err != nil {
		switch mgrName {
		case "web":
			if mgr, err := web.NewSecurityManager(c, realm...); err == nil {
				securityMgr = mgr
			}
		}
	}
	if securityMgr == nil {
		securityMgr = &shiro.DefaultSecurityManager{}
	}
	adaptor, err := newAdaptor(c)
	if err != nil {
		glog.Fatal(err)
	}
	securityMgr.SetAdaptor(adaptor)
	return securityMgr
}

func newAdaptor(c config.Config) (shiro.Adaptor, error) {
	return nil, errors.New("not implemented")
}
