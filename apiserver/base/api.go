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
	"github.com/cloustone/sentel/common"
	jwt "github.com/dgrijalva/jwt-go"

	echo "github.com/labstack/echo"
)

type ApiContext struct {
	echo.Context
	Config com.Config
}

type ApiResponse struct {
	RequestId string      `json:"requestID"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

type JwtApiClaims struct {
	jwt.StandardClaims
	Name string `json:"name"`
}

type ConsoleApi interface {
	// Tenant
	RegisterTenant(ctx echo.Context) error
	LoginTenant(ctx echo.Context) error
	LogoutTenant(ctx echo.Context) error
	DeleteTenant(ctx echo.Context) error
	UpdateTenant(ctx echo.Context) error
	GetTenant(ctx echo.Context) error
	GetTenantProductList(ctx echo.Context) error

	// Product
	RegisterProduct(ctx echo.Context) error
	DeleteProduct(ctx echo.Context) error
	UpdateProduct(ctx echo.Context) error
	GetProduct(ctx echo.Context) error
	GetProductDevices(ctx echo.Context) error

	// Device
	RegisterDevice(ctx echo.Context) error
	BulkRegisterDevices(ctx echo.Context) error
	GetOneDevice(ctx echo.Context) error
	DeleteDevice(ctx echo.Context) error
	UpdateDevice(ctx echo.Context) error

	// Rule
	CreateRule(ctx echo.Context) error
	GetRule(ctx echo.Context) error
	RemoveRule(ctx echo.Context) error
	UpdateRule(ctx echo.Context) error
	StartRule(ctx echo.Context) error
	StopRule(ctx echo.Context) error
	GetProductRules(ctx echo.Context) error

	// Shadow Device
	GetShadowDevice(ctx echo.Context) error
	UpdateShadowDevice(ctx echo.Context) error

	// Message
	SendMessageToDevice(ctx echo.Context) error
	BroadcaseProductMessage(ctx echo.Context) error
	// Statitics
	GetDeviceStatics(ctx echo.Context) error
	GetServiceStatics(ctx echo.Context) error
}

type ManagementApi interface {
}
