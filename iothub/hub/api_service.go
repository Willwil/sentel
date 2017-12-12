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

	"github.com/cloustone/sentel/broker/auth"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	"github.com/labstack/echo"
)

type ApiService struct {
	core.ServiceBase
	echo *echo.Echo
}

type apiContext struct {
	echo.Context
	config core.Config
}

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// ApiServiceFactory
type ApiServiceFactory struct{}

func (p *ApiServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "iothub", "mongo")
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

	// product
	e.POST("iothub/api/v1/tenants/:tid/products/:pid", addProduct)
	e.DELETE("iothub/api/v1/tenants/:tid/products/:pid", deleteProduct)
	e.PUT("iothub/api/v1/tenants/:tid/products/:pid/action?start", startProduct)
	e.PUT("iothub/api/v1/tenants/:tid/products/:pid/action?stop", stopProduct)

	return &ApiService{
		ServiceBase: core.ServiceBase{
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

type addProductRequest struct {
	auth.Options
	Replicas int32 `json:"replicas"`
}

// addProduct add a tenant's product to iothub
func addProduct(ctx echo.Context) error {
	// Authentication
	req := &addProductRequest{}
	if err := ctx.Bind(&req); err != nil {
		glog.Error("addProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(&req.Options); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	replicas := req.Replicas

	glog.Infof("iothub: add product(%s, %s, %d)", tid, pid, replicas)
	hub := getIothub()
	if brokers, err := hub.addProduct(tid, pid, replicas); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &response{Success: true, Result: brokers})
	}
}

// deleteProduct delete a tenant from iothub
func deleteProduct(ctx echo.Context) error {
	glog.Infof("iothub: delete product(tid, %s)", ctx.Param("tid"))

	// Authentication
	req := &auth.Options{}
	if err := ctx.Bind(&req); err != nil {
		glog.Error("deleteProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	hub := getIothub()
	if err := hub.deleteProduct(tid, pid); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{Success: true})
}

func startProduct(ctx echo.Context) error {
	glog.Infof("iothub: start product(tid, %s)", ctx.Param("tid"))

	// Authentication
	req := &auth.Options{}
	if err := ctx.Bind(&req); err != nil {
		glog.Error("deleteProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	hub := getIothub()
	if err := hub.startProduct(tid, pid); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{Success: true})
}

func stopProduct(ctx echo.Context) error {
	glog.Infof("iothub: stop product(tid, %s)", ctx.Param("tid"))

	// Authentication
	req := &auth.Options{}
	if err := ctx.Bind(&req); err != nil {
		glog.Error("deleteProduct failed:%s", err.Error())
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	if err := auth.Authenticate(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}

	tid := ctx.Param("tid")
	pid := ctx.Param("pid")
	hub := getIothub()
	if err := hub.stopProduct(tid, pid); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
