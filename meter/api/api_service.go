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

package api

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/labstack/echo"
)

type apiService struct {
	service.ServiceBase
	echo *echo.Echo
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

// apiServiceFactory
type ServiceFactory struct{}

const APIHEAD = "api/v1/"

// New create apiService service factory
func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (service.Service, error) {
	// try connect with mongo db
	hosts := c.MustString("meter", "mongo")
	timeout := c.MustInt("meter", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, fmt.Errorf("Failed to connect with mongo:'%s'", err.Error())
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

	// Clusters & Node
	e.GET(APIHEAD+"nodes", getAllNodes)
	e.GET(APIHEAD+"nodes/:nodeName", getNodeInfo)
	e.GET(APIHEAD+"nodes/clients", getNodesClientInfo)
	e.GET(APIHEAD+"nodes/:nodeName/clients", getNodeClients)
	e.GET(APIHEAD+"nodes/:nodeName/clients/:clientId", getNodeClientInfo)

	// Client
	e.GET(APIHEAD+"clients/:clientId", getClientInfo)

	// Session
	e.GET(APIHEAD+"nodes/:nodeName/sessions", getNodeSessions)
	e.GET(APIHEAD+"nodes/:nodeName/sessions/:clientId", getNodeSessionsClientInfo)
	e.GET(APIHEAD+"sessions/:clientId", getClusterSessionClientInfo)

	// Subscription
	e.GET(APIHEAD+"nodes/:nodeName/subscriptions", getNodeSubscriptions)
	e.GET(APIHEAD+"nodes/:nodeName/subscriptions/:clientId", getNodeSubscriptionsClientInfo)
	e.GET(APIHEAD+"subscriptions/:clientId", getClusterSubscriptionsInfo)

	// Routes
	e.GET(APIHEAD+"routes", getClusterRoutes)
	e.GET(APIHEAD+"routes/:topic", getTopicRoutes)

	// Publish & Subscribe
	e.POST(APIHEAD+"mqtt/publish", publishMqttMessage)
	e.POST(APIHEAD+"mqtt/subscribe", subscribeMqttMessage)
	e.POST(APIHEAD+"mqtt/unsubscribe", unsubscribeMqttMessage)

	// Plugins
	e.GET(APIHEAD+"nodes/:nodeName/plugins", getNodePluginsInfo)

	// Services
	e.GET(APIHEAD+"services", getClusterServicesInfo)
	e.GET(APIHEAD+"nodes/:nodeName/services", getNodeServicesInfo)

	// Metrics
	e.GET(APIHEAD+"metrics", getClusterMetricsInfo)
	e.GET(APIHEAD+"nodes/:nodeName/metrics", getNodeMetricsInfo)

	// Stats
	e.GET(APIHEAD+"stats", getClusterStats)
	e.GET(APIHEAD+"nodes/:nodeName/stats", getNodeStatsInfo)

	return &apiService{
		ServiceBase: service.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		echo: e,
	}, nil

}

// Name
func (p *apiService) Name() string {
	return "api"
}

// Start
func (p *apiService) Start() error {
	go func(p *apiService) {
		addr := p.Config.MustString("api", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *apiService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}
