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

package management

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloustone/sentel/apiserver/v1api"
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

const (
	tenantId    = "jenson"
	accessId    = "jenson"
	deviceId    = "device_01"
	deviceName  = "device_name_0"
	props       = "props_01"
	productName = "test_product1"
)

var (
	localEcho *echo.Echo = echo.New()
	productId string     = ""
	token     string     = ""
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

func Test_createProduct(t *testing.T) {
	req := struct {
		ProductName string `json:"productName"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}{
		ProductName: productName,
		Category:    "home",
		Description: "test",
	}

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/products", req)
	defer closeRegistry(ctx)
	v1api.CreateProduct(ctx)
	if rsp, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	} else {
		p := rsp.Result.(map[string]interface{})
		if p["TenantId"] != tenantId || p["ProductName"] != productName {
			t.Error("tenantId and productName miss")
		}
		productId = p["ProductId"].(string) // Update global productId
	}
}

func Test_updateProduct(t *testing.T) {
	req := struct {
		ProductId   string `json:"productId"`
		ProductName string `json:"productName"`
		Category    string `json:"category"`
		Description string `json:"description"`
	}{
		ProductId:   productId,
		ProductName: productName,
		Category:    "home",
		Description: "updated product",
	}

	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/products/:productId", req)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.UpdateProduct(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getProductDevices(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/products/:productId/devices", nil)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.GetProductDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

// Device
func Test_createDevice(t *testing.T) {
	req := struct {
		DeviceName string `json:"DeviceName"`
		ProductId  string `json:"ProductId"`
		DeviceId   string `json:"DeviceId"`
	}{
		DeviceName: "test device",
		ProductId:  productId,
	}

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/products/:productId/devices", req)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.CreateDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_bulkApplyDevices(t *testing.T) {
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/products/:productId/devices/bulk", nil)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.BulkApplyDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_bulkApplyGetStatus(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/products/:productId/devices/bulk/:deviceId", nil)
	ctx.SetParamNames("productId", "deviceId")
	ctx.SetParamValues(productId, deviceId)
	defer closeRegistry(ctx)
	v1api.BulkApplyGetStatus(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_bulkApplyGetDevices(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/products/:productId/devices/bulk", nil)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.BulkApplyGetDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_getDeviceList(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/products/:productId/devices", nil)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.GetDeviceList(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_bulkGetDeviceStatus(t *testing.T) {
	url := "/iot/api/v1/products/:productId/devices/bulk/status"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	ctx.SetParamNames("productId")
	ctx.SetParamValues(productId)
	defer closeRegistry(ctx)
	v1api.BulkGetDeviceStatus(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getDeviceByName(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/products/:productId/devices/:deviceName", nil)
	ctx.SetParamNames("productId", "deviceName")
	ctx.SetParamValues(productId, deviceName)
	defer closeRegistry(ctx)
	v1api.GetDeviceByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_saveDevicePropsByName(t *testing.T) {
	url := "/iot/api/v1/products/:productId/devices/:deviceId/props"
	ctx := initializeContext(t, http.MethodPost, url, nil)
	ctx.SetParamNames("productId", "devicId")
	ctx.SetParamValues(productId, deviceId)
	defer closeRegistry(ctx)
	v1api.SaveDevicePropsByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getDevicePropsByName(t *testing.T) {
	url := "/iot/api/v1/products/:productId/devices/:deviceId/props/:props"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	ctx.SetParamNames("productId", "devicId", "props")
	ctx.SetParamValues(productId, deviceId, props)
	defer closeRegistry(ctx)
	v1api.GetDevicePropsByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_removeDevicePropsByName(t *testing.T) {
	url := "/iot/api/v1/products/:productId/devices/:deviceId/props/:props"
	ctx := initializeContext(t, http.MethodDelete, url, nil)
	ctx.SetParamNames("productId", "devicId", "props")
	ctx.SetParamValues(productId, deviceId, props)
	defer closeRegistry(ctx)
	v1api.RemoveDevicePropsByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_sendMessageToDevice(t *testing.T) {
	req := struct {
		ProductId string `json:"productId"`
		DeviceId  string `json:"deviceId"`
		Topic     string `json:"topic"`
		Payload   []byte `json:"payload"`
		Qos       uint8  `json:"qos"`
		Retain    bool   `json:"retain"`
	}{
		ProductId: productId,
		DeviceId:  deviceId,
		Topic:     "hello",
		Payload:   []byte("world"),
	}

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/message", req)
	defer closeRegistry(ctx)
	v1api.SendMessageToDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_broadcastProductMessage(t *testing.T) {
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/message/broadcast", nil)
	defer closeRegistry(ctx)
	v1api.SendMessageToDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_updateShadowDevice(t *testing.T) {
	url := "/iot/api/v1/products/:productId/devices/:deviceId/shadow"
	ctx := initializeContext(t, http.MethodPatch, url, nil)
	ctx.SetParamNames("productId", "devicId")
	ctx.SetParamValues(productId, deviceId)
	defer closeRegistry(ctx)
	v1api.UpdateShadowDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getShadowDevice(t *testing.T) {
	url := "/iot/api/v1/products/:productId/devices/:deviceId/shadow"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	ctx.SetParamNames("productId", "devicId")
	ctx.SetParamValues(productId, deviceId)
	defer closeRegistry(ctx)
	v1api.GetShadowDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
