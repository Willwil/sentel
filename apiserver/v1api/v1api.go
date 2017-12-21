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

const APIHEAD = "api/v1/"

// v1apiService mananage version 1 apis
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

type requestBase struct {
	Format           string `json:"format"`
	AccessKeyId      string `json:"accessKeyID"`
	Signature        string `json:"signature"`
	Timestamp        string `json:"timestamp"`
	SignatureVersion string `json:"signatureVersion"`
	SignatueNonce    string `json:"signatureNonce"`
	RegionId         string `json:"regiionID"`
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

// NewApiManager create api manager instance
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

func (p *v1apiService) Name() string { return "v1apiService" }

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

	// Login
	p.echo.POST("/api/v1/tenant", registerTenant)
	p.echo.POST("/api/v1/login", loginTenant)

	// Tenant
	g := p.echo.Group("/api/v1/tenants/")
	p.setAuth(c, g)
	g.DELETE(":id", deleteTenant)
	g.GET(":id", getTenant)
	g.PUT(":id", updateTenant)

	// Product Api
	g = p.echo.Group("/api/v1/products/")
	p.setAuth(c, g)
	g.POST(":id", registerProduct)
	g.DELETE(":id", deleteProduct)
	g.DELETE(":name/products", deleteProductByName)
	g.PUT("", updateProduct)
	g.GET(":id", getProduct)
	g.GET(":cat/products", getProductsByCat)
	g.GET(":id/devices", getProductDevices)
	g.GET("devices/page", getProductDevicesPage)
	g.GET(":name/devicesname", getProductDevicesByName)
	g.GET(":name/devicesname/page", getProductDevicesPageByName)

	// Rule
	g = p.echo.Group("/api/v1/rules/")
	p.setAuth(c, g)
	g.POST("", createRule)
	g.DELETE(":productId/:ruleName", removeRule)
	g.GET(":productId/:ruleName", getRule)
	g.PATCH(":productId/:ruleName", updateRule)

	// Device Api
	g = p.echo.Group("/api/v1/devices/")
	p.setAuth(c, g)
	g.POST(":id", registerDevice)
	g.POST(":number/bulk", bulkRegisterDevices)
	g.GET(":id", getDevice)
	g.DELETE(":id", deleteDevice)
	g.PUT(":id", updateDevice)
	g.DELETE(":id/commands", purgeCommandQueue)
	g.GET(":name/all", getMultipleDevices)
	g.POST("query", queryDevices)
	// Http Runtip. Api
	g.POST(":id/messages/deviceBound/:etag/abandon", abandonDeviceBoundNotification)
	g.DELETE(":id/messages/devicesBound/:etag", completeDeviceBoundNotification)

	g.POST(":ideviceId/files", createFileUploadSasUri)
	g.GET(":id/message/deviceBound", receiveDeviceBoundNotification)
	g.POST(":deviceId/files/notifications", updateFileUploadStatus)
	g.POST(":id/messages/event", sendDeviceEvent)

	// Statics Api
	g = p.echo.Group("/api/v1/statics/")
	p.setAuth(c, g)
	g.GET("devices", getRegistryStatistics)
	g.GET("service", getServiceStatistics)

	// Device Twin Api
	g = p.echo.Group("/api/v1/twins/")
	p.setAuth(c, g)
	g.GET(":id", getDeviceTwin)
	g.POST(":id/methods", invokeDeviceMethod)
	g.PATCH(":id", updateDeviceTwin)

	// Job Api
	g = p.echo.Group("/api/v1/jobs/")
	p.setAuth(c, g)
	g.POST(":jobid/cancel", cancelJob)
	g.PUT(":jobid", createJob)
	g.GET(":jobid", getJob)
	g.GET("query", queryJobs)

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
