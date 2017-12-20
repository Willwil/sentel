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
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/conductor/executor"
	"github.com/labstack/echo"
)

type RestapiService struct {
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

// RestapiServiceFactory
type RestapiServiceFactory struct{}

func (p *RestapiServiceFactory) New(c com.Config, quit chan os.Signal) (com.Service, error) {
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

	return &RestapiService{
		ServiceBase: com.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		echo: e,
	}, nil

}

// Name
func (p *RestapiService) Name() string { return "api" }

// Start
func (p *RestapiService) Start() error {
	go func(s *RestapiService) {
		addr := p.Config.MustString("restapi", "listen")
		p.echo.Start(addr)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *RestapiService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

func createRule(ctx echo.Context) error {
	r := com.RuleTopic{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// the rule should be stored to database at first
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
func updateRule(ctx echo.Context) error {
	r := com.RuleTopic{}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	r.RuleAction = com.RuleActionUpdate
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}

func removeRule(ctx echo.Context) error {
	r := com.RuleTopic{RuleName: ctx.Param("ruleName"), RuleAction: com.RuleActionRemove}
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}

func startRule(ctx echo.Context) error {
	r := com.RuleTopic{RuleName: ctx.Param("ruleName"), RuleAction: com.RuleActionStart}
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
func stopRule(ctx echo.Context) error {
	r := com.RuleTopic{RuleName: ctx.Param("ruleName"), RuleAction: com.RuleActionStop}
	executor.HandleRuleNotification(&r)
	return ctx.JSON(http.StatusOK, &response{Success: true})
}
