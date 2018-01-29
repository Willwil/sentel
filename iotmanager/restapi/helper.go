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
package restapi

import (
	"net/http"

	"github.com/cloustone/sentel/iotmanager/mgrdb"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/labstack/echo"
)

const (
	OK           = http.StatusOK
	ServerError  = http.StatusInternalServerError
	BadRequest   = http.StatusBadRequest
	NotFound     = http.StatusNotFound
	Unauthorized = http.StatusUnauthorized
)

func getConfig(ctx echo.Context) config.Config {
	return ctx.(*apiContext).config
}

func openManagerDB(ctx echo.Context) (mgrdb.ManagerDB, error) {
	c := getConfig(ctx)
	return mgrdb.New(c)
}
