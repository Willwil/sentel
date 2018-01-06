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
	"github.com/cloustone/sentel/keystone/client"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	jwt "github.com/dgrijalva/jwt-go"

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
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	// get registry store instance by context
	r, err := registry.New("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	t := registry.Tenant{
		TenantId:  req.TenantId,
		Password:  req.Password,
		CreatedAt: time.Now(),
	}
	if err := r.RegisterTenant(&t); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// notify keystone to register the object
	client, _ := client.New(getConfig(ctx))
	client.CreateAccount(req.TenantId)
	// Notify kafka
	asyncProduceMessage(ctx, message.TopicNameTenant,
		&message.TenantTopic{
			TenantId: req.TenantId,
			Action:   message.ObjectActionRegister,
		})

	return reply(ctx, OK, apiResponse{Result: &t})
}

func LoginTenant(ctx echo.Context) error {
	req := registerRequest{}
	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	// Get registry store instance by context
	r, err := registry.New("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	tenant, err := r.GetTenant(req.TenantId)
	if err != nil {
		return reply(ctx, NotFound, apiResponse{Message: err.Error()})
	}
	if req.Password != tenant.Password {
		return echo.ErrUnauthorized
	}
	// Authorized
	claims := &base.JwtApiClaims{
		AccessId: tenant.TenantId,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}
	// Creat token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as base.ApiResponse
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return reply(ctx, OK, apiResponse{Result: echo.Map{"token": t}})
}

func LogoutTenant(ctx echo.Context) error {
	return nil
}

// deleteTenant delete tenant
func DeleteTenant(ctx echo.Context) error {
	req := registerRequest{}
	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	if req.TenantId == "" {
		return reply(ctx, BadRequest, apiResponse{Message: "invalid parameter"})
	}

	// Get registry store instance by context
	r, err := registry.New("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteTenant(req.TenantId); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx, message.TopicNameTenant,
		&message.TenantTopic{
			TenantId: req.TenantId,
			Action:   message.ObjectActionDelete,
		})

	return reply(ctx, OK, apiResponse{})
}

// getTenant return tenant's information
func GetTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tenantIdd")
	// Get registry store instance by context
	r, err := registry.New("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	if t, err := r.GetTenant(tenantId); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	} else {
		return reply(ctx, OK, apiResponse{Result: &t})
	}
}

// updateTenant update tenant's information
func UpdateTenant(ctx echo.Context) error {
	t := registry.Tenant{}
	if err := ctx.Bind(&t); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	// Get registry store instance by context
	r, err := registry.New("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.UpdateTenant(t.TenantId, &t); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx, message.TopicNameTenant,
		&message.TenantTopic{
			TenantId: t.TenantId,
			Action:   message.ObjectActionUpdate,
		})

	return reply(ctx, OK, apiResponse{Result: &t})

}

func GetTenantProductList(ctx echo.Context) error {
	return nil
}
