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

import "github.com/cloustone/sentel/pkg/goshiro/shiro"

var consoleApiPolicies = []shiro.AuthorizePolicy{
	// v1api.CreateProduct
	{
		Path:     "/iot/api/v1/console/products",
		Resource: "/products",
	},
	// v1api.RemoveProduct
	{
		Path:     "/iot/api/v1/console/products/:productId",
		Resource: "/products/:productId",
	},
	// v1api.UpdateProduct
	{
		Path:     "/iot/api/v1/console/products",
		Resource: "/products",
	},
	// v1api.GetProductList
	{
		Path:     "/iot/api/v1/console/products",
		Resource: "/products",
	},
	// v1api.GetProduct
	{
		Path:     "/iot/api/v1/console/products/:productId",
		Resource: "/products/:productId",
	},
	// v1api.GetProductDevices
	{
		Path:     "/iot/api/v1/console/products/:productId/devices",
		Resource: "/products/:productId",
	},
	// v1api.GetProductRules
	{
		Path:     "/iot/api/v1/console/procuts/:productId/rules",
		Resource: "/products/:productId/rules",
	},
	// v1api.GetDeviceStatics
	{
		Path:     "/iot/api/v1/console/products/:productId/devices/statics",
		Resource: "/products/:productId",
	},
	// v1api.CreateDevice
	{
		Path:     "/iot/api/v1/console/devices",
		Resource: "/products/$ProductId",
	},
	// v1api.GetOneDevice
	{
		Path:     "/iot/api/v1/console/products/:productId/devics/:deviceId",
		Resource: "/products/:productId",
	},
	// v1api.RemoveDevice
	{
		Path:     "/iot/api/v1/console/products/:productId/devicds/:deviceId",
		Resource: "/products/:productId",
	},
	// v1api.UpdateDevice
	{
		Path:     "/iot/api/v1/console/devices",
		Resource: "/products/$ProductId/devices",
	},
	// v1api.BulkRegisterDevices
	{
		Path:     "/iot/api/v1/console/devices/bulk",
		Resource: "/products/$ProductId/devices",
	},
	// v1api.UpdateShadowDevice
	{
		Path:     "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Resource: "/products/:productId",
	},
	// v1api.GetShadowDevice
	{
		Path:     "/iot/api/v1/console/products/:productId/devices/:deviceId/shadow",
		Resource: "/products/:productId",
	},
	// v1api.CreateRule
	{
		Path:     "/iot/api/v1/console/rules",
		Resource: "/products/$ProductId",
	},
	// v1api.RemoveRule
	{
		Path:     "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Resource: "/products/:productId",
	},
	// v1api.UpdateRule
	{
		Path:     "/iot/api/v1/console/rules",
		Resource: "/products/$ProductId",
	},
	// v1api.StopRule
	{
		Path:     "/iot/api/v1/console/rules/stop",
		Resource: "/products/$ProductId",
	},
	// v1api.GetRule
	{
		Path:     "/iot/api/v1/console/products/:productId/rules/:ruleName",
		Resource: "/products/$ProductId",
	},
	// v1api.CreateTopicFlavor
	{
		Path:     "/iot/api/v1/console/topicflavors",
		Resource: "/topicflavors",
	},
	// v1api.RemoveTopicFlavor
	{
		Path:     "/iot/api/v1/console/topicflavors/products/:productId",
		Resource: "/topicflavors",
	},
	// v1api.GetProductTopicFlavor
	{
		Path:     "/iot/api/v1/console/topicflavors/products/:productId",
		Resource: "/topicflavors",
	},
	// v1api.GetTenantTopicFlavors
	{
		Path:     "/iot/api/v1/console/topicflavors/tenants/:tenantId",
		Resource: "/topicflavors",
	},
	// v1api.GetBuiltinTopicFlavors
	{
		Path:     "/iot/api/v1/console/topicflavors/builtin",
		Resource: "/topicflavors",
	},
	// v1api.SetProductTopicFlavor
	{
		Path:     "/iot/api/v1/console/topicflavors/:productId?flavor=:topicflavor",
		Resource: "/products/:productId",
	},
	// v1api.SendMessageToDevice
	{
		Path:     "/iot/api/v1/console/message",
		Resource: "/products/$ProductId",
	},
	// v1api.BroadcastProductMessage
	{
		Path:     "/iot/api/v1/console/message/broadcast",
		Resource: "/products/$ProductId",
	},
	// v1api.GetServiceStatics
	{
		Path:     "/iot/api/v1/console/service",
		Resource: "/services",
	},
}
