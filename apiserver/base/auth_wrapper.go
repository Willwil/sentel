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
	"fmt"

	"github.com/cloustone/sentel/keystone/auth"
	"github.com/cloustone/sentel/keystone/rac"
	"github.com/labstack/echo"
)

const (
	ResourceCategoryTenant  = "tenant"
	ResourceCategoryProduct = "product"
	ResourceCategoryDevice  = "device"
	ResourceCategoryRule    = "rule"
)

func CreateResource(ctx echo.Context, category string, resource string, ower string) {
	req := auth.ApiAuthParam{}
	if err := ctx.Bind(&req); err == nil {
		rac.CreateResource(category, resource, ower, req.AccessKey)
	}
}
func AccessResource(ctx echo.Context, category string, resource string, right rac.AccessRight) error {
	req := auth.ApiAuthParam{}
	if err := ctx.Bind(&req); err == nil {
		return rac.AccessResource(category, resource, req.AccessKey, right)
	}
	return fmt.Errorf("in valid access to resource '%s'", resource)
}
func ReleaseResource(ctx echo.Context, category string, resource string) {
	req := auth.ApiAuthParam{}
	if err := ctx.Bind(&req); err == nil {
		rac.ReleaseResource(category, resource, req.AccessKey)
	}
}
