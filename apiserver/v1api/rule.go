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
	"net/http"
	"time"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
	"github.com/cloustone/sentel/keystone/auth"
	"github.com/labstack/echo"
	uuid "github.com/satori/go.uuid"
)

// Rule Api
type ruleRequest struct {
	auth.ApiAuthParam
	rule db.Rule `json:"rule"`
}

// createRule add new rule for product
func CreateRule(ctx echo.Context) error {
	req := ruleRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule := &req.rule
	rule.TimeCreated = time.Now()
	rule.TimeUpdated = time.Now()
	if err := r.RegisterRule(rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			ProductKey: rule.ProductKey,
			RuleName:   rule.RuleName,
			RuleAction: com.RuleActionCreate,
		})
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &rule})
}

func GetProductRules(ctx echo.Context) error {
	productKey := ctx.Param("productKey")
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	if names, err := r.GetProductRuleNames(productKey); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	} else {
		return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: names})
	}
}

// deleteRule delete existed rule
func RemoveRule(ctx echo.Context) error {
	productKey := ctx.Param("productKey")
	ruleName := ctx.Param("ruleName")
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	if err := r.DeleteRule(productKey, ruleName); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	// Notify kafka
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   ruleName,
			ProductKey: productKey,
			RuleAction: com.RuleActionRemove,
		})
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String()})
}

func StartRule(ctx echo.Context) error {
	productKey := ctx.Param("productKey")
	ruleName := ctx.Param("ruleName")
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   ruleName,
			ProductKey: productKey,
			RuleAction: com.RuleActionStart,
		})
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Success: true})
}

func StopRule(ctx echo.Context) error {
	productKey := ctx.Param("productKey")
	ruleName := ctx.Param("ruleName")
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   ruleName,
			ProductKey: productKey,
			RuleAction: com.RuleActionStop,
		})
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Success: true})
}

// UpdateRule update existed rule
func UpdateRule(ctx echo.Context) error {
	req := ruleRequest{}
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule := &req.rule
	if err := r.UpdateRule(rule); err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	com.AsyncProduceMessage(ctx.(*base.ApiContext).Config,
		"rule",
		com.TopicNameRule,
		&com.RuleTopic{
			RuleName:   rule.RuleName,
			ProductKey: rule.ProductKey,
			RuleAction: com.RuleActionUpdate,
		})
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: rule})
}

// getRule retrieve a rule
func GetRule(ctx echo.Context) error {
	productKey := ctx.Param("productKey")
	ruleName := ctx.Param("ruleName")
	// Connect with registry
	r, err := db.NewRegistry("apiserver", ctx.(*base.ApiContext).Config)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	defer r.Release()
	rule, err := r.GetRule(productKey, ruleName)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &base.ApiResponse{Success: false, Message: err.Error()})
	}
	return ctx.JSON(http.StatusOK, &base.ApiResponse{RequestId: uuid.NewV4().String(), Result: &rule})
}
