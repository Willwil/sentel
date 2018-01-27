//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package base

import (
	"github.com/cloustone/sentel/keystone/client"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/config"
)

var (
	noAuth = true
)

func Initialize(c config.Config) error {
	if auth, err := c.String("auth"); err != nil || auth == "none" {
		noAuth = true
		return nil
	}
	return client.Initialize(c)
}

func Authenticate(opts interface{}) error {
	if !noAuth {
		return client.Authenticate(opts)
	}
	return nil
}

func Authorize(accessId string, resource string, action string) error {
	if !noAuth {
		return client.Authorize(accessId, resource, action)
	}
	return nil
}

func CreateResource(accessId string, res ram.ResourceCreateOption) error {
	if !noAuth {
		return client.CreateResource(accessId, res)
	}
	return nil
}

func AccessResource(res string, accessId string, action string) error {
	if !noAuth {
		return client.AccessResource(res, accessId, action)
	}
	return nil
}

func DestroyResource(rid string, accessId string) error {
	if !noAuth {
		return client.DestroyResource(rid, accessId)
	}
	return nil
}

func AddResourceGrantee(res string, accessId string, right string) error {
	if !noAuth {
		return client.AddResourceGrantee(res, accessId, right)
	}
	return nil
}

func CreateAccount(account string) error {
	if !noAuth {
		return client.CreateAccount(account)
	}
	return nil
}

func DestroyAccount(account string) error {
	if !noAuth {
		return client.DestroyAccount(account)
	}
	return nil
}
