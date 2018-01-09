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

package v1api

import (
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/labstack/echo"
)

// CreateTopicFlavor create topic flavor
func CreateTopicFlavor(ctx echo.Context) error {
	flavor := registry.TopicFlavor{}
	if err := ctx.Bind(&flavor); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	r := getRegistry(ctx)
	if err := r.CreateTopicFlavor(&flavor); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{})
}

func GetBuiltinTopicFlavors(ctx echo.Context) error {
	r := getRegistry(ctx)
	flavors := r.GetBuiltinTopicFlavors()
	return ctx.JSON(OK, apiResponse{Result: flavors})
}

func GetProductTopicFlavors(ctx echo.Context) error {
	productId := ctx.Param("productId")
	r := getRegistry(ctx)
	flavors := r.GetProductTopicFlavors(productId)
	return ctx.JSON(OK, apiResponse{Result: flavors})
}

func GetTenantTopicFlavors(ctx echo.Context) error {
	tenantId := ctx.Param("tenantId")
	r := getRegistry(ctx)
	flavors := r.GetTenantTopicFlavors(tenantId)
	return ctx.JSON(OK, apiResponse{Result: flavors})
}

func RemoveProductTopicFlavor(ctx echo.Context) error {
	productId := ctx.Param("productId")
	r := getRegistry(ctx)
	if p, err := r.GetProduct(productId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	} else {
		p.TopicFlavor = ""
		if err := r.UpdateProduct(p); err != nil {
			return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
		}
	}
	return ctx.JSON(OK, apiResponse{})
}

func SetProductTopicFlavor(ctx echo.Context) error {
	productId := ctx.Param("productId")
	topicFlavor := ctx.QueryParam("topicflavor")
	flavor := registry.TopicFlavor{}
	if err := ctx.Bind(&flavor); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	r := getRegistry(ctx)
	if p, err := r.GetProduct(productId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	} else {
		p.TopicFlavor = topicFlavor
		if err := r.UpdateProduct(p); err != nil {
			return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
		}
	}
	return ctx.JSON(OK, apiResponse{})
}
