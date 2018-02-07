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
	"encoding/json"
	"net/http"
	"sync"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/cloustone/sentel/whaler/engine"
	"github.com/labstack/echo"
)

type restapiService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	echo      *echo.Echo
}

type response struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

const SERVICE_NAME = "restapi"

// restapiServiceFactory
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// Create echo instance and setup router
	e := echo.New()

	// Rule
	e.POST("whaler/api/v1/rules", createRule)
	e.DELETE("whaler/api/v1/rules", removeRule)
	e.PUT("whaler/api/v1/rules", controlRule)
	e.PATCH("whaler/api/v1/rules", updateRule)

	// Data
	e.POST("whaler/api/v1/topic", publishTopic)

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
		addr := s.config.MustStringWithSection("restapi", "listen")
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
	r := engine.RuleContext{
		Action: engine.RuleActionCreate,
		Resp:   make(chan error),
	}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	if err := engine.HandleRuleNotification(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, response{})
}

func updateRule(ctx echo.Context) error {
	r := engine.RuleContext{
		Action: engine.RuleActionUpdate,
		Resp:   make(chan error),
	}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	if err := engine.HandleRuleNotification(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response{})
}

func removeRule(ctx echo.Context) error {
	r := engine.RuleContext{
		Action: engine.RuleActionRemove,
		Resp:   make(chan error),
	}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	if err := engine.HandleRuleNotification(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}

	return ctx.JSON(http.StatusOK, response{})
}

func controlRule(ctx echo.Context) error {
	r := engine.RuleContext{
		Resp: make(chan error),
	}
	if err := ctx.Bind(&r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	if err := engine.HandleRuleNotification(r); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, response{})
}

type topicData struct {
	ClientId  string                 `json:"clientId"`
	ProductId string                 `json:"productId"`
	Topic     string                 `json:"topic"`
	Payload   map[string]interface{} `json:"payload"`
}

// Just for test
func publishTopic(ctx echo.Context) error {
	t := topicData{}
	if err := ctx.Bind(&t); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	if t.ProductId == "" || t.Topic == "" {
		return ctx.JSON(http.StatusBadRequest, response{Message: "invalid parameter"})
	}
	payload, _ := json.Marshal(t.Payload)
	e := event.TopicPublishEvent{
		Type:     event.TopicPublish,
		ClientId: t.ClientId,
		Topic:    t.Topic,
		Payload:  payload,
	}

	value, _ := event.Encode(&e, event.JSONCodec)
	msg := &message.Broker{
		EventType: event.TopicPublish,
		Payload:   value,
	}
	if err := engine.HandleTopicNotification(t.ProductId, msg); err != nil {
		return ctx.JSON(http.StatusBadRequest, response{Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, response{})
}
