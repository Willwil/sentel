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

package auth

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/common"
	auth "github.com/cloustone/sentel/iothub/auth"
	"github.com/labstack/echo"
)

type restapiService struct {
	com.ServiceBase
	echo *echo.Echo
}

type apiContext struct {
	echo.Context
	config com.Config
}

type authResponse struct {
	Success bool `json:"success"`
}

// restapiServiceFactory
type ServiceFactory struct{}

func (p ServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
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
	e.POST("keystone/api/v1/auth/tenant", authenticateTenant)
	e.POST("keystone/api/v1/auth/device", authenticateDevice)

	// Authorization
	e.POST("keystone/api/v1/authoriize/tenant", authorizeTenant)
	e.POST("keystone/api/v1/authorize/device", authorizeTenant)

	return &restapiService{
		ServiceBase: com.ServiceBase{
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

// addTenant
type addTenantRequest struct {
	auth.Options
}

func createToken(ctx echo.Context) error {
	return nil
}

func removeToken(ctx echo.Context) error {
	return nil
}

func authenticateTenant(ctx echo.Context) error {
	r := ApiAuthParam{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &authResponse{Success: false})
	}
	err := authenticate(r)
	return ctx.JSON(http.StatusOK, &authResponse{Success: (err == nil)})

}

func authenticateDevice(ctx echo.Context) error {
	r := DeviceAuthParam{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &authResponse{Success: false})
	}
	err := authenticate(r)
	return ctx.JSON(http.StatusOK, &authResponse{Success: (err == nil)})
}

func authorizeTenant(ctx echo.Context) error {
	return nil
}

func authorizeDevice(ctx echo.Context) error {
	return nil
}
