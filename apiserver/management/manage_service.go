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
package management

import (
	"fmt"
	"sync"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/apiserver/middleware"
	"github.com/cloustone/sentel/apiserver/util"
	"github.com/cloustone/sentel/apiserver/v1api"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/goshiro/web"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type managementService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	version     string
	echo        *echo.Echo
	securityMgr shiro.SecurityManager
}

type ServiceFactory struct{}

var resourceMaps = make(map[string]string)

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// create resource maps
	for _, res := range mngApiPolicies {
		resourceMaps[res.Path] = res.Resource
	}

	// loading customized realm
	realm, err := base.NewAuthorizeRealm(c)
	if err != nil {
		return nil, err
	}
	securityMgr, err := goshiro.NewSecurityManager(c, realm)
	if err != nil {
		return nil, err
	}
	securityMgr.Load()
	return &managementService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		echo:        echo.New(),
		securityMgr: securityMgr,
	}, nil

}

func (p *managementService) Name() string { return "management" }
func (p *managementService) Initialize() error {
	c := p.config
	if err := registry.Initialize(c); err != nil {
		return fmt.Errorf("registry initialize failed:%v", err)
	}
	glog.Infof("Registry is initialized successfuly")

	p.echo.HideBanner = true
	p.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("SecurityManager", p.securityMgr)
			return h(ctx)
		}
	})

	// Initialize middleware
	p.echo.Use(middleware.RegistryWithConfig(c))
	if util.AuthNeed(c) {
		p.echo.Use(authenticationWithConfig(c))
		p.echo.Use(authorizeWithConfig(c))
	}
	p.echo.Use(mw.RequestID())
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Cross-Origin
	p.echo.Use(mw.CORSWithConfig(mw.DefaultCORSConfig))

	// API for end users with restapi support
	g := p.echo.Group("/iot/api/v1")
	g.POST("/products", v1api.CreateProduct)
	g.PATCH("/products/:productId", v1api.UpdateProduct)
	g.GET("/products/:productId/devices", v1api.GetProductDevices)

	g.POST("/products/:productId/devices", v1api.CreateDevice)
	g.POST("/prodcuts/:productId/devices/bulk", v1api.BulkApplyDevices)
	g.GET("/products/:productId/devices/bulk/:deviceId", v1api.BulkApplyGetStatus)
	g.GET("/products/:productId/devices/bulk", v1api.BulkApplyGetDevices)
	g.GET("/products/:productId/devices/bulk/status", v1api.BulkGetDeviceStatus)
	g.GET("/products/:productId/devices", v1api.GetDeviceList)
	g.GET("/products/:productId/devices/:deviceName", v1api.GetDeviceByName)

	g.POST("/products/:productId/devices/:deviceId/props", v1api.SaveDevicePropsByName)
	g.GET("/products/:productId/devices/:deviceId/props/:props", v1api.GetDevicePropsByName)
	g.DELETE("/products/:productId/devices/:deviceId/props/:props", v1api.RemoveDevicePropsByName)

	g.POST("/message", v1api.SendMessageToDevice)
	g.POST("/message/broadcast", v1api.BroadcastProductMessage)
	g.GET("/products/:productId/devices/:deviceId/shardow", v1api.GetShadowDevice)
	g.PATCH("/products/:productId/devices/:deviceId/shadow", v1api.UpdateShadowDevice)

	return nil

}

// Start
func (p *managementService) Start() error {
	p.waitgroup.Add(1)
	go func(s *managementService) {
		addr := p.config.MustStringWithSection("management", "listen")
		p.echo.Start(addr)
		p.waitgroup.Done()
	}(p)
	return nil
}

// Stop
func (p *managementService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}

func authenticationWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			securityManager := base.GetSecurityManager(ctx)
			authToken := web.NewRequestToken(ctx)
			if _, err := securityManager.Login(authToken); err != nil {
				return err
			}
			ctx.Set("AccessId", authToken.AccessId)
			return next(ctx)
		}
	}
}

func authorizeWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			securityManager := base.GetSecurityManager(ctx)
			token := web.NewRequestToken(ctx)
			resource, action := base.GetRequestInfo(ctx, resourceMaps)
			if err := securityManager.Authorize(token, resource, action); err != nil {
				return err
			}

			return next(ctx)
		}
	}
}
