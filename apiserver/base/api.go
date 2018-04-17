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
	"net/http"
	"strings"

	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
)

type ApiContext struct {
	echo.Context
	//Config      config.Config
	SecurityMgr shiro.SecurityManager
}

type ApiJWTClaims struct {
	jwt.StandardClaims
	AccessId string `json:"accessId"`
}

func GetSecurityManager(ctx echo.Context) shiro.SecurityManager {
	return ctx.Get("SecurityManager").(shiro.SecurityManager)
}

func GetRequestInfo(ctx echo.Context, resourceMaps map[string]string) (string, string) {
	action := ""
	switch ctx.Request().Method {
	case http.MethodGet:
		action = shiro.ActionRead
	default:
		action = shiro.ActionWrite
	}
	if resourceName, found := resourceMaps[ctx.Path()]; found {
		// generate the requested real resource, using parameter from context to replace identifier in resource field
		resource := "/"
		items := strings.Split(resourceName, "/")
		for _, item := range items {
			if item == ""{continue}
			switch item[0] {
			case '$':
				val := ctx.Get(item[1:]).(string)
				resource += val
				resource += "/"
			case ':':
				val := ctx.Param(item[1:])
				resource += val
				resource += "/"
			default:
				resource += item
			}
			strings.TrimRight(resource, "/")
		}

		return resource, action
	}
	return "", action
}
