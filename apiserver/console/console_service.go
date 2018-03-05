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
	"errors"
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
	jwt "github.com/dgrijalva/jwt-go"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type consoleService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	version     string
	echo        *echo.Echo
	securityMgr shiro.SecurityManager
}
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	realm, err := base.NewAuthorizeRealm(c)
	if err != nil {
		return nil, err
	}
	return &consoleService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		echo:        echo.New(),
		securityMgr: goshiro.NewSecurityManager(c, consoleApiPolicies, realm),
	}, nil
}

func (p *consoleService) Name() string { return "console" }
func (p *consoleService) Initialize() error {
	c := p.config
	if err := registry.Initialize(c); err != nil {
		return fmt.Errorf("registry initialize failed:%v", err)
	}
	p.echo.HideBanner = true
	p.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &base.ApiContext{
				Context:     e,
				Config:      c,
				SecurityMgr: p.securityMgr,
			}
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
	p.echo.POST("/iot/api/v1/console/tenants", v1api.RegisterTenant)
	p.echo.POST("/iot/api/v1/console/tenants/login", v1api.LoginTenant)

	g := p.echo.Group("/iot/api/v1/console")
	p.setAuth(c, g)
	g.POST("/tenants/logout", v1api.LogoutTenant)
	g.DELETE("/tenants/:tenantId", v1api.DeleteTenant)
	g.GET("/tenants/:tenantId", v1api.GetTenant)
	g.PATCH("/tenants", v1api.UpdateTenant)

	// Product
	g.POST("/products", v1api.CreateProduct)
	g.DELETE("/products/:productId", v1api.RemoveProduct)
	g.PATCH("/products", v1api.UpdateProduct)
	g.GET("/products", v1api.GetProductList)
	g.GET("/products/:productId", v1api.GetProduct)
	g.GET("/products/:productId/devices", v1api.GetProductDevices)
	g.GET("/products/:productId/rules", v1api.GetProductRules)
	g.GET("/products/:productId/devices/statics", v1api.GetDeviceStatics)

	// Device
	g.POST("/devices", v1api.CreateDevice)
	g.GET("/products/:productId/devices/:deviceId", v1api.GetOneDevice)
	g.DELETE("/products/:productId/devices/:deviceId", v1api.RemoveDevice)
	g.PATCH("/devices", v1api.UpdateDevice)
	g.POST("/devices/bulk", v1api.BulkRegisterDevices)
	g.GET("/products/:productId/devices/:deviceId/shardow", v1api.GetShadowDevice)
	g.PATCH("/products/:productId/devices/:deviceId/shadow", v1api.UpdateShadowDevice)

	// Rules
	g.POST("/rules", v1api.CreateRule)
	g.DELETE("/products/:productId/rules/:ruleName", v1api.RemoveRule)
	g.PATCH("/rules", v1api.UpdateRule)
	g.PUT("/rules/start", v1api.StartRule)
	g.PUT("/rules/stop", v1api.StopRule)
	g.GET("/products/:productId/rules/:ruleName", v1api.GetRule)

	// Topic Flavor
	g.POST("/topicflavors", v1api.CreateTopicFlavor)
	g.DELETE("/topicflavors/products/:productId", v1api.RemoveProductTopicFlavor)
	g.GET("/topicflavors/products/:productId", v1api.GetProductTopicFlavors)
	g.GET("/topicflavors/tenants/:tenantId", v1api.GetTenantTopicFlavors)
	g.GET("topicflavors/builtin", v1api.GetBuiltinTopicFlavors)
	g.PUT("/topicflavors/:productId?flavor=:topicflavor", v1api.SetProductTopicFlavor)

	// Runtime
	g.POST("/message", v1api.SendMessageToDevice)
	g.POST("/message/broadcast", v1api.BroadcastProductMessage)

	g.GET("/service", v1api.GetServiceStatics)

	return nil

}

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

// setAuth setup api group 's authentication method
func (p *consoleService) setAuth(c config.Config, g *echo.Group) {
	auth := util.StringConfigWithDefaultValue(c, "auth", "jwt")
	switch auth {
	case "jwt":
		// Authentication config
		config := mw.JWTConfig{
			Claims:     &base.ApiJWTClaims{},
			SigningKey: []byte("secret"),
		}
		g.Use(mw.JWTWithConfig(config))
		g.Use(accessIdWithConfig(c))
		if util.AuthNeed(c) {
			g.Use(authorizeWithConfig(c))
		}
	default:
	}
}

func (p *consoleService) Authenticate(token shiro.AuthenticationToken) error {
	principal := token.GetPrincipal().(string)
	crenditals := token.GetCrenditals().(string)
	r, err := registry.New(p.config)
	if err != nil {
		return err
	}
	if t, err := r.GetTenant(principal); err == nil {
		if t.TenantId == principal && t.Password == crenditals {
			return nil
		}
	}
	return errors.New("invalid user")
}

func accessIdWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// After authenticated by gateway,the authentication paramters must bevalid
			if user, ok := ctx.Get("user").(*jwt.Token); ok {
				claims := user.Claims.(*base.ApiJWTClaims)
				ctx.Set("AccessId", claims.AccessId)
			}
			return next(ctx)
		}
	}
}

func authorizeWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			securityManager := base.GetSecurityManager(ctx)
			accessId := ctx.Get("AccessId").(string)
			token := web.JWTToken{Username: accessId}
			subject, err := securityManager.GetSubject(token)
			if err != nil {
				return errors.New("no valid subject exist")
			}
			if err := securityManager.IsPermitted(subject, web.NewRequest(ctx)); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
