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

	"github.com/cloustone/sentel/keystone/subject"
)

var authorizations = []subject.Declaration{
	// v1api.LogoutTenant
	{
		Url:     "/iot/api/v1/console/tenants/logout",
		Method:  http.MethodPost,
		Subject: "/login",
		Action:  "write",
	},
	// v1api.DeleteTenant
	{
		Url:     "/iot/api/v1/console/tenants/:tenantId",
		Method:  http.MethodDelete,
		Subject: "/tenants/:tenantId",
		Action:  "delete",
	},
	// v1api.GetTenant
	{
		Url:     "/iot/api/v1/console/tenants/:tenantId",
		Method:  http.MethodGet,
		Subject: "/tenants/:tenantId",
		Action:  "read",
	},
	// v1api.UpdateTenant
	{
		Url:     "/iot/api/v1/console/tenants",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId",
		Action:  "write",
	},

	// v1api.CreateProduct
	{
		Url:     "/iot/api/v1/console/products",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products",
		Action:  "create",
	},
	// v1api.RemoveProduct
	{
		Url:     "/iot/api/v1/console/products/:productId",
		Method:  http.MethodDelete,
		Subject: "/tenants/$AccessId/products/:productId",
		Action:  "delete",
	},
	// v1api.UpdateProduct
	{
		Url:     "/iot/api/v1/console/products",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId/products",
		Action:  "write",
	},
	// v1api.GetProductList
	{
		Url:     "/iot/api/v1/console/products",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products",
		Action:  "read",
	},
	// v1api.GetProduct
	{
		Url:     "/iot/api/v1/console/products/:productId",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId",
		Action:  "read",
	},
	// v1api.GetProductDevices
	{
		Url:     "/iot/api/v1/console/products/:productId/devices",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.GetProductRules
	{
		Url:     "/iot/api/v1/console/procuts/:productId/rules",
		Method:  http.MethodGet,
		Subject: "/tenants/$accessId/products/:productId/rules",
		Action:  "read",
	},
	// v1api.GetDeviceStatics
	{
		Url:     "/iot/api/v1/console/products/:productId/devices/statics",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.CreateDevice
	{
		Url:     "/iot/api/v1/console/devices",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/devices",
		Action:  "create",
	},
	// v1api.GetOneDevice
	{
		Url:     "/iot/api/v1/console/products/:productId/devics/:deviceId",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.RemoveDevice
	{
		Url:     "/iot/api/v1/console/products/:productId/devicds/:deviceId",
		Method:  http.MethodDelete,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "delete",
	},
	// v1api.UpdateDevice
	{
		Url:     "/iot/api/v1/console/devices",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId/products/$ProductId/devices",
		Action:  "write",
	},
	// v1api.BulkRegisterDevices
	{
		Url:     "/iot/api/v1/console/devices/bulk",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/devices",
		Action:  "write",
	},
	// v1api.UpdateShadowDevice
	{
		Url:     "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId/products/:productId/devices",
		Action:  "write",
	},
	// v1api.GetShadowDevice
	{
		Url:     "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Method:  http.MethodGet,
		Subject: "/tenants/:$AccessId/products/:productId/devices",
		Action:  "read",
	},
	// v1api.CreateRule
	{
		Url:     "/iot/api/v1/console/rules",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/rules",
		Action:  "create",
	},
	// v1api.RemoveRule
	{
		Url:     "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Method:  http.MethodDelete,
		Subject: "/tenants/$AccessId/products/$ProductId/rules",
		Action:  "delete",
	},
	// v1api.UpdateRule
	{
		Url:     "/iot/api/v1/console/rules",
		Method:  http.MethodPatch,
		Subject: "/tenants/$AccessId/products/$ProductId/rules",
		Action:  "write",
	},
	// v1api.StartRule
	{
		Url:     "/iot/api/v1/console/rules/start",
		Method:  http.MethodPut,
		Subject: "/tenants/$AccessId/products/$ProductId/rules",
		Action:  "write",
	},
	// v1api.StopRule
	{
		Url:     "/iot/api/v1/console/rules/stop",
		Method:  http.MethodPut,
		Subject: "/tenants/$AccessId/products/$ProductId/rules",
		Action:  "write",
	},
	// v1api.GetRule
	{
		Url:     "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/products/$ProductId/rules",
		Action:  "read",
	},
	// v1api.CreateTopicFlavor
	{
		Url:     "/iot/api/v1/console/topicflavors",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/topicflavors",
		Action:  "create",
	},
	// v1api.RemoveTopicFlavor
	{
		Url:     "/iot/api/v1/console/topicflavors/products/:productId",
		Method:  http.MethodDelete,
		Subject: "/tenants/$AccessId/topicflavors",
		Action:  "delete",
	},
	// v1api.GetProductTopicFlavor
	{
		Url:     "/iot/api/v1/console/topicflavors/products/:productId",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/topicflavors",
		Action:  "read",
	},
	// v1api.GetTenantTopicFlavors
	{
		Url:     "/iot/api/v1/console/topicflavors/tenants/:tenantId",
		Method:  http.MethodGet,
		Subject: "/tenants/:tenantId/topicflavors",
		Action:  "read",
	},
	// v1api.GetBuiltinTopicFlavors
	{
		Url:     "/iot/api/v1/console/topicflavors/builtin",
		Method:  http.MethodGet,
		Subject: "/topicflavors/builtin",
		Action:  "read",
	},
	// v1api.SetProductTopicFlavor
	{
		Url:     "/iot/api/v1/console/topicflavors/:productId?flavor=:topicflavor",
		Method:  http.MethodPut,
		Subject: "/tenants/$AccessId/topicflavors",
		Action:  "write",
	},
	// v1api.SendMessageToDevice
	{
		Url:     "/iot/api/v1/console/message",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/devices/$DeviceId/message",
		Action:  "write",
	},
	// v1api.BroadcastProductMessage
	{
		Url:     "/iot/api/v1/console/message/broadcast",
		Method:  http.MethodPost,
		Subject: "/tenants/$AccessId/products/$ProductId/message",
		Action:  "write",
	},
	// v1api.GetServiceStatics
	{
		Url:     "/iot/api/v1/console/service",
		Method:  http.MethodGet,
		Subject: "/tenants/$AccessId/services",
		Action:  "read",
	},
}
