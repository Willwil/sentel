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
	"net/http"

	"github.com/cloustone/sentel/apiserver/base"
)

var consoleResources = []base.Resource{
	// v1api.LogoutTenant
	{
		Url:          "/iot/api/v1/console/tenants/logout",
		Method:       http.MethodPost,
		Resource:     "/logout",
		AllowedRoles: "user",
	},
	// v1api.DeleteTenant
	{
		Url:          "/iot/api/v1/console/tenants/:tenantId",
		Method:       http.MethodDelete,
		Resource:     "/tenants/:tenantId",
		AllowedRoles: "user",
	},
	// v1api.GetTenant
	{
		Url:          "/iot/api/v1/console/tenants/:tenantId",
		Method:       http.MethodGet,
		Resource:     "/tenants/:tenantId",
		AllowedRoles: "user",
	},
	// v1api.UpdateTenant
	{
		Url:          "/iot/api/v1/console/tenants",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId",
		AllowedRoles: "user",
	},

	// v1api.CreateProduct
	{
		Url:          "/iot/api/v1/console/products",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "user",
	},
	// v1api.RemoveProduct
	{
		Url:          "/iot/api/v1/console/products/:productId",
		Method:       http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/:productId",
		AllowedRoles: "user",
	},
	// v1api.UpdateProduct
	{
		Url:          "/iot/api/v1/console/products",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "user",
	},
	// v1api.GetProductList
	{
		Url:          "/iot/api/v1/console/products",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "user",
	},
	// v1api.GetProduct
	{
		Url:          "/iot/api/v1/console/products/:productId",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId",
		AllowedRoles: "user",
	},
	// v1api.GetProductDevices
	{
		Url:          "/iot/api/v1/console/products/:productId/devices",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.GetProductRules
	{
		Url:          "/iot/api/v1/console/procuts/:productId/rules",
		Method:       http.MethodGet,
		Resource:     "/tenants/$accessId/products/:productId/rules",
		AllowedRoles: "user",
	},
	// v1api.GetDeviceStatics
	{
		Url:          "/iot/api/v1/console/products/:productId/devices/statics",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.CreateDevice
	{
		Url:          "/iot/api/v1/console/devices",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices",
		AllowedRoles: "user",
	},
	// v1api.GetOneDevice
	{
		Url:          "/iot/api/v1/console/products/:productId/devics/:deviceId",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.RemoveDevice
	{
		Url:          "/iot/api/v1/console/products/:productId/devicds/:deviceId",
		Method:       http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.UpdateDevice
	{
		Url:          "/iot/api/v1/console/devices",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices",
		AllowedRoles: "user",
	},
	// v1api.BulkRegisterDevices
	{
		Url:          "/iot/api/v1/console/devices/bulk",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices",
		AllowedRoles: "user",
	},
	// v1api.UpdateShadowDevice
	{
		Url:          "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.GetShadowDevice
	{
		Url:          "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Method:       http.MethodGet,
		Resource:     "/tenants/:$AccessId/products/:productId/devices",
		AllowedRoles: "user",
	},
	// v1api.CreateRule
	{
		Url:          "/iot/api/v1/console/rules",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "user",
	},
	// v1api.RemoveRule
	{
		Url:          "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Method:       http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "user",
	},
	// v1api.UpdateRule
	{
		Url:          "/iot/api/v1/console/rules",
		Method:       http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "user",
	},
	// v1api.StopRule
	{
		Url:          "/iot/api/v1/console/rules/stop",
		Method:       http.MethodPut,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "user",
	},
	// v1api.GetRule
	{
		Url:          "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "user",
	},
	// v1api.CreateTopicFlavor
	{
		Url:          "/iot/api/v1/console/topicflavors",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "user",
	},
	// v1api.RemoveTopicFlavor
	{
		Url:          "/iot/api/v1/console/topicflavors/products/:productId",
		Method:       http.MethodDelete,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "user",
	},
	// v1api.GetProductTopicFlavor
	{
		Url:          "/iot/api/v1/console/topicflavors/products/:productId",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "user",
	},
	// v1api.GetTenantTopicFlavors
	{
		Url:          "/iot/api/v1/console/topicflavors/tenants/:tenantId",
		Method:       http.MethodGet,
		Resource:     "/tenants/:tenantId/topicflavors",
		AllowedRoles: "user",
	},
	// v1api.GetBuiltinTopicFlavors
	{
		Url:          "/iot/api/v1/console/topicflavors/builtin",
		Method:       http.MethodGet,
		Resource:     "/topicflavors/builtin",
		AllowedRoles: "user",
	},
	// v1api.SetProductTopicFlavor
	{
		Url:          "/iot/api/v1/console/topicflavors/:productId?flavor=:topicflavor",
		Method:       http.MethodPut,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "user",
	},
	// v1api.SendMessageToDevice
	{
		Url:          "/iot/api/v1/console/message",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices/$DeviceId/message",
		AllowedRoles: "user",
	},
	// v1api.BroadcastProductMessage
	{
		Url:          "/iot/api/v1/console/message/broadcast",
		Method:       http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/message",
		AllowedRoles: "user",
	},
	// v1api.GetServiceStatics
	{
		Url:          "/iot/api/v1/console/service",
		Method:       http.MethodGet,
		Resource:     "/tenants/$AccessId/services",
		AllowedRoles: "user",
	},
}
