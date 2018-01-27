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

package restapi

import (
	"fmt"
	"sync"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

type apiService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	echo      *echo.Echo
}

type apiContext struct {
	echo.Context
	config config.Config
}

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

const SERVICE_NAME = "restapi"

// apiServiceFactory
type ServiceFactory struct{}

// New create apiService service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// try connect with mongo db
	hosts := c.MustString("mongo")
	timeout := c.MustInt("connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("failed to connect with mongo:'%s'", err.Error())
	}
	session.Close()

	// Create echo instance and setup router
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})
	//Cross-Origin
	e.Use(middleware.CORSWithConfig(middleware.DefaultCORSConfig))

	e.Use(middleware.RequestID())
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))

	g := e.Group("/iot/api/v1")
	// Clusters & Node
	g.GET("/nodes", getAllNodes)
	g.GET("/nodes/:nodeName", getNodeInfo)
	g.GET("/nodes/clients", getNodesClientInfo)
	g.GET("/nodes/:nodeName/clients", getNodeClients)
	g.GET("/nodes/:nodeName/clients/:clientId", getNodeClientInfo)

	// Client
	g.GET("/clients/:clientId", getClientInfo)

	// Session
	g.GET("/nodes/:nodeName/sessions", getNodeSessions)
	g.GET("/nodes/:nodeName/sessions/:clientId", getNodeSessionsClientInfo)
	g.GET("/sessions/:clientId", getClusterSessionClientInfo)

	// Subscription
	g.GET("/nodes/:nodeName/subscriptions", getNodeSubscriptions)
	g.GET("/nodes/:nodeName/subscriptions/:clientId", getNodeSubscriptionsClientInfo)
	g.GET("/subscriptions/:clientId", getClusterSubscriptionsInfo)

	// Routes
	g.GET("/routes", getClusterRoutes)
	g.GET("/routes/:topic", getTopicRoutes)

	// Publish & Subscribe
	g.POST("/mqtt/publish", publishMqttMessage)
	g.POST("/mqtt/subscribe", subscribeMqttMessage)
	g.POST("/mqtt/unsubscribe", unsubscribeMqttMessage)

	// Plugins
	g.GET("/nodes/:nodeName/plugins", getNodePluginsInfo)

	// Services
	g.GET("/services", getClusterServicesInfo)
	g.GET("/nodes/:nodeName/services", getNodeServicesInfo)

	// Metrics
	g.GET("/metrics", getClusterMetricsInfo)
	g.GET("/nodes/:nodeName/metrics", getNodeMetricsInfo)

	// Stats
	g.GET("/stats", getClusterStats)
	g.GET("/nodes/:nodeName/stats", getNodeStatsInfo)

	// Tenant
	g.POST("/tenants", createTenant)
	g.DELETE("/tenants/:tid", removeTenant)
	g.POST("/tenants/:tid/products", createProduct)
	g.DELETE("/tenants/:tid/products/:pid", removeProduct)

	return &apiService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		echo:      e,
	}, nil

}

// Name
func (p *apiService) Name() string      { return SERVICE_NAME }
func (p *apiService) Initialize() error { return nil }

// Start
func (p *apiService) Start() error {
	p.waitgroup.Add(1)
	go func(p *apiService) {
		defer p.waitgroup.Done()
		addr := p.config.MustStringWithSection("api", "listen")
		p.echo.Start(addr)
	}(p)
	return nil
}

// Stop
func (p *apiService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}
