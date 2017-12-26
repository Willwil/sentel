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
	"github.com/cloustone/sentel/apiserver/v1api"
	"github.com/cloustone/sentel/keystone/auth"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

func setAccessId(ctx echo.Context) {
	// After authenticated by gateway,
	// the authentication paramters must bevalid
	param := auth.ApiAuthParam{}
	ctx.Bind(&param)
	ctx.Set("AccessId", param.AccessId)
	ctx.Set("RequestId", uuid.NewV4().String())
}

// Product
func createProduct(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.CreateProduct(ctx)
}
func removeProduct(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.RemoveProduct(ctx)
}
func updateProduct(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.UpdateProduct(ctx)
}
func getProductList(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetProductList(ctx)
}
func getProduct(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetProduct(ctx)
}
func getProductDevices(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetProductDevices(ctx)
}

func getDeviceStatics(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetDeviceStatics(ctx)
}

func registerDevice(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.RegisterDevice(ctx)
}
func bulkApplyDevices(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.BulkApplyDevices(ctx)
}
func bulkApplyGetStatus(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.BulkApplyGetStatus(ctx)
}
func bulkApplyGetDevices(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.BulkApplyGetDevices(ctx)
}
func getDeviceList(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetDeviceList(ctx)
}
func bulkGetDeviceStatus(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.BulkGetDeviceStatus(ctx)
}
func getDeviceByName(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetDeviceByName(ctx)
}
func saveDevicePropsByName(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.SaveDevicePropsByName(ctx)
}
func getDevicePropsByName(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetDevicePropsByName(ctx)
}
func getShadowDevice(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetShadowDevice(ctx)
}
func updateShadowDevice(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.UpdateShadowDevice(ctx)
}

// Runtime & Service
func sendMessageToDevice(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.SendMessageToDevice(ctx)
}
func broadcastProductMessage(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.BroadcastProductMessage(ctx)
}
func getServiceStatics(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetServiceStatics(ctx)
}