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

package hub

import (
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/common"
	auth "github.com/cloustone/sentel/iothub/auth"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

type ApiService struct {
	com.ServiceBase
	echo *echo.Echo
}

type apiContext struct {
	echo.Context
	config com.Config
}

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// ApiServiceFactory
type ApiServiceFactory struct{}

func (p ApiServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("iothub", "mongo")
	timeout := c.MustInt("api", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
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

	// Tenant
	e.POST("iothub/api/v1/tenants", createTenant)
	e.DELETE("iothub/api/v1/tenants/:tid", removeTenant)

	// Product
	e.POST("iothub/api/v1/tenants/:tid/products", createProduct)
	e.DELETE("iothub/api/v1/tenants/:tid/products/:pid", removeProduct)

	return &ApiService{
		ServiceBase: com.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		echo: e,
	}, nil

}

// Name
func (p *ApiService) Name() string { return "api" }

// Start
func (p *ApiService) Start() error {
	go func(s *ApiService) {
		addr := p.Config.MustString("api", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *ApiService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

// addTenant
type addTenantRequest struct {
	auth.Options
}

func createTenant(ctx echo.Context) error {
	// Authentication
	req := addTenantRequest{}
	if err := ctx.Bind(&req); err != nil {
		glog.Errorf("addTenant failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(&req.Options); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	hub := getIothub()
	if err := hub.createTenant(req.TenantId); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true})
	}
}

func removeTenant(ctx echo.Context) error {
	// Authentication
	req := auth.Options{}
	if err := ctx.Bind(&req); err != nil {
		glog.Errorf("deleteProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	hub := getIothub()
	if err := hub.removeTenant(tid); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true})
	}
}

type addProductRequest struct {
	auth.Options
	Replicas int32 `json:"replicas"`
}

func createProduct(ctx echo.Context) error {
	// Authentication
	req := &addProductRequest{}
	if err := ctx.Bind(&req); err != nil {
		glog.Errorf("addProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(&req.Options); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := req.ProductId
	replicas := req.Replicas

	if pid == "" || replicas == 0 {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: "Invalid Parameter"})
	}

	glog.Infof("iothub: add product(%s, %s, %d)", tid, pid, replicas)
	hub := getIothub()
	if brokers, err := hub.createProduct(tid, pid, replicas); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true, Result: brokers})
	}
}

func removeProduct(ctx echo.Context) error {
	glog.Infof("iothub: delete product(tid, %s)", ctx.Param("tid"))

	// Authentication
	req := &auth.Options{}
	if err := ctx.Bind(&req); err != nil {
		glog.Errorf("deleteProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	hub := getIothub()
	if err := hub.removeProduct(tid, pid); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
