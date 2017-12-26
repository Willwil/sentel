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
	com "github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/cloustone/sentel/keystone/auth"
	"github.com/cloustone/sentel/keystone/orm"

	"github.com/labstack/echo"
)

type productCreateRequest struct {
	ProductName string `json:"productName"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

func CreateProduct(ctx echo.Context) error {
	accessId := ctx.Get("AccessId").(string)
	req := productCreateRequest{}
	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductName == "" {
		return reply(ctx, BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// inquery object access management service to authorization
	if err := orm.AccessObject(base.RootObjectProducts, accessId, orm.AccessRightFull); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	p := &db.Product{
		ProductId:   orm.NewObjectId(),
		ProductName: req.ProductName,
		TimeCreated: time.Now(),
		Description: req.Description,
		Category:    req.Category,
	}

	if err = r.RegisterProduct(p); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// Notify keystone
	orm.CreateObject(&orm.Object{
		ObjectName:     req.ProductName,
		ObjectId:       p.ProductId,
		ParentObjectId: accessId,
		Category:       "product",
		CreatedTime:    time.Now(),
		Owner:          accessId,
	})

	// Notify kafka
	asyncProduceMessage(ctx,
		com.TopicNameProduct,
		&com.ProductTopic{
			ProductId: p.ProductId,
			Action:    com.ObjectActionRegister,
		})
	return reply(ctx, OK, apiResponse{Result: &p})
}

// removeProduct delete product from registry store
func RemoveProduct(ctx echo.Context) error {
	req := auth.ApiAuthParam{}
	productId := ctx.Param("productId")
	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	if err := orm.AccessObject(productId, req.AccessId, orm.AccessRightFull); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteProduct(productId); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	syncProduceMessage(ctx, com.TopicNameProduct,
		&com.ProductTopic{
			ProductId: productId,
			Action:    com.ObjectActionDelete,
		})

	return reply(ctx, OK, apiResponse{})
}

type productUpdateRequest struct {
	auth.ApiAuthParam
	ProductId   string `json:"productId"`
	ProductName string `json:"productName"`
	Category    string `json:"category"`
	Description string `json:"description"`
}

// updateProduct update product information in registry
func UpdateProduct(ctx echo.Context) error {
	req := productUpdateRequest{}
	productId := ctx.Param("productId")

	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	if err := orm.AccessObject(req.ProductId, req.AccessId, orm.AccessRightFull); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	p := &db.Product{
		ProductId:   productId,
		ProductName: req.ProductName,
		TimeUpdated: time.Now(),
		Description: req.Description,
		Category:    req.Category,
	}

	if err := r.UpdateProduct(p); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx, com.TopicNameProduct,
		&com.ProductTopic{
			ProductId: p.ProductId,
			Action:    com.ObjectActionUpdate,
		})

	return reply(ctx, OK, apiResponse{})
}

func GetProductList(ctx echo.Context) error {
	accessId := ctx.QueryParam("accessId")
	if err := orm.AccessObject(base.RootObjectProducts, accessId, orm.AccessRightReadOnly); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	if products, err := r.GetProducts(accessId); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	} else {
		type product struct {
			ProductId   string `json:"productId"`
			ProductName string `json:"productName"`
		}
		result := []product{}
		for _, p := range products {
			result = append(result, product{ProductId: p.ProductId, ProductName: p.ProductName})
		}
		return reply(ctx, OK, apiResponse{Result: result})
	}
}

// getProduct retrieve production information from registry store
func GetProduct(ctx echo.Context) error {
	productId := ctx.Param("productId")
	accessId := ctx.QueryParam("accessId")

	if err := orm.AccessObject(productId, accessId, orm.AccessRightReadOnly); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	p, err := r.GetProduct(productId)
	if err != nil {
		return reply(ctx, NotFound, apiResponse{Message: err.Error()})
	}
	return reply(ctx, OK, apiResponse{Result: p})
}

// getProductDevices retrieve product devices list from registry store
func GetProductDevices(ctx echo.Context) error {
	accessId := ctx.QueryParam("accessId")
	productId := ctx.Param("productId")
	if err := orm.AccessObject(productId, accessId, orm.AccessRightReadOnly); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	pdevices, err := r.GetProductDevices(productId)
	if err != nil {
		return reply(ctx, OK, apiResponse{Message: err.Error()})
	}
	rdevices := []device{}
	for _, dev := range pdevices {
		rdevices = append(rdevices, device{DeviceId: dev.DeviceId, Status: dev.DeviceStatus})
	}
	return reply(ctx, OK, apiResponse{Result: rdevices})
}
