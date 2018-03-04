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

	web "github.com/cloustone/sentel/pkg/goshiro/web"
)

var mngApiPolicies = []web.ApiAuthorizePolicy{
	// v1api.CreateProduct
	{
		Url:          "/iot/api/v1/products",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "user",
	},
	// v1api.UpdateProduct
	{
		Url:          "/iot/api/v1/products/:productId",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "user",
	},
	// v1api.GetProductDevices
	{
		Url:          "/iot/api/v1/products/:productId/devices",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.CreateDevice
	{
		Url:          "/iot/api/v1/products/:productId/devices",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.BulkApplyDevices
	{
		Url:          "/iot/api/v1/products/:productId/devices/bulk",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.BulkApplyGetStatus
	{
		Url:          "/iot/api/v1/products/:productId/devices/bulk/:deviceId",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.BulkApplyGetDevices
	{
		Url:          "/iot/api/v1/products/:productId/devices/bulk",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.BulkApplyGetDeviceStatus
	{
		Url:          "/iot/api/v1/products/:productId/devices/bulk/status",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.GetDeviceList
	{
		Url:          "/iot/api/v1/products/:productId/devices",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.GetDeviceByName
	{
		Url:          "/iot/api/v1/products/:productId/devices/:deviceName",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},

	// v1api.SaveDevicePropsByName
	{
		Url:          "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		AllowedRoles: "user",
	},
	// v1api.GetDevicePropsByName
	{
		Url:          "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		AllowedRoles: "user",
	},
	// v1api.RemoveDevicePropsByName
	{
		Url:          "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Method:       http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		AllowedRoles: "user",
	},

	// v1api.UpdateShadowDevice
	{
		Url:          "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.GetShadowDevice
	{
		Url:          "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Method:       http.MethodGet,
		Resource:     "/tenants/:$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.SendMessageToDevice
	{
		Url:          "/iot/api/v1/message",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices/$DeviceId/message",
		AllowedRoles: "user",
	},
	// v1api.BroadcastProductMessage
	{
		Url:          "/iot/api/v1/message/broadcast",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/message",
		AllowedRoles: "user",
	},
}
