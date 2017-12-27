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
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	mgo "gopkg.in/mgo.v2"

	auth "github.com/cloustone/sentel/keystone/auth"
	"github.com/cloustone/sentel/keystone/orm"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/labstack/echo"
)

type restapiService struct {
	service.ServiceBase
	echo *echo.Echo
}

type apiContext struct {
	echo.Context
	config config.Config
}

type authResponse struct {
	Message string `json:"message"`
}

// restapiServiceFactory
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (service.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("keystone", "mongo")
	timeout := c.MustInt("keystone", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.Close()

	// Create echo instance and setup router
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})

	// Token
	e.POST("keystone/api/v1/token", createToken)
	e.DELETE("keystone/api/v1/token/:token", removeToken)

	// Authentication
	e.POST("keystone/api/v1/auth/api", authenticateApi)
	e.POST("keystone/api/v1/auth/device", authenticateDevice)

	// Authorization
	e.POST("keystone/api/v1/orm/object", createOrmObject)
	e.GET("keystone/api/v1/orm/:objectId", accessOrmObject)
	e.PUT("keystone/api/v1/orm/:objectId", assignOrmObjectRight)
	e.DELETE("keystone/api/v1/orm/:objectId", destroyOrmObject)

	return &restapiService{
		ServiceBase: service.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		echo: e,
	}, nil

}

// Name
func (p *restapiService) Name() string { return "restapi" }

// Start
func (p *restapiService) Start() error {
	go func(s *restapiService) {
		addr := p.Config.MustString("restapi", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *restapiService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

func createToken(ctx echo.Context) error {
	return nil
}

func removeToken(ctx echo.Context) error {
	return nil
}

func authenticateApi(ctx echo.Context) error {
	r := auth.ApiAuthParam{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &authResponse{Message: err.Error()})
	}
	if err := auth.Authenticate(r); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &authResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &authResponse{})

}

func authenticateDevice(ctx echo.Context) error {
	r := auth.DeviceAuthParam{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &authResponse{Message: err.Error()})
	}
	if err := auth.Authenticate(r); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &authResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &authResponse{})
}

func createOrmObject(ctx echo.Context) error {
	r := orm.Object{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &authResponse{Message: err.Error()})
	}
	if err := orm.CreateObject(r); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &authResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &authResponse{})
}

func accessOrmObject(ctx echo.Context) error {
	objid := ctx.Param("objectId")
	accessId := ctx.QueryParam("accessId")
	right, _ := strconv.Atoi(ctx.QueryParam("accessRight"))

	if err := orm.AccessObject(objid, accessId, orm.AccessRight(right)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &authResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &authResponse{})
}
func assignOrmObjectRight(ctx echo.Context) error {
	objid := ctx.Param("objectId")
	accessId := ctx.QueryParam("accessId")
	right, _ := strconv.Atoi(ctx.QueryParam("accessRight"))

	if err := orm.AssignObjectRight(objid, accessId, orm.AccessRight(right)); err != nil {
		return ctx.JSON(http.StatusUnauthorized, &authResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &authResponse{})

}
func destroyOrmObject(ctx echo.Context) error {
	objid := ctx.Param("objectId")
	accessId := ctx.QueryParam("accessId")
	if err := orm.DestroyObject(objid, accessId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &authResponse{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &authResponse{})
}
