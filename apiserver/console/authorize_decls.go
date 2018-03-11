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

	"github.com/cloustone/sentel/pkg/goshiro/shiro"
)

var consoleApiPolicies = []shiro.AuthorizePolicy{
	// v1api.LogoutTenant
	{
		Path:         "/iot/api/v1/console/tenants/logout",
		Methods:      http.MethodPost,
		Resource:     "/logout",
		AllowedRoles: "tenant",
	},
	// v1api.DeleteTenant
	{
		Path:         "/iot/api/v1/console/tenants/:tenantId",
		Methods:      http.MethodDelete,
		Resource:     "/tenants/$AccessId",
		AllowedRoles: "tenant",
	},
	// v1api.GetTenant
	{
		Path:         "/iot/api/v1/console/tenants/:tenantId",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId",
		AllowedRoles: "tenant",
	},
	// v1api.UpdateTenant
	{
		Path:         "/iot/api/v1/console/tenants",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId",
		AllowedRoles: "tenant",
	},

	// v1api.CreateProduct
	{
		Path:         "/iot/api/v1/console/products",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "tenant",
	},
	// v1api.RemoveProduct
	{
		Path:         "/iot/api/v1/console/products/:productId",
		Methods:      http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/:productId",
		AllowedRoles: "tenant",
	},
	// v1api.UpdateProduct
	{
		Path:         "/iot/api/v1/console/products",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "tenant",
	},
	// v1api.GetProductList
	{
		Path:         "/iot/api/v1/console/products",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products",
		AllowedRoles: "tenant",
	},
	// v1api.GetProduct
	{
		Path:         "/iot/api/v1/console/products/:productId",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId",
		AllowedRoles: "tenant",
	},
	// v1api.GetProductDevices
	{
		Path:         "/iot/api/v1/console/products/:productId/devices",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.GetProductRules
	{
		Path:         "/iot/api/v1/console/procuts/:productId/rules",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/rules",
		AllowedRoles: "tenant",
	},
	// v1api.GetDeviceStatics
	{
		Path:         "/iot/api/v1/console/products/:productId/devices/statics",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.CreateDevice
	{
		Path:         "/iot/api/v1/console/devices",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.GetOneDevice
	{
		Path:         "/iot/api/v1/console/products/:productId/devics/:deviceId",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.RemoveDevice
	{
		Path:         "/iot/api/v1/console/products/:productId/devicds/:deviceId",
		Methods:      http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.UpdateDevice
	{
		Path:         "/iot/api/v1/console/devices",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.BulkRegisterDevices
	{
		Path:         "/iot/api/v1/console/devices/bulk",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.UpdateShadowDevice
	{
		Path:         "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.GetShadowDevice
	{
		Path:         "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/:productId/devices",
		AllowedRoles: "tenant",
	},
	// v1api.CreateRule
	{
		Path:         "/iot/api/v1/console/rules",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "tenant",
	},
	// v1api.RemoveRule
	{
		Path:         "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Methods:      http.MethodDelete,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "tenant",
	},
	// v1api.UpdateRule
	{
		Path:         "/iot/api/v1/console/rules",
		Methods:      http.MethodPatch,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "tenant",
	},
	// v1api.StopRule
	{
		Path:         "/iot/api/v1/console/rules/stop",
		Methods:      http.MethodPut,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "tenant",
	},
	// v1api.GetRule
	{
		Path:         "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/products/$ProductId/rules",
		AllowedRoles: "tenant",
	},
	// v1api.CreateTopicFlavor
	{
		Path:         "/iot/api/v1/console/topicflavors",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "tenant",
	},
	// v1api.RemoveTopicFlavor
	{
		Path:         "/iot/api/v1/console/topicflavors/products/:productId",
		Methods:      http.MethodDelete,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "tenant",
	},
	// v1api.GetProductTopicFlavor
	{
		Path:         "/iot/api/v1/console/topicflavors/products/:productId",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "tenant",
	},
	// v1api.GetTenantTopicFlavors
	{
		Path:         "/iot/api/v1/console/topicflavors/tenants/:tenantId",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "tenant",
	},
	// v1api.GetBuiltinTopicFlavors
	{
		Path:         "/iot/api/v1/console/topicflavors/builtin",
		Methods:      http.MethodGet,
		Resource:     "/topicflavors/builtin",
		AllowedRoles: "tenant",
	},
	// v1api.SetProductTopicFlavor
	{
		Path:         "/iot/api/v1/console/topicflavors/:productId?flavor=:topicflavor",
		Methods:      http.MethodPut,
		Resource:     "/tenants/$AccessId/topicflavors",
		AllowedRoles: "tenant",
	},
	// v1api.SendMessageToDevice
	{
		Path:         "/iot/api/v1/console/message",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/devices/$DeviceId/message",
		AllowedRoles: "tenant",
	},
	// v1api.BroadcastProductMessage
	{
		Path:         "/iot/api/v1/console/message/broadcast",
		Methods:      http.MethodPost,
		Resource:     "/tenants/$AccessId/products/$ProductId/message",
		AllowedRoles: "tenant",
	},
	// v1api.GetServiceStatics
	{
		Path:         "/iot/api/v1/console/service",
		Methods:      http.MethodGet,
		Resource:     "/tenants/$AccessId/services",
		AllowedRoles: "tenant",
	},
}
