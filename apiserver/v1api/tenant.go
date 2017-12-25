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
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo"
)

type tenantAddRequest struct {
	TenantId string `json:"tenantId"`
	Password string `json:"password"`
}

// Get a device twin
func LoginTenant(ctx echo.Context) error {
	name := ctx.FormValue("username")
	pwd := ctx.FormValue("password")

	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	tenant, err := r.GetTenant(name)
	if err == db.ErrorNotFound {
		return ctx.JSON(http.StatusNotFound, &base.ApiResponse{Message: err.Error()})
	}
	if pwd != tenant.Password {
		return echo.ErrUnauthorized
	}

	// Authorized
	claims := &base.JwtApiClaims{
		Name: tenant.TenantId,
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

// addTenant add a new tenant
func RegisterTenant(ctx echo.Context) error {
	req := new(tenantAddRequest)
	if err := ctx.Bind(req); err != nil {
		glog.Error("addTenant:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	// Get registry store instance by context
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

	// Check name available
	if err := r.CheckTenantNameAvailable(req.TenantId); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: "Same tenant identifier already exist"})
	}
	if err := r.RegisterTenant(&t); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"tenant",
		com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: req.TenantId,
			Action:   com.ObjectActionRegister,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &t})
}

// deleteTenant delete tenant
func DeleteTenant(ctx echo.Context) error {
	id := ctx.Param("id")
	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteTenant(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"tenant",
		com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: id,
			Action:   com.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{})
}

// getTenant return tenant's information
func GetTenant(ctx echo.Context) error {
	id := ctx.Param("id")
	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteTenant(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"tenant",
		com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: id,
			Action:   com.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{})
}

// updateTenant update tenant's information
func UpdateTenant(ctx echo.Context) error {
	req := new(tenantAddRequest)
	if err := ctx.Bind(req); err != nil {
		glog.Error("updateTenant:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Message: err.Error()})
	}
	// Get registry store instance by context
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	defer r.Release()

	t := db.Tenant{
		TenantId:  req.TenantId,
		Password:  req.Password,
		UpdatedAt: time.Now(),
	}

	if err := r.UpdateTenant(req.TenantId, &t); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"tenant",
		com.TopicNameTenant,
		&com.TenantTopic{
			TenantId: req.TenantId,
			Action:   com.ObjectActionUpdate,
		})

	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &t})

}

func GetTenantProductList(ctx echo.Context) error {
	return nil
}
