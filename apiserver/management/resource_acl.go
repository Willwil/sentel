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

import "github.com/cloustone/sentel/apiserver/base"

var mngApiPolicies = []base.ResourceAC{
	// v1api.CreateProduct
	{
		Path:     "/iot/api/v1/products",
		Resource: "/products",
	},
	// v1api.UpdateProduct
	{
		Path:     "/iot/api/v1/products/:productId",
		Resource: "/products",
	},
	// v1api.GetProductDevices
	{
		Path:     "/iot/api/v1/products/:productId/devices",
		Resource: "/products/:productId",
	},
	// v1api.CreateDevice
	{
		Path:     "/iot/api/v1/products/:productId/devices",
		Resource: "/products/:productId",
	},
	// v1api.BulkApplyDevices
	{
		Path:     "/iot/api/v1/products/:productId/devices/bulk",
		Resource: "/products/:productId",
	},
	// v1api.BulkApplyGetStatus
	{
		Path:     "/iot/api/v1/products/:productId/devices/bulk/:deviceId",
		Resource: "/products/:productId",
	},
	// v1api.BulkApplyGetDevices
	{
		Path:     "/iot/api/v1/products/:productId/devices/bulk",
		Resource: "/products/:productId",
	},
	// v1api.BulkApplyGetDeviceStatus
	{
		Path:     "/iot/api/v1/products/:productId/devices/bulk/status",
		Resource: "/products/:productId",
	},
	// v1api.GetDeviceList
	{
		Path:     "/iot/api/v1/products/:productId/devices",
		Resource: "/products/:productId",
	},
	// v1api.GetDeviceByName
	{
		Path:     "/iot/api/v1/products/:productId/devices/:deviceName",
		Resource: "/products/:productId",
	},

	// v1api.SaveDevicePropsByName
	{
		Path:     "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Resource: "/products/:productId",
	},
	// v1api.GetDevicePropsByName
	{
		Path:     "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Resource: "/products/:productId",
	},
	// v1api.RemoveDevicePropsByName
	{
		Path:     "/iot/api/v1/products/:productId/devics/:deviceId/props",
		Resource: "/products/:productId",
	},

	// v1api.UpdateShadowDevice
	{
		Path:     "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Resource: "/products/:productId",
	},
	// v1api.GetShadowDevice
	{
		Path:     "/iot/api/v1/products/:productId/devices/:deviceId/shadow",
		Resource: "/products/:productId",
	},
	// v1api.SendMessageToDevice
	{
		Path:     "/iot/api/v1/message",
		Resource: "/products/:productId",
	},
	// v1api.BroadcastProductMessage
	{
		Path:     "/iot/api/v1/message/broadcast",
		Resource: "/products/:productId",
	},
}
