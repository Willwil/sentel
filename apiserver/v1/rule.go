//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
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
	"net/http"
	"time"

	"github.com/cloustone/sentel/core"
	"github.com/cloustone/sentel/core/db"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Rule Api
type ruleAddRequest struct {
	requestBase
	db.Rule
}

// addRule add new rule for product
func createRule(ctx echo.Context) error {
	req := new(ruleAddRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	// TODO: add rule detail information
	rule := db.Rule{
		RuleName:    req.RuleName,
		ProductId:   req.ProductId,
		TimeCreated: time.Now(),
		TimeUpdated: time.Now(),
	}
	if err := r.RegisterRule(&rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	core.AsyncProduceMessage(ctx.(*apiContext).config,
		"rule",
		core.TopicNameRule,
		&core.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: core.RuleActionCreate,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &rule})
}

// deleteRule delete existed rule
func removeRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	if err := r.DeleteRule(productId, ruleName); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	core.AsyncProduceMessage(ctx.(*apiContext).config,
		"rule",
		core.TopicNameRule,
		&core.RuleTopic{
			RuleName:   ruleName,
			ProductId:  productId,
			RuleAction: core.RuleActionRemove,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String()})
}

// UpdateRule update existed rule
func updateRule(ctx echo.Context) error {
	req := new(ruleAddRequest)
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &response{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule := db.Rule{
		RuleName:    req.RuleName,
		ProductId:   req.ProductId,
		TimeUpdated: time.Now(),
	}
	if err := r.UpdateRule(&rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	// Notify kafka
	core.AsyncProduceMessage(ctx.(*apiContext).config,
		"rule",
		core.TopicNameRule,
		&core.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: core.RuleActionUpdate,
		})
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &rule})
}

// getRule retrieve a rule
func getRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*apiContext).config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule, err := r.GetRule(productId, ruleName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &response{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &response{RequestId: uuid.NewV4().String(), Result: &rule})
}
