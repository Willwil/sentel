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
	"net/http"
	"time"

	"github.com/cloustone/sentel/apiserver/base"
	com "github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo"
)

func CreateProduct(ctx echo.Context) error {
	p := db.Product{}
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	if p.ProductId == "" || p.TenantId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	p.ProductKey = uuid.NewV4().String()
	p.TimeCreated = time.Now()
	if err = r.RegisterProduct(&p); err != nil {
		return ctx.JSON(http.StatusOK,
			&base.ApiResponse{RequestId: uuid.NewV4().String(), Success: false, Message: err.Error()})
	}

	// Notify kafka
	base.AsyncProduceMessage(ctx,
		com.TopicNameProduct,
		&com.ProductTopic{
			ProductKey: p.ProductKey,
			Action:     com.ObjectActionRegister,
		})
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &p})
}

// removeProduct delete product from registry store
func RemoveProduct(ctx echo.Context) error {
	p := db.Product{}
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	} else if p.TenantId == "" || p.ProductKey == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}

	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteProduct(&p); err != nil {
		return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	base.SyncProduceMessage(ctx, com.TopicNameProduct,
		&com.ProductTopic{
			ProductKey: p.ProductKey,
			Action:     com.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Success: true})
}

// updateProduct update product information in registry
func UpdateProduct(ctx echo.Context) error {
	p := db.Product{}
	if err := ctx.Bind(&p); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	} else if p.TenantId == "" || p.ProductId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()

	p.TimeModified = time.Now()
	if err := r.UpdateProduct(&p); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	// Notify kafka
	base.AsyncProduceMessage(ctx, com.TopicNameProduct,
		&com.ProductTopic{
			ProductKey: p.ProductKey,
			Action:     com.ObjectActionUpdate,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Success: true})
}

func GetProductList(ctx echo.Context) error {
	tenantId := ctx.QueryParam("tenantId")
	if tenantId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	if products, err := r.GetProducts(tenantId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	} else {
		type product struct {
			ProductId  string `json:"productId"`
			ProductKey string `json:"productKey"`
		}
		result := []product{}
		for _, p := range products {
			result = append(result, product{ProductId: p.ProductId, ProductKey: p.ProductKey})
		}
		return ctx.JSON(http.StatusOK,
			&base.ApiResponse{RequestId: uuid.NewV4().String(), Success: true, Result: result})
	}
}

// getProduct retrieve production information from registry store
func GetProduct(ctx echo.Context) error {
	productId := ctx.Param("productId")
	tenantId := ctx.QueryParam("tenantId")
	if tenantId == "" || productId == "" {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "invalid parameter"})
	}

	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()

	p, err := r.GetProduct(tenantId, productId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, &base.ApiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK,
		&base.ApiResponse{RequestId: uuid.NewV4().String(), Success: true, Result: p})
}

// getProductDevices retrieve product devices list from registry store
func GetProductDevices(ctx echo.Context) error {
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()

	pdevices, err := r.GetProductDevices(ctx.Param("productId"))
	if err != nil {
		return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	rdevices := []device{}
	for _, dev := range pdevices {
		rdevices = append(rdevices, device{DeviceId: dev.DeviceId, Status: dev.DeviceStatus})
	}
	return ctx.JSON(http.StatusOK,
		&base.ApiResponse{
			RequestId: uuid.NewV4().String(),
			Success:   true,
			Result:    rdevices,
		})
}
