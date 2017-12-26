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
	"github.com/cloustone/sentel/apiserver/v1api"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

func setAccessId(ctx echo.Context) {
	user := ctx.Get("user").(*jwt.Token)
	claims := user.Claims.(jwt.MapClaims)
	accessId := claims["accessId"].(string)
	ctx.Set("AccessId", accessId)
	ctx.Set("RequestId", uuid.NewV4().String())
}

func registerTenant(ctx echo.Context) error {
	return v1api.RegisterTenant(ctx)
}

func loginTenant(ctx echo.Context) error {
	return v1api.LoginTenant(ctx)
}

func logoutTenant(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.LogoutTenant(ctx)
}

func deleteTenant(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.DeleteTenant(ctx)
}

func getTenant(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetTenant(ctx)
}

func updateTenant(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.UpdateTenant(ctx)
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
func getProductRules(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetProductRules(ctx)
}
func getDeviceStatics(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetDeviceStatics(ctx)
}

// Rules
func createRule(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.CreateRule(ctx)
}
func removeRule(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.RemoveRule(ctx)
}
func updateRule(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.UpdateRule(ctx)
}
func startRule(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.StartRule(ctx)
}
func stopRule(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.StopRule(ctx)
}
func getRule(ctx echo.Context) error {
	setAccessId(ctx)
	return v1api.GetRule(ctx)
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
