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

	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/labstack/echo"
)

// createRule add new rule for product
func CreateRule(ctx echo.Context) error {
	rule := registry.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId != ctx.Param("productId") ||
		rule.RuleName != ctx.Param("ruleName") {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	rule.TimeCreated = time.Now()
	rule.TimeUpdated = time.Now()
	r := getRegistry(ctx)
	if err := r.RegisterRule(&rule); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx,
		&message.Rule{
			TopicName: message.TopicNameRule,
			ProductId: rule.ProductId,
			RuleName:  rule.RuleName,
			Action:    message.ActionCreate,
		})
	return ctx.JSON(OK, apiResponse{Result: &rule})
}

// deleteRule delete existed rule
func RemoveRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	r := getRegistry(ctx)
	if err := r.RemoveRule(productId, ruleName); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx,
		&message.Rule{
			TopicName: message.TopicNameRule,
			RuleName:  ruleName,
			ProductId: productId,
			Action:    message.ActionRemove,
		})
	return ctx.JSON(OK, apiResponse{})
}

// UpdateRule update existed rule
func UpdateRule(ctx echo.Context) error {
	rule := registry.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if rule.ProductId != ctx.Param("productId") ||
		rule.RuleName != ctx.Param("ruleName") {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}

	r := getRegistry(ctx)
	if err := r.UpdateRule(&rule); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx,
		&message.Rule{
			TopicName: message.TopicNameRule,
			RuleName:  rule.RuleName,
			ProductId: rule.ProductId,
			Action:    message.ActionUpdate,
		})
	return ctx.JSON(OK, apiResponse{Result: rule})
}

func StartRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	asyncProduceMessage(ctx,
		&message.Rule{
			TopicName: message.TopicNameRule,
			RuleName:  ruleName,
			ProductId: productId,
			Action:    message.ActionStart,
		})
	return ctx.JSON(OK, apiResponse{})
}

func StopRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	asyncProduceMessage(ctx,
		&message.Rule{
			TopicName: message.TopicNameRule,
			RuleName:  ruleName,
			ProductId: productId,
			Action:    message.ActionStop,
		})
	return ctx.JSON(OK, apiResponse{})
}

// getRule retrieve a rule
func GetRule(ctx echo.Context) error {
	productId := ctx.Param("productId")
	ruleName := ctx.Param("ruleName")
	r := getRegistry(ctx)
	rule, err := r.GetRule(productId, ruleName)
	if err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{Result: rule})
}

func GetProductRules(ctx echo.Context) error {
	productId := ctx.Param("productId")
	r := getRegistry(ctx)
	if names, err := r.GetProductRuleNames(productId); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	} else {
		return ctx.JSON(OK, apiResponse{Result: names})
	}
}
