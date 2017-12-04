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
package v1

import (
	"github.com/cloustone/sentel/apiserver"
	"github.com/cloustone/sentel/core"

	"github.com/dgrijalva/jwt-go"
	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

const APIHEAD = "api/v1/"

// v1apiManager mananage version 1 apis
type v1apiManager struct {
	version string
	config  core.Config
	echo    *echo.Echo
}

type apiContext struct {
	echo.Context
	config core.Config
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

// NewApiManager create api manager instance
func NewApiManager() apiserver.ApiManager {
	return &v1apiManager{
		echo:    echo.New(),
		version: "v1",
	}
}

// GetVersion return api's version
func (p *v1apiManager) GetVersion() string { return p.version }

// Run loop to wait api server to terminate
func (p *v1apiManager) Run() error {
	address := p.config.MustString("apiserver", "listen")
	return p.echo.Start(address)
}

// Initialize initialize api manager with configuration
func (p *v1apiManager) Initialize(c core.Config) error {
	p.config = c
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
	g.GET(":id", getProduct)
	g.GET(":id/devices", getProductDevices)
	g.GET(":name", getProductDevicesByName)

	// Rule
	g = p.echo.Group("/api/v1/rules/")
	p.setAuth(c, g)
	g.POST("", addRule)
	g.DELETE(":id", deleteRule)
	g.GET(":id", getRule)
	g.PATCH(":id", updateRule)

	// Device Api
	g = p.echo.Group("/api/v1/devices/")
	p.setAuth(c, g)
	g.POST(":id", registerDevice)
	g.GET(":id", getDevice)
	g.DELETE(":id", deleteDevice)
	g.PUT(":id", updateDevice)
	g.DELETE(":id/commands", purgeCommandQueue)
	g.GET("", getMultipleDevices)
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
func (p *v1apiManager) setAuth(c core.Config, g *echo.Group) {
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
