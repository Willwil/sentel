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
	"github.com/cloustone/sentel/keystone/client"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type managementService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	version   string
	echo      *echo.Echo
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	if err := client.Initialize(c); err != nil {
		return nil, fmt.Errorf("keystone connection failed")
	}
	service := &managementService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		echo:      echo.New(),
	}
	if err := service.initialize(c); err != nil {
		return nil, err
	}
	return service, nil
}

func (p *managementService) Name() string      { return "management" }
func (p *managementService) Initialize() error { return nil }

// Start
func (p *managementService) Start() error {
	p.waitgroup.Add(1)
	go func(s *managementService) {
		addr := p.config.MustString("management", "listen")
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

// Initialize initialize api manager with configuration
func (p *managementService) initialize(c config.Config) error {
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
	p.echo.Use(middleware.RegistryWithConfig(c))
	p.echo.Use(mw.RequestID())
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Cross-Origin
	p.echo.Use(mw.CORSWithConfig(mw.DefaultCORSConfig))

	// API for end users with restapi support
	g := p.echo.Group("/iot/api/v1")
	g.POST("/products", createProduct)
	g.PATCH("/products/:productId", updateProduct)
	g.GET("/products/:productId/devices", getProductDevices)

	g.POST("/products/:productId/device", registerDevice)
	g.POST("/prodcuts/:productId/devices/bulk", bulkApplyDevices)
	g.GET("/products/:productId/devices/bulk/:id", bulkApplyGetStatus)
	g.GET("/products/:productId/devices/bulk", bulkApplyGetDevices)
	g.GET("/products/:productId/devices/list", getDeviceList)
	g.GET("/products/:productId/devices/status", bulkGetDeviceStatus)
	g.GET("/products/:productId/devices/:deviceName", getDeviceByName)

	g.POST("/products/:productId/devices/:deviceId/props", saveDevicePropsByName)
	g.GET("/products/:productId/devices/:deviceId/props/:props", getDevicePropsByName)
	g.DELETE("/products/:productId/devices/:deviceId/props/:props", getDevicePropsByName)

	g.POST("/message", sendMessageToDevice)
	g.POST("/message/broadcast", broadcastProductMessage)
	g.GET("/products/:productId/devices/:deviceId/shardow", getShadowDevice)
	g.PATCH("/products/:productId/devices/:deviceId/shadow", updateShadowDevice)

	return nil
}
