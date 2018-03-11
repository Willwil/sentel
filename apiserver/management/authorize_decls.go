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
	"net/http"

	"github.com/cloustone/sentel/pkg/goshiro/shiro"
)

var mngApiPolicies = []shiro.AuthorizePolicy{
	// v1api.CreateProduct
	{
		Path:         "/iot/api/v1/products",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "tenant",
	},
	// v1api.UpdateProduct
	{
		Path:         "/iot/api/v1/products/:productId",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "tenant",
	},
	// v1api.GetProductDevices
	{
		Path:         "/iot/api/v1/products/:productId/devices",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.CreateDevice
	{
		Path:         "/iot/api/v1/products/:productId/devices",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.BulkApplyDevices
	{
		Path:         "/iot/api/v1/products/:productId/devices/bulk",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.BulkApplyGetStatus
	{
		Path:         "/iot/api/v1/products/:productId/devices/bulk/:deviceId",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.BulkApplyGetDevices
	{
		Path:         "/iot/api/v1/products/:productId/devices/bulk",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.BulkApplyGetDeviceStatus
	{
		Path:         "/iot/api/v1/products/:productId/devices/bulk/status",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.GetDeviceList
	{
		Path:         "/iot/api/v1/products/:productId/devices",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.GetDeviceByName
	{
		Path:         "/iot/api/v1/products/:productId/devices/:deviceName",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},

	// v1api.SaveDevicePropsByName
	{
		Path:         "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		AllowedRoles: "tenant",
	},
	// v1api.GetDevicePropsByName
	{
		Path:         "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		AllowedRoles: "tenant",
	},
	// v1api.RemoveDevicePropsByName
	{
		Path:         "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Methods:      http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		AllowedRoles: "tenant",
	},

	// v1api.UpdateShadowDevice
	{
		Path:         "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.GetShadowDevice
	{
		Path:         "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Methods:      http.MethodGet,
		Resource:     "/tenants/:$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.SendMessageToDevice
	{
		Path:         "/iot/api/v1/message",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices/$DeviceId/message",
		AllowedRoles: "tenant",
	},
	// v1api.BroadcastProductMessage
	{
		Path:         "/iot/api/v1/message/broadcast",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/message",
		AllowedRoles: "tenant",
	},
}
