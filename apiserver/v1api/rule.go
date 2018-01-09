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

package v1api

import (
	"time"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/labstack/echo"
)

// createRule add new rule for product
func CreateRule(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	rule := registry.Rule{}

	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId == "" || rule.RuleName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	objname := rule.ProductId + "/rules"
	if err := base.Authorize(objname, accessId, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	rule.TimeCreated = time.Now()
	rule.TimeUpdated = time.Now()
	r := getRegistry(ctx)
	if err := r.RegisterRule(&rule); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx, message.TopicNameRule,
		&message.RuleTopic{
			ProductId:  rule.ProductId,
			RuleName:   rule.RuleName,
			RuleAction: message.RuleActionCreate,
		})
	return ctx.JSON(OK, apiResponse{Result: &rule})
}

// deleteRule delete existed rule
func RemoveRule(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	rule := registry.Rule{}

	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId == "" || rule.RuleName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	objname := rule.ProductId + "/rules"
	if err := base.Authorize(objname, accessId, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	if err := r.DeleteRule(rule.ProductId, rule.RuleName); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx, message.TopicNameRule,
		&message.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: message.RuleActionRemove,
		})
	return ctx.JSON(OK, apiResponse{})
}

// UpdateRule update existed rule
func UpdateRule(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	rule := registry.Rule{}

	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId == "" || rule.RuleName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	objname := rule.ProductId + "/rules"
	if err := base.Authorize(objname, accessId, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	if err := r.UpdateRule(&rule); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx,
		message.TopicNameRule,
		&message.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: message.RuleActionUpdate,
		})
	return ctx.JSON(OK, apiResponse{Result: rule})
}

func StartRule(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	rule := registry.Rule{}

	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId == "" || rule.RuleName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	objname := rule.ProductId + "/rules"
	if err := base.Authorize(objname, accessId, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	asyncProduceMessage(ctx,
		message.TopicNameRule,
		&message.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: message.RuleActionStart,
		})
	return ctx.JSON(OK, apiResponse{})
}

func StopRule(ctx echo.Context) error {
	accessId := getAccessId(ctx)
	rule := registry.Rule{}

	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId == "" || rule.RuleName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	objname := rule.ProductId + "/rules"
	if err := base.Authorize(objname, accessId, "w"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx, message.TopicNameRule,
		&message.RuleTopic{
			RuleName:   rule.RuleName,
			ProductId:  rule.ProductId,
			RuleAction: message.RuleActionStop,
		})
	return ctx.JSON(OK, apiResponse{})
}

// getRule retrieve a rule
func GetRule(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	productId := ctx.QueryParam("productId")
	ruleName := ctx.Param("ruleName")
	if productId == "" || ruleName == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}

	objname := productId + "/rules"
	if err := base.Authorize(objname, accessId, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	rule, err := r.GetRule(productId, ruleName)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: &rule})
}

func GetProductRules(ctx echo.Context) error {
	accessId := getAccessId(ctx)

	productId := ctx.Param("productId")
	if productId == "" {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	objname := productId + "/rules"
	if err := base.Authorize(objname, accessId, "r"); err != nil {
		return ctx.JSON(Unauthorized, apiResponse{Message: err.Error()})
	}

	r := getRegistry(ctx)
	if names, err := r.GetProductRuleNames(productId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	} else {
		return ctx.JSON(OK, apiResponse{Result: names})
	}
}
