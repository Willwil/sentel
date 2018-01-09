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
	"time"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"

	"github.com/labstack/echo"
)

type productCreateRequest struct {
	ProductName string `json:"productName"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func CreateProduct(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	req := productCreateRequest{}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}

	p := &registry.Product{
		ProductId:   ram.NewObjectId(),
		ProductName: req.ProductName,
		TimeCreated: time.Now(),
		Description: req.Description,
		Category:    req.Category,
	}

	r := getRegistry(ctx)
	if err := r.RegisterProduct(p); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify keystone
	base.CreateResource(p.ProductId, ram.ResourceCreateOption{
		Name:       req.ProductName,
		ObjectId:   p.ProductId,
		Creator:    accessId,
		Category:   "product",
		Attributes: []string{"rules", "devices"},
	})
	// Notify kafka
	asyncProduceMessage(ctx,
		message.TopicNameProduct,
		&message.ProductTopic{
			ProductId: p.ProductId,
			Action:    message.ObjectActionRegister,
		})
	return ctx.JSON(OK, apiResponse{Result: &p})
}

// removeProduct delete product from registry store
func RemoveProduct(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	productId := ctx.Param("productId")

	// Authrozie
	if err := base.Authorize(productId, accessId, "x"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	if err := r.DeleteProduct(productId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	syncProduceMessage(ctx, message.TopicNameProduct,
		&message.ProductTopic{
			ProductId: productId,
			Action:    message.ObjectActionDelete,
		})

	return ctx.JSON(OK, apiResponse{})
}

type productUpdateRequest struct {
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// updateProduct update product information in registry
func UpdateProduct(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	productId := ctx.Param("productId")

	// Authrozie
	if err := base.Authorize(productId, accessId, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	req := productUpdateRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}

	p := &registry.Product{
		ProductId:   productId,
		ProductName: req.ProductName,
		TimeUpdated: time.Now(),
		Description: req.Description,
		Category:    req.Category,
	}

	r := getRegistry(ctx)
	if err := r.UpdateProduct(p); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx, message.TopicNameProduct,
		&message.ProductTopic{
			ProductId: p.ProductId,
			Action:    message.ObjectActionUpdate,
		})

	return ctx.JSON(OK, apiResponse{})
}

func GetProductList(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	if err := base.Authorize(accessId+"/products", accessId, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	if products, err := r.GetProducts(accessId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	} else {
		type product struct {
			ProductId   string `json:"productId"`
			ProductName string `json:"productName"`
		}
		result := []product{}
		for _, p := range products {
			result = append(result, product{ProductId: p.ProductId, ProductName: p.ProductName})
		}
		return ctx.JSON(OK, apiResponse{Result: result})
	}
}

// getProduct retrieve production information from registry store
func GetProduct(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	productId := ctx.Param("productId")

	if err := base.Authorize(productId, accessId, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	p, err := r.GetProduct(productId)
	if err != nil {
		return ctx.JSON(NotFound, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: p})
}

// getProductDevices retrieve product devices list from registry store
func GetProductDevices(ctx echo.Context) error {
	accessId := ctx.QueryParam("accessId")
	productId := ctx.Param("productId")

	if err := base.Authorize(productId+"/devices", accessId, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	pdevices, err := r.GetProductDevices(productId)
	if err != nil {
		return ctx.JSON(OK, apiResponse{Message: err.Error()})
	}
	rdevices := []device{}
	for _, dev := range pdevices {
		rdevices = append(rdevices, device{DeviceId: dev.DeviceId, Status: dev.DeviceStatus})
	}
	return ctx.JSON(OK, apiResponse{Result: rdevices})
}
