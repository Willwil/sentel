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
	rr := httptest.NewRecorder()
	ctx := localEcho.NewContext(req, rr)
	ctx.Set("registry", r)
	ctx.Set("config", c)
	ctx.Set("AccessId", accessId)
	return ctx
}

func getApiResponse(ctx echo.Context) (apiResponse, error) {
	rsp := apiResponse{}
	rr := ctx.Response().Writer.(*httptest.ResponseRecorder)
	if rr.Code != http.StatusOK {
		return rsp, fmt.Errorf("exepected status 200, return %d", rr.Code)
	} else {
		// unmarshal product detail
		body, _ := ioutil.ReadAll(rr.Body)
		if err := json.Unmarshal(body, &rsp); err != nil {
			return rsp, err
		}
		return rsp, nil
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
	createProduct(ctx)
	if rsp, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	} else {
		p := rsp.Result.(registry.Product)
		if p.TenantId != tenantId || p.ProductName != productName {
			t.Error("tenantId and productName miss")
		}
		productId = p.ProductId // Update global productId
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

	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/products/"+productId, req)
	updateProduct(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getProductDevices(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	getProductDevices(ctx)
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

	url := "/iot/api/v1/products/" + productId + "/devices"
	ctx := initializeContext(t, http.MethodPost, url, req)
	createDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_bulkApplyDevices(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/bulk"
	ctx := initializeContext(t, http.MethodPost, url, nil)
	bulkApplyDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_bulkApplyGetStatus(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/bulk/" + deviceId
	ctx := initializeContext(t, http.MethodGet, url, nil)
	bulkApplyGetStatus(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_bulkApplyGetDevices(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/bulk"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	bulkApplyGetDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_getDeviceList(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	getDeviceList(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_bulkGetDeviceStatus(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/bulk/status"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	bulkGetDeviceStatus(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getDeviceByName(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/" + deviceName
	ctx := initializeContext(t, http.MethodGet, url, nil)
	getDeviceByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_saveDevicePropsByName(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/" + deviceId + "/props"
	ctx := initializeContext(t, http.MethodPost, url, nil)
	saveDevicePropsByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getDevicePropsByName(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/" + deviceId + "/props/" + props
	ctx := initializeContext(t, http.MethodGet, url, nil)
	getDevicePropsByName(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_removeDevicePropsByName(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/" + deviceId + "/props/" + props
	ctx := initializeContext(t, http.MethodDelete, url, nil)
	removeDevicePropsByName(ctx)
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
	sendMessageToDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_broadcastProductMessage(t *testing.T) {
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/message/broadcast", nil)
	sendMessageToDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_updateShadowDevice(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/" + deviceId + "/shadow"
	ctx := initializeContext(t, http.MethodPatch, url, nil)
	updateShadowDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getShadowDevice(t *testing.T) {
	url := "/iot/api/v1/products/" + productId + "/devices/" + deviceId + "/shadow"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	getShadowDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
