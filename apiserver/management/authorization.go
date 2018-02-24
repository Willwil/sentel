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

	"github.com/cloustone/sentel/keystone/subject"
)

var authorizations = []subject.Declaration{
	// v1api.CreateProduct
	{
		Url:     "/iot/api/v1/products",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products",
		Action:  "create",
	},
	// v1api.UpdateProduct
	{
		Url:     "/iot/api/v1/products/:productId",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId/products",
		Action:  "write",
	},
	// v1api.GetProductDevices
	{
		Url:     "/iot/api/v1/products/:productId/devices",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.CreateDevice
	{
		Url:     "/iot/api/v1/products/:productId/devices",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "create",
	},
	// v1api.BulkApplyDevices
	{
		Url:     "/iot/api/v1/products/:productId/devices/bulk",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.BulkApplyGetStatus
	{
		Url:     "/iot/api/v1/products/:productId/devices/bulk/:deviceId",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.BulkApplyGetDevices
	{
		Url:     "/iot/api/v1/products/:productId/devices/bulk",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.BulkApplyGetDeviceStatus
	{
		Url:     "/iot/api/v1/products/:productId/devices/bulk/status",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.GetDeviceList
	{
		Url:     "/iot/api/v1/products/:productId/devices",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.GetDeviceByName
	{
		Url:     "/iot/api/v1/products/:productId/devices/:deviceName",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},

	// v1api.SaveDevicePropsByName
	{
		Url:     "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		Action:  "create",
	},
	// v1api.GetDevicePropsByName
	{
		Url:     "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		Action:  "read",
	},
	// v1api.RemoveDevicePropsByName
	{
		Url:     "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Method:  http.MethodDelete,
		Subject: "/tenants/$AccessId/products/:productId/devices/:deviceId/props",
		Action:  "delete",
	},

	// v1api.UpdateShadowDevice
	{
		Url:     "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "write",
	},
	// v1api.GetShadowDevice
	{
		Url:     "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Method:  http.MethodGet,
		Subject: "/tenants/:$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.SendMessageToDevice
	{
		Url:     "/iot/api/v1/message",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/devices/$DeviceId/message",
		Action:  "write",
	},
	// v1api.BroadcastProductMessage
	{
		Url:     "/iot/api/v1/message/broadcast",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/message",
		Action:  "write",
	},
}
