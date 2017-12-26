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
package v1api

import (
	"net/http"

	"github.com/Shopify/sarama"
	"github.com/cloustone/sentel/apiserver/base"
	com "github.com/cloustone/sentel/common"
	"github.com/labstack/echo"
)

const (
	OK           = http.StatusOK
	ServerError  = http.StatusInternalServerError
	BadRequest   = http.StatusBadRequest
	NotFound     = http.StatusNotFound
	Unauthorized = http.StatusUnauthorized
)

type apiResponse struct {
	RequestId string      `json:"requestID"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

func getAccessId(ctx echo.Context) string {
	return ctx.Get("AccessId").(string)
}

func reply(ctx echo.Context, code int, r interface{}) error {
	switch r.(type) {
	case apiResponse:
		requestId := ctx.Get("RequestId").(string)
		r.(*apiResponse).RequestId = requestId
	}
	return ctx.JSON(code, r)
}

func getConfig(ctx echo.Context) com.Config {
	return ctx.(*base.ApiContext).Config
}

func syncProduceMessage(ctx echo.Context, topic string, value sarama.Encoder) error {
	c := getConfig(ctx)
	hosts := c.MustString("apiserver", "kafka")
	key := "iot"
	return com.SyncProduceMessage(hosts, key, topic, value)
}

func asyncProduceMessage(ctx echo.Context, topic string, value sarama.Encoder) error {
	c := getConfig(ctx)
	hosts := c.MustString("apiserver", "kafka")
	key := "iot"
	return com.AsyncProduceMessage(hosts, key, topic, value)
}
