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
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/goshiro/web"
	"github.com/golang/glog"
)

func NewSecurityManager(c config.Config, realm ...shiro.Realm) shiro.SecurityManager {
	mgrName, err := c.StringWithSection("security_manager", "manager")
	if err != nil {
		glog.Error("invalid security manager configuration, using default...")
		return &shiro.DefaultSecurityManager{}
	}
	switch mgrName {
	case "web":
		if mgr, err := web.NewSecurityManager(c, realm...); err == nil {
			return mgr
		} else {
			glog.Error("web security manager creattion failed, using default")
		}
	}
	return &shiro.DefaultSecurityManager{}
}
