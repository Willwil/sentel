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

package console

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/labstack/echo"
)

var defaultConfigs = config.M{
	"apiserver": {
		"loglevel": "debug",
		"kafka":    "localhost:9092",
		"version":  "v1",
		"auth":     "jwt",
		"mongo":    "localhost:27017",
		"swagger":  "0.0.0.0:53384",
	},
	"registry": {
		"hosts":    "localhost:27017",
		"loglevel": "debug",
	},
	"security": {
		"cafile":              "",
		"capath":              "",
		"certfile":            "",
		"keyfile":             "",
		"require_certificate": false,
	},
}

const accessId = "jenson"

var localEcho *echo.Echo = echo.New()

type apiResponse struct {
	RequestId string      `json:"requestId"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

func initializeContext(t *testing.T, method string, url string, reqdata interface{}) echo.Context {
	c := config.New("apiserver")
	c.AddConfig(defaultConfigs)
	r, err := registry.New(c)
	if err != nil {
		t.Fatal(err)
	}

	body, _ := json.Marshal(reqdata)
	req := httptest.NewRequest(method, url, bytes.NewReader(body))
	rr := httptest.NewRecorder()
	ctx := localEcho.NewContext(req, rr)
	ctx.Set("registry", r)
	ctx.Set("config", c)
	return ctx
}

func Test_registerTenant(t *testing.T) {
	req := struct {
		TenantId string `json:"tenantId"`
		Password string `json:"password"`
	}{
		TenantId: "jenson",
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/tenants", req)
	defer closeRegistry(ctx)
	registerTenant(ctx)
	r := ctx.Response()
	if r.Status != http.StatusOK {
		t.Errorf("exepected status 200, return %d", r.Status)
	}
}

func Test_loginTenant(t *testing.T) {
	req := struct {
		TenantId string `json:"tenantId"`
		Password string `json:"password"`
	}{
		TenantId: "jenson",
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/tenants/login", req)
	defer closeRegistry(ctx)
	loginTenant(ctx)
	r := ctx.Response()
	if r.Status != http.StatusOK {
		t.Errorf("exepected status 200, return %d", r.Status)
	}
}

func Test_logoutTenant(t *testing.T) {
	req := struct {
		TenantId string `json:"tenantId"`
		Password string `json:"password"`
	}{
		TenantId: "jenson",
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/tenants/logout", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	logoutTenant(ctx)
	r := ctx.Response()
	if r.Status != http.StatusOK {
		t.Errorf("exepected status 200, return %d", r.Status)
	}
}
