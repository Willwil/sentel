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
package v1api

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/golang/glog"

	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

type v1apiService struct {
	com.ServiceBase
	version string
	config  com.Config
	echo    *echo.Echo
}

type apiContext struct {
	echo.Context
	config com.Config
}

type response struct {
	RequestId string      `json:"requestID"`
	Success   bool        `json:"success"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

type jwtApiClaims struct {
	jwt.StandardClaims
	Name string `json:"name"`
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
	service := &v1apiService{
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

func (p *v1apiService) Name() string { return "v1api" }

// Start
func (p *v1apiService) Start() error {
	go func(s *v1apiService) {
		addr := p.Config.MustString("apiserver", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *v1apiService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

// Initialize initialize api manager with configuration
func (p *v1apiService) initialize(c com.Config) error {
	if err := db.InitializeRegistry(c); err != nil {
		return fmt.Errorf("registry initialize failed:%v", err)
	}
	glog.Infof("Registry is initialized successfuly")

	p.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})

	// Initialize middleware
	// p.echo.Use(middleware.ApiVersion(p.version))
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	//Cross-Origin
	//ie: curl -XGET -H'Origin: *' "http://localhost:4145/api/v1/devices/mytest?api-version=v1"
	p.echo.Use(mw.CORSWithConfig(mw.DefaultCORSConfig))

	g := p.echo.Group("/iot/api/v1/sys")
	g.POST("/tenants", registerTenant)
	g.POST("/tenants/:tenantiId/login", loginTenant)
	g.POST("/tenants/:tenantiId/logout", logoutTenant)
	g.DELETE("/tenants/:tenantId", deleteTenant)
	g.GET("/tenants/:tenantId", getTenant)
	g.PATCH("/tenants/:tenantId", updateTenant)
	g.POST("/tenants/:tenantId/products", registerProduct)
	g.DELETE("/products/:productKey", deleteProduct)
	g.PATCH("/products/:productKey", updateProduct)
	g.GET("/products/:productKey", getOneProduct)
	g.GET("/products/:productKey/devices", getProductDevices)
	g.GET("/tenants/:tenantId/products", getTenantProductList)
	g.POST("/prodcuts/:productKey/devices/bulk", bulkRegisterDevices)
	g.POST("/products/:productKey/device", registerDevice)
	g.GET("/products/:productKey/devices/:deviceId", getOneDevice)
	g.DELETE("/products/:productKey/devices/:deviceId", deleteOneDevice)
	g.PATCH("/products/:productKey/device/:deviceId", updateDevice)
	g.POST("/products/:productKey/rules", createRule)
	g.DELETE("/products/:prouctKey/rules/:ruleName", removeRule)
	g.GET("/products/:productKey/rules/:ruleName", getRule)
	g.PATCH("/products/:productKey/rules/:ruleName", updateRule)
	g.PUT("/products/:productKey/rules/:ruleName/start", startRule)
	g.PUT("/products/:productKey/rules/:ruleName/stop", stopRule)

	g.POST("/message/device", sendMessageToDevice)
	g.POST("/message/broadcast", broadcastMessage)
	g.GET("/products/:productKey/devices/:deviceId/shardow", getDeviceTwin)
	g.PATCH("/products/:productKey/devices/:deviceId/shadow", updateDeviceTwin)

	g.GET("/products/:productKey/devices/statics", getRegistryStatistics)
	g.GET("/service", getServiceStatistics)

	g = p.echo.Group("/iot/api/v1")
	g.POST("/tenants/:tenantId/products", registerProduct)
	g.POST("/prodcuts/:productKey/devices/bulk", bulkRegisterDevices)
	g.POST("/products/:productKey/device", registerDevice)

	g.POST("/message/device", sendMessageToDevice)
	g.POST("/message/broadcast", broadcastMessage)
	g.GET("/products/:productKey/devices/:deviceId/shardow", getDeviceTwin)
	g.PATCH("/products/:productKey/devices/:deviceId/shadow", updateDeviceTwin)

	return nil
}

// setAuth setup api group 's authentication method
func (p *v1apiService) setAuth(c com.Config, g *echo.Group) {
	auth := "jwt"
	if _, err := c.String("apiserver", "auth"); err == nil {
		auth = c.MustString("apiserver", "auth")
	}
	switch auth {
	case "jwt":
		// Authentication config
		config := mw.JWTConfig{
			Claims:     &jwtApiClaims{},
			SigningKey: []byte("secret"),
		}
		g.Use(mw.JWTWithConfig(config))
	default:
	}
}
