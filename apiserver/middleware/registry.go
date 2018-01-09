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

package middleware

import (
	"errors"
	"net/http"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/labstack/echo"
)

func RegistryWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			r, err := registry.New("apiserver", config)
			if err != nil {
				return ctx.JSON(http.StatusInternalServerError, errors.New("registry access failed"))
			}
			ctx.Set("registry", r)
			ctx.Set("config", config)
			return next(ctx)
		}
	}
}
