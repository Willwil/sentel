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

package apiservice

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/labstack/echo"
)

var defaultConfigs = config.M{
	"mns": {
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

const (
	accessId = "jenson"
	tenantId = "jenson"
)

var (
	localEcho   *echo.Echo = echo.New()
	token       string     = ""
	securityMgr shiro.SecurityManager
)

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

	var req *http.Request
	if reqdata != nil {
		body, _ := json.Marshal(reqdata)
		req = httptest.NewRequest(method, url, bytes.NewReader(body))
	} else {
		req = httptest.NewRequest(method, url, nil)
	}
	req.Header.Set("Content-Type", "application/json")
	rr := httptest.NewRecorder()
	ctx := localEcho.NewContext(req, rr)
	ctx.Set("registry", r)
	ctx.Set("config", c)
	ctx.Set("AccessId", accessId)
	// Initialize security Manager
	if securityMgr == nil {
		for _, res := range apiPolicies {
			resourceMaps[res.Path] = res.Resource
		}
		realm, _ := base.NewAuthorizeRealm(c)
		mgr, _ := goshiro.NewSecurityManager(c, realm)
		securityMgr = mgr
		securityMgr.Load()
	}
	ctx.Set("SecurityManager", securityMgr)
	return ctx
}

func closeRegistry(ctx echo.Context) {
	r := ctx.Get("registry").(*registry.Registry)
	r.Close()
}

func getApiResponse(ctx echo.Context) (apiResponse, error) {
	rsp := apiResponse{}
	rr := ctx.Response().Writer.(*httptest.ResponseRecorder)
	if rr.Code != http.StatusOK {
		return rsp, fmt.Errorf("exepected status 200, return %d", rr.Code)
	} else {
		// unmarshal product detail
		result := rr.Result()
		if result.StatusCode != http.StatusOK {
			return rsp, fmt.Errorf("expected status 200, result %d", result.StatusCode)
		} else {
			body, _ := ioutil.ReadAll(result.Body)
			if err := json.Unmarshal(body, &rsp); err != nil {
				return rsp, err
			}
			return rsp, nil
		}

	}
}
