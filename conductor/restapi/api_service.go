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

	"github.com/cloustone/sentel/conductor/executor"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
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
	e.DELETE("conductor/api/v1/rules/:ruleName", removeRule)
	e.PUT("conductor/api/v1/rules/:ruleName?action=start", startRule)
	e.PUT("conductor/api/v1/rules/:ruleName?action=stop", stopRule)
	e.PATCH("conductor/api/v1/rules/:ruleName", updateRule)

	return &restapiService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		echo:      e,
	}, nil

}

// Name
func (p *restapiService) Name() string      { return "restapi" }
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
	r := message.RuleTopic{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// the rule should be stored to database at first
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
func updateRule(ctx echo.Context) error {
	r := message.RuleTopic{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	r.RuleAction = message.RuleActionUpdate
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}

func removeRule(ctx echo.Context) error {
	r := message.RuleTopic{RuleName: ctx.Param("ruleName"), RuleAction: message.RuleActionRemove}
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}

func startRule(ctx echo.Context) error {
	r := message.RuleTopic{RuleName: ctx.Param("ruleName"), RuleAction: message.RuleActionStart}
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
func stopRule(ctx echo.Context) error {
	r := message.RuleTopic{RuleName: ctx.Param("ruleName"), RuleAction: message.RuleActionStop}
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
