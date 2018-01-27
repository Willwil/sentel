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
	"github.com/cloustone/sentel/apiserver/middleware"
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
		addr := p.config.MustStringWithSection("console", "listen")
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

	p.echo.HideBanner = true
	p.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &base.ApiContext{Context: e, Config: c}
			return h(cc)
		}
	})

	// Initialize middleware
	//Cross-Origin
	p.echo.Use(mw.CORSWithConfig(mw.DefaultCORSConfig))

	p.echo.Use(mw.RequestID())
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))
	p.echo.Use(middleware.RegistryWithConfig(c))

	// Api for console
	p.echo.POST("/iot/api/v1/console/tenants", registerTenant)
	p.echo.POST("/iot/api/v1/console/tenants/login", loginTenant)

	g := p.echo.Group("/iot/api/v1/console")
	p.setAuth(c, g)
	g.POST("/tenants/logout", logoutTenant)
	g.DELETE("/tenants/:tenantId", deleteTenant)
	g.GET("/tenants/:tenantId", getTenant)
	g.PATCH("/tenants", updateTenant)

	// Product
	g.POST("/products", createProduct)
	g.DELETE("/products/:productId", removeProduct)
	g.PATCH("/products", updateProduct)
	g.GET("/products/all", getProductList)
	g.GET("/products/:productId", getProduct)
	g.GET("/products/:productId/devices", getProductDevices)
	g.GET("/products/:productId/rules", getProductRules)
	g.GET("/products/:productId/devices/statics", getDeviceStatics)

	// Device
	g.POST("/devices", createDevice)
	g.GET("/products/:productId/devices/:deviceId", getOneDevice)
	g.DELETE("/products/:productId/devices/:deviceId", removeDevice)
	g.PATCH("/devices", updateDevice)
	g.POST("/devices/bulk", bulkRegisterDevices)
	g.PATCH("/devices/shadow", updateShadowDevice)
	g.GET("/products/:productId/devices/:deviceId/shardow", getShadowDevice)

	// Rules
	g.POST("/rules", createRule)
	g.DELETE("/rules", removeRule)
	g.PATCH("/rules", updateRule)
	g.PUT("/rules/start", startRule)
	g.PUT("/rules/stop", stopRule)
	g.GET("/rules/:ruleName", getRule)

	// Topic Flavor
	g.POST("/topicflavors", createTopicFlavor)
	g.DELETE("/topicflavors/:productId", removeProductTopicFlavor)
	g.GET("/topicflavors/product/:productId", getProductTopicFlavors)
	g.GET("/topicflavors/tenant/:tenantId", getTenantTopicFlavors)
	g.GET("topicflavors/builtin", getBuiltinTopicFlavors)
	g.PUT("/topicflavors/:productId?flavor=:topicflavor", setProductTopicFlavor)

	// Runtime
	g.POST("/message", sendMessageToDevice)
	g.POST("/message/broadcast", broadcastProductMessage)

	g.GET("/service", getServiceStatics)

	return nil
}

// setAuth setup api group 's authentication method
func (p *consoleService) setAuth(c config.Config, g *echo.Group) {
	auth := "jwt"
	if _, err := c.String("auth"); err == nil {
		auth = c.MustString("auth")
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
