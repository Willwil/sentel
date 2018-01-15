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
	"net/http"
	"sync"

	"github.com/cloustone/sentel/conductor/engine"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/labstack/echo"
)

type restapiService struct {
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

// restapiServiceFactory
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// Create echo instance and setup router
	e := echo.New()
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(e echo.Context) error {
			cc := &apiContext{Context: e, config: c}
			return h(cc)
		}
	})

	// Rule
	e.POST("conductor/api/v1/rules", createRule)
	e.DELETE("conductor/api/v1/rules", removeRule)
	e.PUT("conductor/api/v1/rules?action=start", startRule)
	e.PUT("conductor/api/v1/rules?action=stop", stopRule)
	e.PATCH("conductor/api/v1/rules", updateRule)

	return &restapiService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		echo:      e,
	}, nil

}

// Name
func (p *restapiService) Name() string      { return SERVICE_NAME }
func (p *restapiService) Initialize() error { return nil }

// Start
func (p *restapiService) Start() error {
	p.waitgroup.Add(1)
	go func(s *restapiService) {
		defer s.waitgroup.Done()
		addr := s.config.MustString("restapi", "listen")
		s.echo.Start(addr)
	}(p)
	return nil
}

// Stop
func (p *restapiService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}

func createRule(ctx echo.Context) error {
	r := engine.RuleContext{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	r.Action = engine.RuleActionCreate
	engine.HandleRuleNotification(r)
	return ctx.JSON(http.StatusOK, response{})
}
func updateRule(ctx echo.Context) error {
	r := engine.RuleContext{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	r.Action = engine.RuleActionUpdate
	engine.HandleRuleNotification(r)
	return ctx.JSON(http.StatusOK, response{})
}

func removeRule(ctx echo.Context) error {
	r := engine.RuleContext{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	r.Action = engine.RuleActionRemove
	engine.HandleRuleNotification(r)
	return ctx.JSON(http.StatusOK, response{})
}

func startRule(ctx echo.Context) error {
	r := engine.RuleContext{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	r.Action = engine.RuleActionStart
	engine.HandleRuleNotification(r)
	return ctx.JSON(http.StatusOK, response{})
}
func stopRule(ctx echo.Context) error {
	r := engine.RuleContext{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	r.Action = engine.RuleActionStop
	engine.HandleRuleNotification(r)
	return ctx.JSON(http.StatusOK, response{})
}
