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
	"sync"

	auth "github.com/cloustone/sentel/keystone/auth"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/labstack/echo"
)

type restapiService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	echo      *echo.Echo
}

type apiContext struct {
	config config.Config
	echo.Context
}

type apiResponse struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// restapiServiceFactory
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	if err := ram.Initialize(c, "direct"); err != nil {
		return nil, err
	}
	// Create echo instance and setup router
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})

	// Authentication
	e.POST("keystone/api/v1/auth/api", authenticateApi)

	// Account
	e.POST("keystone/api/v1/ram/account/:account", createAccount)
	e.DELETE("keystone/api/v1/ram/account/:account", destroyAccount)

	// Authorization
	e.POST("keystone/api/v1/ram/resource", createResource)
	e.GET("keystone/api/v1/ram/resource", accessResource)
	e.PUT("keystone/api/v1/ram/resource", addResourceGrantee)
	e.DELETE("keystone/api/v1/ram/resource", destroyResource)

	return &restapiService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		echo:      e,
	}, nil
}

// Name
func (p *restapiService) Name() string      { return "restapi" }
func (p *restapiService) Initialize() error { return nil }

// Start
func (p *restapiService) Start() error {
	p.waitgroup.Add(1)
	go func(s *restapiService) {
		defer s.waitgroup.Done()
		addr := s.config.MustString("listen")
		s.echo.Start(addr)
	}(p)
	return nil
}

// Stop
func (p *restapiService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}

func authenticateApi(ctx echo.Context) error {
	r := auth.ApiAuthParam{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &apiResponse{Message: err.Error()})
	}
	if err := auth.Authenticate(r); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &apiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &apiResponse{})

}

func createResource(ctx echo.Context) error {
	accessId := ctx.QueryParam("accessId")
	r := ram.ResourceCreateOption{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &apiResponse{Message: err.Error()})
	}
	if r, err := ram.CreateResource(accessId, &r); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &apiResponse{Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &apiResponse{Result: r})
	}
}

func accessResource(ctx echo.Context) error {
	resource := ctx.QueryParam("resource")
	accessId := ctx.QueryParam("accessId")
	action := ctx.QueryParam("action")

	if err := ram.Authorize(resource, accessId, action); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &apiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &apiResponse{})
}
func addResourceGrantee(ctx echo.Context) error {
	resource := ctx.QueryParam("resource")
	accessId := ctx.QueryParam("accessId")
	right := ctx.QueryParam("right")
	if err := ram.AddResourceGrantee(resource, accessId, right); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &apiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &apiResponse{})

}
func destroyResource(ctx echo.Context) error {
	r := ram.ResourceDestroyOption{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &apiResponse{Message: err.Error()})
	}
	if err := ram.DestroyResource(r.ObjectId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &apiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &apiResponse{})
}

func createAccount(ctx echo.Context) error {
	account := ctx.Param("account")
	if err := ram.CreateAccount(account); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &apiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &apiResponse{})
}

func destroyAccount(ctx echo.Context) error {
	account := ctx.Param("account")
	if err := ram.DestroyAccount(account); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &apiResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &apiResponse{})
}
