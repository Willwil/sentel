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

type managementService struct {
	com.ServiceBase
	version string
	config  com.Config
	echo    *echo.Echo
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
	service := &managementService{
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

func (p *managementService) Name() string { return "management" }

// Start
func (p *managementService) Start() error {
	go func(s *managementService) {
		addr := p.Config.MustString("management", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *managementService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

// Initialize initialize api manager with configuration
func (p *managementService) initialize(c com.Config) error {
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

	// API for end users with restapi support
	g := p.echo.Group("/iot/api/v1")
	p.setAuth(c, g)
	g.POST("/products", v1api.CreateProduct)
	g.PATCH("/products/:productKey", v1api.UpdateProduct)
	g.GET("/products/:productKey/devices", v1api.GetProductDevices)

	g.POST("/products/:productKey/device", v1api.RegisterDevice)
	g.POST("/prodcuts/:productKey/devices/bulk", v1api.BulkApplyDevices)
	g.GET("/products/:productKey/devices/bulk/:id", v1api.BulkApplyGetStatus)
	g.GET("/products/:productKey/devices/bulk", v1api.BulkApplyGetDevices)
	g.GET("/products/:productKey/devices/list", v1api.GetDeviceList)
	g.GET("/products/:productKey/devices/status", v1api.BulkGetDeviceStatus)
	g.GET("/products/:productKey/devices/:deviceName", v1api.GetDeviceByName)

	g.POST("/products/:productKey/devices/:deviceId/props", v1api.SaveDevicePropsByName)
	g.GET("/products/:productKey/devices/:deviceId/props/:props", v1api.GetDevicePropsByName)
	g.DELETE("/products/:productKey/devices/:deviceId/props/:props", v1api.GetDevicePropsByName)

	g.POST("/products/:productKey/devices/:deviceId/message", v1api.SendMessageToDevice)
	g.POST("/products/:productKey/message", v1api.BroadcastProductMessage)
	g.GET("/products/:productKey/devices/:deviceId/shardow", v1api.GetShadowDevice)
	g.PATCH("/products/:productKey/devices/:deviceId/shadow", v1api.UpdateShadowDevice)

	return nil
}

// setAuth setup api group 's authentication method
func (p *managementService) setAuth(c com.Config, g *echo.Group) {
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
