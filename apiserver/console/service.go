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
package console

import (
	"fmt"
	"sync"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/apiserver/v1api"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type consoleService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	version   string
	echo      *echo.Echo
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	service := &consoleService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		echo:      echo.New(),
	}
	if err := service.initialize(c); err != nil {
		return nil, err
	}
	return service, nil
}

func (p *consoleService) Name() string      { return "console" }
func (p *consoleService) Initialize() error { return nil }

// Start
func (p *consoleService) Start() error {
	p.waitgroup.Add(1)
	go func(s *consoleService) {
		addr := p.config.MustString("console", "listen")
		p.echo.Start(addr)
		p.waitgroup.Done()
	}(p)
	return nil
}

// Stop
func (p *consoleService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}

// Initialize initialize api manager with configuration
func (p *consoleService) initialize(c config.Config) error {
	if err := registry.Initialize(c); err != nil {
		return fmt.Errorf("registry initialize failed:%v", err)
	}
	glog.Infof("Registry is initialized successfuly")

	p.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &base.ApiContext{Context: e, Config: c}
			return h(cc)
		}
	})

	// Initialize middleware
	// p.echo.Use(middleware.ApiVersion(p.version))
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Cross-Origin
	p.echo.Use(mw.CORSWithConfig(mw.DefaultCORSConfig))

	// Api for console
	g := p.echo.Group("/iot/api/v1/console")
	p.setAuth(c, g)
	g.POST("/tenants", registerTenant)
	g.POST("/tenants/login", loginTenant)
	g.POST("/tenants/logout", logoutTenant)
	g.DELETE("/tenants", deleteTenant)
	g.GET("/tenants/:tenantId", getTenant)
	g.PATCH("/tenants", updateTenant)

	// Product
	g.POST("/products", createProduct)
	g.DELETE("/products/:productId", removeProduct)
	g.PATCH("/products/:productId", updateProduct)
	g.GET("/products/", getProductList)
	g.GET("/products/:productId", getProduct)
	g.GET("/products/:productId/devices", getProductDevices)
	g.GET("/products/:productId/rules", getProductRules)
	g.GET("/products/:productId/devices/statics", getDeviceStatics)

	// Device
	g.POST("/devices", v1api.RegisterDevice)
	g.GET("/devices/:deviceId", v1api.GetOneDevice)
	g.DELETE("/devices", v1api.DeleteDevice)
	g.PATCH("/devices", v1api.UpdateDevice)
	g.POST("/devices/bulk", v1api.BulkRegisterDevices)
	g.PATCH("/devices/:deviceId/shadow", v1api.UpdateShadowDevice)
	g.GET("/devices/:deviceId/shardow", v1api.GetShadowDevice)

	// Rules
	g.POST("/rules", createRule)
	g.DELETE("/rules", removeRule)
	g.PATCH("/rules", updateRule)
	g.PUT("/rules/start", startRule)
	g.PUT("/rules/stop", stopRule)
	g.GET("/rules/:ruleName", getRule)

	// Runtime
	g.POST("/message", sendMessageToDevice)
	g.POST("/message/broadcast", broadcastProductMessage)

	g.GET("/service", getServiceStatics)

	return nil
}

// setAuth setup api group 's authentication method
func (p *consoleService) setAuth(c config.Config, g *echo.Group) {
	auth := "jwt"
	if _, err := c.String("apiserver", "auth"); err == nil {
		auth = c.MustString("apiserver", "auth")
	}
	switch auth {
	case "jwt":
		// Authentication config
		config := mw.JWTConfig{
			Claims:     &base.JwtApiClaims{},
			SigningKey: []byte("secret"),
		}
		g.Use(mw.JWTWithConfig(config))
	default:
	}
}
