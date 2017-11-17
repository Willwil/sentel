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
	// p.echo.Use(mw.KeyAuthWithConfig(middleware.DefaultKeyAuthConfig))
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	config := mw.JWTConfig{
		Claims:     &jwtApiClaims{},
		SigningKey: []byte("secrect"),
	}

	// Login
	p.echo.POST("/api/v1/login", login)

	// Tenant
	g := p.echo.Group("/api/v1/tenants")
	g.Use(mw.JWTWithConfig(config))
	g.POST("/api/v1/tenants/:id", addTenant)
	g.DELETE("/api/v1/tenants/:id", deleteTenant)
	g.GET("/api/v1/tenants/:id", getTenant)
	g.PUT("/api/v1/tenants/:id", updateTenant)

	// Product Api
	g = p.echo.Group("/api/v1/products")
	g.Use(mw.JWTWithConfig(config))
	g.POST("api/v1/products/:id", registerProduct)
	g.DELETE("/api/v1/products/:id", deleteProduct)
	g.GET("/api/v1/products/:id", getProduct)
	g.GET("/api/v1/products/:id/devices", getProductDevices)

	// Rule
	g = p.echo.Group("/api/v1/rules")
	g.Use(mw.JWTWithConfig(config))
	g.POST("api/v1/rules/", addRule)
	g.DELETE("/api/v1/rules/:id", deleteRule)
	g.GET("/api/v1/rules/:id", getRule)
	g.PATCH("/api/v1/rules/:id", updateRule)

	// Device Api
	g = p.echo.Group("/api/v1/devices")
	g.Use(mw.JWTWithConfig(config))
	g.POST("api/v1/devices/:id", registerDevice)
	g.GET("/devices/:id", getDevice)
	g.DELETE("api/v1/devices/:id", deleteDevice)
	g.PUT("api/v1/devices/:id", updateDevice)
	g.DELETE("api/v1/devices/:id/commands", purgeCommandQueue)
	g.GET("api/v1/devices/", getMultipleDevices)
	g.POST("api/v1/devices/query", queryDevices)
	// Http Runtip. Api
	g.POST("api/v1/devices/:id/messages/deviceBound/:etag/abandon", abandonDeviceBoundNotification)
	g.DELETE("api/v1/devices/:id/messages/devicesBound/:etag", completeDeviceBoundNotification)

	g.POST("api/v1/devices/:ideviceId/files", createFileUploadSasUri)
	g.GET("api/v1/devices/:id/message/deviceBound", receiveDeviceBoundNotification)
	g.POST("api/v1/devices/:deviceId/files/notifications", updateFileUploadStatus)
	g.POST("api/v1/devices/:id/messages/event", sendDeviceEvent)

	// Statics Api
	g = p.echo.Group("/api/v1/statics")
	g.Use(mw.JWTWithConfig(config))
	g.GET("api/v1/statistics/devices", getRegistryStatistics)
	g.GET("api/v1/statistics/service", getServiceStatistics)

	// Device Twin Api
	g = p.echo.Group("/api/v1/twins")
	g.Use(mw.JWTWithConfig(config))
	g.GET("api/v1/twins/:id", getDeviceTwin)
	g.POST("api/v1/twins/:id/methods", invokeDeviceMethod)
	g.PATCH("api/v1/twins/:id", updateDeviceTwin)

	// Job Api
	g = p.echo.Group("/api/v1/jobs")
	g.Use(mw.JWTWithConfig(config))
	g.POST("api/v1/jobs/:jobid/cancel", cancelJob)
	g.PUT("api/v1/jobs/:jobid", createJob)
	g.GET("api/v1/jobs/:jobid", getJob)
	g.GET("api/v1/jobs/query", queryJobs)

	return nil
}
