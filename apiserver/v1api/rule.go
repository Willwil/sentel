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

	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/labstack/echo"
)

// createRule add new rule for product
func CreateRule(ctx echo.Context) error {
	rule := db.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	rule.TimeCreated = time.Now()
	rule.TimeUpdated = time.Now()
	if err := r.RegisterRule(&rule); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx,
		com.TopicNameRule,
		&com.RuleTopic{
			ProductKey: rule.ProductKey,
			RuleName:   rule.RuleName,
			RuleAction: com.RuleActionCreate,
		})
	return reply(ctx, OK, apiResponse{Result: &rule})
}

// deleteRule delete existed rule
func RemoveRule(ctx echo.Context) error {
	rule := db.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	if err := r.DeleteRule(rule.ProductKey, rule.RuleName); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	// Notify kafka
	asyncProduceMessage(ctx,
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductKey: rule.ProductKey,
			RuleAction: com.RuleActionRemove,
		})
	return reply(ctx, OK, apiResponse{})
}

// UpdateRule update existed rule
func UpdateRule(ctx echo.Context) error {
	rule := db.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}

	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	if err := r.UpdateRule(&rule); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx,
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductKey: rule.ProductKey,
			RuleAction: com.RuleActionUpdate,
		})
	return reply(ctx, OK, apiResponse{Result: rule})
}

func StartRule(ctx echo.Context) error {
	rule := db.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx,
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductKey: rule.ProductKey,
			RuleAction: com.RuleActionStart,
		})
	return reply(ctx, OK, apiResponse{})
}

func StopRule(ctx echo.Context) error {
	rule := db.Rule{}
	if err := ctx.Bind(&rule); err != nil {
		return reply(ctx, BadRequest, apiResponse{Message: err.Error()})
	}
	asyncProduceMessage(ctx,
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductKey: rule.ProductKey,
			RuleAction: com.RuleActionStop,
		})
	return reply(ctx, OK, apiResponse{})
}

// getRule retrieve a rule
func GetRule(ctx echo.Context) error {
	tenantId := ctx.QueryParam("tenantId")
	productKey := ctx.QueryParam("productKey")
	ruleName := ctx.Param("ruleName")
	if tenantId == "" || productKey == "" || ruleName == "" {
		return reply(ctx, BadRequest, apiResponse{Message: "invalid parameter"})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	rule, err := r.GetRule(productKey, ruleName)
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	return reply(ctx, OK, apiResponse{Result: &rule})
}

func GetProductRules(ctx echo.Context) error {
	tenantId := ctx.QueryParam("tenantId")
	productKey := ctx.QueryParam("productKey")
	if tenantId == "" || productKey == "" {
		return reply(ctx, BadRequest, apiResponse{Message: "invalid parameter"})
	}

	// Connect with registry
	r, err := db.NewRegistry("apiserver", getConfig(ctx))
	if err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	}
	defer r.Release()
	if names, err := r.GetProductRuleNames(productKey); err != nil {
		return reply(ctx, ServerError, apiResponse{Message: err.Error()})
	} else {
		return reply(ctx, OK, apiResponse{Result: names})
	}
}
