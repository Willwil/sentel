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
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/apiserver/v1api"
	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/golang/glog"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type consoleService struct {
	com.ServiceBase
	version string
	config  com.Config
	echo    *echo.Echo
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
	service := &consoleService{
		ServiceBase: com.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		echo: echo.New(),
	}
	if err := service.initialize(c); err != nil {
		return nil, err
	}
	return service, nil
}

func (p *consoleService) Name() string { return "console" }

// Start
func (p *consoleService) Start() error {
	go func(s *consoleService) {
		addr := p.Config.MustString("console", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *consoleService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

// Initialize initialize api manager with configuration
func (p *consoleService) initialize(c com.Config) error {
	if err := db.InitializeRegistry(c); err != nil {
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
	g.POST("/tenants", v1api.RegisterTenant)
	g.POST("/tenants/login", v1api.LoginTenant)
	g.POST("/tenants/logout", v1api.LogoutTenant)
	g.DELETE("/tenants/:tenantId", v1api.DeleteTenant)
	g.GET("/tenants/:tenantId", v1api.GetTenant)
	g.PATCH("/tenants/:tenantId", v1api.UpdateTenant)

	// Product
	g.POST("/products", v1api.CreateProduct)
	g.DELETE("/products", v1api.RemoveProduct)
	g.PATCH("/products", v1api.UpdateProduct)
	g.GET("/products", v1api.GetProductList)
	g.GET("/products/:productKey", v1api.GetProduct)
	g.GET("/products/:productKey/devices", v1api.GetProductDevices)
	g.GET("/products/:productKey/rules", v1api.GetProductRules)

	// Device
	g.POST("/devices", v1api.RegisterDevice)
	g.GET("/devices/:deviceId", v1api.GetOneDevice)
	g.DELETE("/devices", v1api.DeleteDevice)
	g.PATCH("/devices", v1api.UpdateDevice)
	g.POST("/devices/bulk", v1api.BulkRegisterDevices)

	// Rules
	g.POST("/rules", v1api.CreateRule)
	g.DELETE("/rules", v1api.RemoveRule)
	g.PATCH("/rules", v1api.UpdateRule)
	g.PUT("/rules/start", v1api.StartRule)
	g.PUT("/rules/stop", v1api.StopRule)
	g.GET("/rules/:ruleName", v1api.GetRule)

	// Runtime
	g.POST("/products/:productKey/devices/:deviceId/message", v1api.SendMessageToDevice)
	g.POST("/products/:productKey/message", v1api.BroadcastProductMessage)
	g.GET("/products/:productKey/devices/:deviceId/shardow", v1api.GetShadowDevice)
	g.PATCH("/products/:productKey/devices/:deviceId/shadow", v1api.UpdateShadowDevice)

	g.GET("/products/:productKey/devices/statics", v1api.GetDeviceStatics)
	g.GET("/service", v1api.GetServiceStatics)

	return nil
}

// setAuth setup api group 's authentication method
func (p *consoleService) setAuth(c com.Config, g *echo.Group) {
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
