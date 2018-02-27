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

	"github.com/cloustone/sentel/goshiro"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"

	"github.com/labstack/echo"
)

type registerRequest struct {
	TenantId string `json:"tenantId"`
	Password string `json:"password"`
}

// RegisterTenant add a new tenant
func RegisterTenant(ctx echo.Context) error {
	req := registerRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	t := registry.Tenant{
		TenantId:  req.TenantId,
		Password:  req.Password,
		CreatedAt: time.Now(),
	}
	r := getRegistry(ctx)
	if err := r.RegisterTenant(&t); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx,
		&message.Tenant{
			TopicName: message.TopicNameTenant,
			TenantId:  req.TenantId,
			Action:    message.ActionCreate,
		})

	return ctx.JSON(OK, apiResponse{Result: &t})
}

func LoginTenant(ctx echo.Context) error {
	req := registerRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	// combined with goshiro
	authToken := goshiro.JwtToken{Username: req.TenantId, Password: req.Password}
	subject, _ := goshiro.CreateSubject()
	if err := subject.Login(authToken); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}
	subject.Save()
	// Authorized
	t, _ := authToken.GetJwtToken()
	return ctx.JSON(OK, apiResponse{Result: echo.Map{"token": t}})
}

func LogoutTenant(ctx echo.Context) error {
	return ctx.JSON(OK, apiResponse{})
}

// deleteTenant delete tenant
func DeleteTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tenantId")
	if tenantId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}

	r := getRegistry(ctx)
	if err := r.DeleteTenant(tenantId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx,
		&message.Tenant{
			TopicName: message.TopicNameTenant,
			TenantId:  tenantId,
			Action:    message.ActionRemove,
		})

	return ctx.JSON(OK, apiResponse{})
}

// getTenant return tenant's information
func GetTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tenantId")
	if tenantId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	r := getRegistry(ctx)
	if t, err := r.GetTenant(tenantId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	} else {
		return ctx.JSON(OK, apiResponse{Result: &t})
	}
}

// updateTenant update tenant's information
func UpdateTenant(ctx echo.Context) error {
	t := registry.Tenant{}
	if err := ctx.Bind(&t); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	if err := r.UpdateTenant(t.TenantId, &t); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx,
		&message.Tenant{
			TopicName: message.TopicNameTenant,
			TenantId:  t.TenantId,
			Action:    message.ActionUpdate,
		})

	return ctx.JSON(OK, apiResponse{Result: &t})

}

func GetTenantProductList(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
