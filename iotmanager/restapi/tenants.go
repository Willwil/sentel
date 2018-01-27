//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package restapi

import (
	"net/http"

	"github.com/cloustone/sentel/iotmanager/scheduler"
	"github.com/labstack/echo"
)

// addTenant
type addTenantRequest struct {
	TenantId string `json:"tenantId"`
}

func createTenant(ctx echo.Context) error {
	req := &addTenantRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := scheduler.CreateTenant(req.TenantId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true})
	}
}

func removeTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tid")
	if err := scheduler.RemoveTenant(tenantId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true})
	}
}

type addProductRequest struct {
	ProductId string `json:"productId"`
	Replicas  int32  `json:"replicas"`
}

func createProduct(ctx echo.Context) error {
	// Authentication
	req := &addProductRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := req.ProductId
	replicas := req.Replicas
	if pid == "" || replicas == 0 {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: "Invalid Parameter"})
	}

	if brokers, err := scheduler.CreateProduct(tid, pid, replicas); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true, Result: brokers})
	}
}

func removeProduct(ctx echo.Context) error {
	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	if err := scheduler.RemoveProduct(tid, pid); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
