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

package v1

import (
	"net/http"
	"time"

	"github.com/cloustone/sentel/apiserver/db"
	"github.com/cloustone/sentel/apiserver/util"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/golang/glog"
	uuid "github.com/satori/go.uuid"

	"github.com/labstack/echo"
)

type tenantAddRequest struct {
	Id       string `json:"id"`
	Password string `json:"password"`
}

// Get a device twin
func loginTenant(ctx echo.Context) error {
	name := ctx.FormValue("username")
	pwd := ctx.FormValue("password")

	// Get registry store instance by context
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	defer r.Release()

	tenant, err := r.GetTenant(name)
	if err != nil || tenant.Password != pwd {
		return echo.ErrUnauthorized
	}

	// Authorized
	claims := &jwtApiClaims{
		Name: tenant.Id,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		},
	}

	// Creat token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Generate encoded token and send it as response
	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return err
	}
	return ctx.JSON(http.StatusOK, echo.Map{"token": t})
}

// addTenant add a new tenant
func registerTenant(ctx echo.Context) error {
	req := new(tenantAddRequest)
	if err := ctx.Bind(req); err != nil {
		glog.Error("addTenant:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Message: err.Error()})
	}
	// Get registry store instance by context
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	defer r.Release()

	t := db.Tenant{
		Id:        req.Id,
		Password:  req.Password,
		CreatedAt: time.Now(),
	}

	// Check name available
	if err := r.CheckTenantNameAvailable(req.Id); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Message: "Same tenant identifier already exist"})
	}
	if err := r.RegisterTenant(&t); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"tenant",
		util.TopicNameTenant,
		&util.TenantTopic{
			TenantId: req.Id,
			Action:   util.ObjectActionRegister,
		})

	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &t})
}

// deleteTenant delete tenant
func deleteTenant(ctx echo.Context) error {
	id := ctx.Param("id")
	// Get registry store instance by context
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteTenant(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"tenant",
		util.TopicNameTenant,
		&util.TenantTopic{
			TenantId: id,
			Action:   util.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK, &response{})
}

// getTenant return tenant's information
func getTenant(ctx echo.Context) error {
	id := ctx.Param("id")
	// Get registry store instance by context
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	defer r.Release()

	if err := r.DeleteTenant(id); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"tenant",
		util.TopicNameTenant,
		&util.TenantTopic{
			TenantId: id,
			Action:   util.ObjectActionDelete,
		})

	return ctx.JSON(http.StatusOK, &response{})
}

// updateTenant update tenant's information
func updateTenant(ctx echo.Context) error {
	req := new(tenantAddRequest)
	if err := ctx.Bind(req); err != nil {
		glog.Error("updateTenant:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Message: err.Error()})
	}
	// Get registry store instance by context
	r, err := db.NewRegistry(ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	defer r.Release()

	t := db.Tenant{
		Id:        req.Id,
		Password:  req.Password,
		UpdatedAt: time.Now(),
	}

	if err := r.UpdateTenant(req.Id, &t); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Message: err.Error()})
	}
	// Notify kafka
	util.AsyncProduceMessage(ctx.(*apiContext).config,
		"tenant",
		util.TopicNameTenant,
		&util.TenantTopic{
			TenantId: req.Id,
			Action:   util.ObjectActionUpdate,
		})

	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &t})
}
