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
	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/cloustone/sentel/keystone/orm"
	jwt "github.com/dgrijalva/jwt-go"
	uuid "github.com/satori/go.uuid"

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
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	// get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()
	t := db.Tenant{
		TenantId:  req.TenantId,
		Password:  req.Password,
		CreatedAt: time.Now(),
	}
	if err := r.RegisterTenant(&t); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	// notify keystone to register the object
	orm.CreateObject(&orm.Object{
		ObjectName:  req.TenantId,
		ObjectId:    req.TenantId,
		Category:    "tenant",
		CreatedTime: time.Now(),
		Owner:       req.TenantId,
	})
	// Notify kafka
	base.AsyncProduceMessage(ctx, com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: req.TenantId,
			Action:   com.ObjectActionRegister,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &t})
}

func LoginTenant(ctx echo.Context) error {
	req := registerRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}

	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	tenant, err := r.GetTenant(req.TenantId)
	if err == db.ErrorNotFound {
		return ctx.JSON(http.StatusNotFound, &base.ApiResponse{Message: err.Error()})
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
	return ctx.JSON(http.StatusOK, echo.Map{"token": t})
}

func LogoutTenant(ctx echo.Context) error {
	return nil
}

// deleteTenant delete tenant
func DeleteTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tenantId")
	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteTenant(tenantId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	// Notify kafka
	base.AsyncProduceMessage(ctx, com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: tenantId,
			Action:   com.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{})
}

// getTenant return tenant's information
func GetTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tenantIdd")
	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	if t, err := r.GetTenant(tenantId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &base.ApiResponse{Success: true, Result: &t})
	}
}

// updateTenant update tenant's information
func UpdateTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tenantId")
	t := db.Tenant{}
	if err := ctx.Bind(&t); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.UpdateTenant(tenantId, &t); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	base.AsyncProduceMessage(ctx, com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: tenantId,
			Action:   com.ObjectActionUpdate,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &t})

}

func GetTenantProductList(ctx echo.Context) error {
	return nil
}
