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
	"sync"
	"time"

	"github.com/cloustone/sentel/iotmanager/scheduler"
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
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

const SERVICE_NAME = "restapi"

// apiServiceFactory
type ServiceFactory struct{}

// New create apiService service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
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
		addr := p.config.MustStringWithSection(SERVICE_NAME, "listen")
		p.echo.Start(addr)
	}(p)
	return nil
}

// Stop
func (p *apiService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}

// Nodes
// getAllNodes return all nodes in clusters
func getAllNodes(ctx echo.Context) error {
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	nodes := dbc.GetAllNodes()
	return ctx.JSON(OK, response{Result: nodes})
}

// getNodeInfo return a node's detail info
func getNodeInfo(ctx echo.Context) error {
	nodeId := ctx.Param("nodeId")
	if nodeId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()

	node, err := dbc.GetNode(nodeId)
	if err != nil {
		return ctx.JSON(NotFound, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{Result: node})
}

// getNodesClientInfoWithinTimeScope return each node's client info in specified time scope
func getNodesClientInfoWithinTimeScope(ctx echo.Context) error {
	from, err1 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("from"))
	to, err2 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("to"))
	duration, err3 := time.ParseDuration(ctx.Param("unit"))
	if err1 != nil || err2 != nil || err3 != nil {
		return ctx.JSON(BadRequest, response{Message: "time format is wrong"})
	}

	if to.Sub(from) < duration {
		return ctx.JSON(BadRequest, response{Message: "time format is wrong"})
	}

	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()

	// Get all nodes, for each node, query clients collection to get client's count
	nodes := dbc.GetAllNodes()
	results := map[string][]int{}
	for _, node := range nodes {
		f := from
		result := []int{}
		for {
			t := f.Add(duration)
			if to.Sub(t) <= duration {
				break
			}
			clients := dbc.GetNodesClientWithTimeScope(node.NodeId, f, t)
			result = append(result, len(clients))
			f = f.Add(duration)
		}
		results[node.NodeId] = result
	}
	return ctx.JSON(OK, response{Result: results})

}

//getNodesClientInfo return clients static infor for each node
func getNodesClientInfo(ctx echo.Context) error {
	// Deal specifully if timescope is specified
	from := ctx.Param("from")
	if from != "" {
		return getNodesClientInfoWithinTimeScope(ctx)
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()

	// Retrun last statics for each node
	nodes := dbc.GetAllNodes()
	// For each node, query clients collection to get client's count
	result := map[string]int{}
	for _, node := range nodes {
		clients := dbc.GetNodeClients(node.NodeId)
		result[node.NodeId] = len(clients)
	}
	return ctx.JSON(OK, response{Result: result})
}

// getNodeClientsWithinTimeScope return a node's clients statics within
// timescope
func getNodeClientsWithinTimeScope(ctx echo.Context) error {
	// Check parameter's validity
	from, err1 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("from"))
	to, err2 := time.Parse("yyyy-mm-dd hh:mm:ss", ctx.Param("to"))
	duration, err3 := time.ParseDuration(ctx.Param("unit"))
	nodeId := ctx.Param("nodeId")
	if err1 != nil || err2 != nil || err3 != nil || nodeId == "" {
		return ctx.JSON(BadRequest, response{Message: "time format is wrong"})
	}

	if to.Sub(from) < duration {
		return ctx.JSON(BadRequest, response{Message: "time format is wrong"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()

	result := []int{}
	for {
		t := from.Add(duration)
		if to.Sub(t) <= duration {
			break
		}
		clients := dbc.GetNodesClientWithTimeScope(nodeId, from, to)
		result = append(result, len(clients))
	}
	return ctx.JSON(OK, response{Result: result})
}

// getNodeClients return a node's all clients
func getNodeClients(ctx echo.Context) error {
	// Deal specifully if timescope is specified
	from := ctx.Param("from")
	if from != "" {
		return getNodeClientsWithinTimeScope(ctx)
	}

	// Retrun last statics for this node
	nodeId := ctx.Param("nodeId")
	if nodeId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}

	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	clients := dbc.GetNodeClients(nodeId)
	return ctx.JSON(OK, response{Result: clients})
}

// getNodeClientInfo return spcicified client infor on a node
func getNodeClientInfo(ctx echo.Context) error {
	nodeId := ctx.Param("nodeId")
	clientId := ctx.Param("clientId")
	if nodeId == "" || clientId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	client, err := dbc.GetClientWithNode(nodeId, clientId)
	if err != nil {
		return ctx.JSON(NotFound, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{Result: client})
}

// Clients
// getClusterClientInfo return clients info in cluster
func getClientInfo(ctx echo.Context) error {
	clientId := ctx.Param("clientId")
	if clientId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	client, err := dbc.GetClient(clientId)
	if err != nil {
		return ctx.JSON(NotFound, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{Result: client})
}

// Metrics

// getClusterMetricsInfo return cluster metrics
func getClusterMetricsInfo(ctx echo.Context) error {
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	metrics := dbc.GetMetrics()
	services := map[string]map[string]uint64{}
	for _, metric := range metrics {
		if service, ok := services[metric.Service]; !ok { // not found
			services[metric.Service] = metric.Values
		} else {
			for key, val := range metric.Values {
				if _, ok := service[key]; !ok {
					service[key] = val
				} else {
					service[key] += val
				}
			}
		}
	}
	return ctx.JSON(OK, response{Result: services})
}

// getNodeMetricsInfo return a node's metrics
func getNodeMetricsInfo(ctx echo.Context) error {
	nodeId := ctx.Param("nodeId")
	if nodeId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	metric, err := dbc.GetNodeMetric(nodeId)
	if err != nil {
		return ctx.JSON(NotFound, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{Result: metric})
}

// Sessions
// getNodeSessions return a node's session
func getNodeSessions(ctx echo.Context) error {
	nodeId := ctx.Param("nodeId")
	if nodeId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}

	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	sessions := dbc.GetNodeSessions(nodeId)
	return ctx.JSON(OK, response{Result: sessions})
}

// getNodeSessionsClient return client infor in a node's sessions
func getNodeSessionsClientInfo(ctx echo.Context) error {
	nodeId := ctx.Param("nodeId")
	clientId := ctx.Param("clientId")
	if nodeId == "" || clientId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()

	session, err := dbc.GetSessionWithNode(nodeId, clientId)
	if err != nil {
		return ctx.JSON(NotFound, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{Result: session})
}

// getClusterSessionInfor return client info in cluster session
func getClusterSessionClientInfo(ctx echo.Context) error {
	clientId := ctx.Param("clientId")
	if clientId == "" {
		return ctx.JSON(BadRequest, response{Message: "Invalid parameter"})
	}
	dbc, err := openManagerDB(ctx)
	if err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	defer dbc.Close()
	session, err := dbc.GetSession(clientId)
	if err != nil {
		return ctx.JSON(NotFound, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{Result: session})
}

// getClusterStats return cluster stats
func getClusterStats(ctx echo.Context) error {
	return ctx.JSON(OK, response{})
}

//getNodeStatsInfo return a node's stats
func getNodeStatsInfo(ctx echo.Context) error {
	return ctx.JSON(OK, response{})
}

// getNodeSubscriptions return a node's subscriptions
func getNodeSubscriptions(ctx echo.Context) error {
	return ctx.JSON(OK, response{})
}

// getNodeSubscriptionsClientInfo return client info in node's subscriptions
func getNodeSubscriptionsClientInfo(ctx echo.Context) error {
	return ctx.JSON(OK, response{})
}

// getClusterSubscriptionsInfo return client info in cluster subscriptions
func getClusterSubscriptionsInfo(ctx echo.Context) error {
	return ctx.JSON(OK, response{})
}

// Testing &Tenant
// addTenant
type addTenantRequest struct {
	TenantId string `json:"tenantId"`
}

func createTenant(ctx echo.Context) error {
	req := &addTenantRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, response{Message: err.Error()})
	}
	if err := scheduler.CreateTenant(req.TenantId); err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	} else {
		return ctx.JSON(OK, response{})
	}
}

func removeTenant(ctx echo.Context) error {
	tenantId := ctx.Param("tid")
	if err := scheduler.RemoveTenant(tenantId); err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	} else {
		return ctx.JSON(OK, response{})
	}
}

type addProductRequest struct {
	ProductId string `json:"productId"`
	Replicas  int32  `json:"replicas"`
}

func createProduct(ctx echo.Context) error {
	// Authentication
	req := &addProductRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, response{Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := req.ProductId
	replicas := req.Replicas
	if pid == "" || replicas == 0 {
		return ctx.JSON(BadRequest, response{Message: "Invalid Parameter"})
	}

	if brokers, err := scheduler.CreateProduct(tid, pid, replicas); err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	} else {
		return ctx.JSON(OK, response{Result: brokers})
	}
}

func removeProduct(ctx echo.Context) error {
	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	if err := scheduler.RemoveProduct(tid, pid); err != nil {
		return ctx.JSON(ServerError, response{Message: err.Error()})
	}
	return ctx.JSON(OK, response{})
}

// Routes

// getClusterRoutes return cluster's routes table
func getClusterRoutes(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// getTopicRoutes return a topic's route
func getTopicRoutes(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// Publish & Subscribe

// publishMqttMessage will publish a mqtt message
func publishMqttMessage(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// subscribeMqttMessage subscribe a mqtt topic
func subscribeMqttMessage(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// unsubscribeMqttMessage unsubsribe mqtt topic
func unsubscribeMqttMessage(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// Plugins

// getNodePluginsInfo return plugins info for a node
func getNodePluginsInfo(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// Services

// getClusterServicesInfo return all services infor in cluster
func getClusterServicesInfo(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}

// getNodeServicesInfo return a node's service info
func getNodeServicesInfo(ctx echo.Context) error {
	return ctx.JSON(NotFound, response{})
}
