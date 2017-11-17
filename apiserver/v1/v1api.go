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
	"github.com/cloustone/sentel/apiserver/middleware"
	"github.com/cloustone/sentel/core"

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
	p.echo.Use(middleware.ApiVersion(p.version))
	p.echo.Use(mw.KeyAuthWithConfig(middleware.DefaultKeyAuthConfig))
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	// Initialize api routes

	// Tenant
	p.echo.POST("/api/v1/tenants/:id", addTenant)
	p.echo.DELETE("/api/v1/tenants/:id", deleteTenant)
	p.echo.GET("/api/v1/tenants/:id", getTenant)
	p.echo.PUT("/api/v1/tenants/:id", updateTenant)

	// Product Api
	p.echo.POST("api/v1/products/:id", registerProduct)
	p.echo.DELETE("/api/v1/products/:id", deleteProduct)
	p.echo.GET("/api/v1/products/:id", getProduct)
	p.echo.GET("/api/v1/products/:id/devices", getProductDevices)

	// Rule
	p.echo.POST("api/v1/rules/", addRule)
	p.echo.DELETE("/api/v1/rules/:id", deleteRule)
	p.echo.GET("/api/v1/rules/:id", getRule)
	p.echo.PATCH("/api/v1/rules/:id", updateRule)

	// Device Api
	p.echo.POST("api/v1/devices/:id", registerDevice)
	p.echo.GET("/devices/:id", getDevice)
	p.echo.DELETE("api/v1/devices/:id", deleteDevice)
	p.echo.PUT("api/v1/devices/:id", updateDevice)
	p.echo.DELETE("api/v1/devices/:id/commands", purgeCommandQueue)
	p.echo.GET("api/v1/devices/", getMultipleDevices)
	p.echo.POST("api/v1/devices/query", queryDevices)

	// Statics Api
	p.echo.GET("api/v1/statistics/devices", getRegistryStatistics)
	p.echo.GET("api/v1/statistics/service", getServiceStatistics)

	// Device Twin Api
	p.echo.GET("api/v1/twins/:id", getDeviceTwin)
	p.echo.POST("api/v1/twins/:id/methods", invokeDeviceMethod)
	p.echo.PATCH("api/v1/twins/:id", updateDeviceTwin)

	// Http Runtip. Api
	p.echo.POST("api/v1/devices/:id/messages/deviceBound/:etag/abandon", abandonDeviceBoundNotification)
	p.echo.DELETE("api/v1/devices/:id/messages/devicesBound/:etag", completeDeviceBoundNotification)

	p.echo.POST("api/v1/devices/:ideviceId/files", createFileUploadSasUri)
	p.echo.GET("api/v1/devices/:id/message/deviceBound", receiveDeviceBoundNotification)
	p.echo.POST("api/v1/devices/:deviceId/files/notifications", updateFileUploadStatus)
	p.echo.POST("api/v1/devices/:id/messages/event", sendDeviceEvent)

	// Job Api
	p.echo.POST("api/v1/jobs/:jobid/cancel", cancelJob)
	p.echo.PUT("api/v1/jobs/:jobid", createJob)
	p.echo.GET("api/v1/jobs/:jobid", getJob)
	p.echo.GET("api/v1/jobs/query", queryJobs)

	return nil
}
