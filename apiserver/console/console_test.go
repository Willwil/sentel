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
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	jwt "github.com/dgrijalva/jwt-go"
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
	accessId        = "jenson"
	tenantId        = "jenson"
	productName     = "test_product1"
	ruleName        = "test_rule"
	topicFlavorName = "test_topicflavor"
)

var (
	localEcho *echo.Echo = echo.New()
	productId string     = ""
	deviceId  string     = ""
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
	token := jwt.Token{
		Claims: &base.JwtApiClaims{
			AccessId: accessId,
		},
	}
	ctx.Set("user", &token)
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

func Test_registerTenant(t *testing.T) {
	req := struct {
		TenantId string `json:"tenantId"`
		Password string `json:"password"`
	}{
		TenantId: tenantId,
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/tenants", req)
	defer closeRegistry(ctx)
	registerTenant(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_loginTenant(t *testing.T) {
	req := struct {
		TenantId string `json:"tenantId"`
		Password string `json:"password"`
	}{
		TenantId: tenantId,
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/tenants/login", req)
	defer closeRegistry(ctx)
	loginTenant(ctx)
	if rsp, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	} else {
		token = rsp.Result.(echo.Map)["token"].(string)
	}
}

func Test_getTenant(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/console/tenants/"+tenantId, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getTenant(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_updateTenant(t *testing.T) {
	req := registry.Tenant{
		TenantId: tenantId,
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/console/tenants/logout", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	updateTenant(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
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

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/products", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
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

	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/console/products/"+productId, req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	updateProduct(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getProductList(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/console/products", nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getProductList(ctx)
	if rsp, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	} else {
		type product struct {
			ProductId   string `json:"productId"`
			ProductName string `json:"productName"`
		}
		products := rsp.Result.([]interface{})
		if len(products) == 0 {
			t.Error("product number is zero")
		}
	}
}

func Test_getProduct(t *testing.T) {
	ctx := initializeContext(t, http.MethodGet, "/iot/api/v1/console/products/"+productId, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getProduct(ctx)
	if rsp, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	} else {
		p := rsp.Result.(registry.Product)
		if p.ProductId != productId {
			t.Error("wrong product id")
		}
	}
}

func Test_getProductDevices(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "/devices"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getProductDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getProductRules(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "/rules"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getProductRules(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getDeviceStatics(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "/devices/statics"
	ctx := initializeContext(t, http.MethodPost, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getDeviceStatics(ctx)
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

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/devices", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	createDevice(ctx)
	if rsp, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	} else {
		device := rsp.Result.(registry.Device)
		if device.DeviceName != req.DeviceName || device.ProductId != req.ProductId {
			t.Error("wrong device retrived")
		}
		deviceId = device.DeviceId
	}
}

func Test_getOneDevice(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "devices/" + deviceId
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getOneDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_updateDevice(t *testing.T) {
	req := struct {
		DeviceName string `json:"DeviceName"`
		ProductId  string `json:"ProductId"`
		DeviceId   string `json:"DeviceId"`
	}{
		DeviceName: "updated device name",
		ProductId:  productId,
		DeviceId:   deviceId,
	}
	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/console/devices", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	updateDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_bulkRegisterDevices(t *testing.T) {
	req := struct {
		DeviceName string `json:"deviceName"`
		ProductId  string `json:"productId"`
		Number     string `json:"number"`
	}{
		DeviceName: "test device",
		ProductId:  productId,
		Number:     "10",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/devices/bulk", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	bulkRegisterDevices(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_updateShadowDevice(t *testing.T) {
	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/console/device/shadow", nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	updateShadowDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getShadowDevice(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "devices/" + deviceId + "/shadow"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getShadowDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

// Rule
func Test_createRule(t *testing.T) {
	rule1 := registry.Rule{
		ProductId:   productId,
		RuleName:    ruleName,
		DataFormat:  "json",
		Description: "test rule1",
		DataProcess: registry.RuleDataProcess{
			Topic:     "hello",
			Condition: "",
			Fields:    []string{"field1", "field2"},
		},
		DataTarget: registry.RuleDataTarget{
			Type:  registry.DataTargetTypeTopic,
			Topic: "world",
		},
	}

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/rules", rule1)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	createRule(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_updateRule(t *testing.T) {
	rule1 := registry.Rule{
		ProductId:   productId,
		RuleName:    ruleName,
		DataFormat:  "json",
		Description: "updated test rule1",
		DataProcess: registry.RuleDataProcess{
			Topic:     "hello",
			Condition: "",
			Fields:    []string{"field1", "field2"},
		},
		DataTarget: registry.RuleDataTarget{
			Type:  registry.DataTargetTypeTopic,
			Topic: "world",
		},
	}
	ctx := initializeContext(t, http.MethodPatch, "/iot/api/v1/console/rules", rule1)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	removeRule(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_startRule(t *testing.T) {
	rule1 := registry.Rule{
		ProductId: productId,
		RuleName:  ruleName,
	}

	ctx := initializeContext(t, http.MethodPut, "/iot/api/v1/console/rules/start", rule1)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	startRule(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_stopRule(t *testing.T) {
	rule1 := registry.Rule{
		ProductId: productId,
		RuleName:  ruleName,
	}

	ctx := initializeContext(t, http.MethodPut, "/iot/api/v1/console/rules/stop", rule1)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	stopRule(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getRule(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "/rules/" + ruleName
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getRule(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

// Topic Flavor
func Test_createTopicFlavor(t *testing.T) {
	req := registry.TopicFlavor{
		FlavorName: topicFlavorName,
		Builtin:    false,
		TenantId:   tenantId,
		Topics: map[string]uint8{
			"hello": 1,
			"world": 0,
		},
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/topicflavors", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	createTopicFlavor(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getProductTopicFlavors(t *testing.T) {
	url := "/iot/api/v1/console/topicflavors/" + productId
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getProductTopicFlavors(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getTenantTopicFlavors(t *testing.T) {
	url := "/iot/api/v1/console/topicflavors/tenants/" + tenantId
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getTenantTopicFlavors(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getBuiltinTopicFlavors(t *testing.T) {
	url := "/iot/api/v1/console/topicflavors/builtin"
	ctx := initializeContext(t, http.MethodGet, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	getBuiltinTopicFlavors(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_setProductTopicFlavor(t *testing.T) {
	url := "/iot/api/v1/console/topicflavors/" + productId + "?flavor=" + topicFlavorName
	ctx := initializeContext(t, http.MethodPut, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	setProductTopicFlavor(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

// Runtime
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

	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/message", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	sendMessageToDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_broadcastProductMessage(t *testing.T) {
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/message/broadcast", nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	sendMessageToDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_getServiceStatics(t *testing.T) {
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
	getServiceStatics(ctx)
	r := ctx.Response()
	if r.Status != http.StatusOK {
		t.Errorf("exepected status 200, return %d", r.Status)
	}
}

func Test_removeProductTopicFlavor(t *testing.T) {
	url := "/iot/api/v1/console/topicflavors/" + productId
	ctx := initializeContext(t, http.MethodDelete, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	removeProductTopicFlavor(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
func Test_removeRule(t *testing.T) {
	url := "/iot/api/v1/console/" + productId + "/rules/" + ruleName
	ctx := initializeContext(t, http.MethodDelete, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	removeRule(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_removeDevice(t *testing.T) {
	url := "/iot/api/v1/console/products/" + productId + "devices/" + deviceId
	ctx := initializeContext(t, http.MethodDelete, url, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	removeDevice(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_removeProduct(t *testing.T) {
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/products/"+productId, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	removeProduct(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_logoutTenant(t *testing.T) {
	req := struct {
		TenantId string `json:"tenantId"`
		Password string `json:"password"`
	}{
		TenantId: tenantId,
		Password: "default",
	}
	ctx := initializeContext(t, http.MethodPost, "/iot/api/v1/console/tenants/logout", req)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	logoutTenant(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}

func Test_deleteTenant(t *testing.T) {
	ctx := initializeContext(t, http.MethodDelete, "/iot/api/v1/console/tenants/"+tenantId, nil)
	defer closeRegistry(ctx)
	ctx.Set("AccessId", accessId)
	deleteTenant(ctx)
	if _, err := getApiResponse(ctx); err != nil {
		t.Error(err)
	}
}
