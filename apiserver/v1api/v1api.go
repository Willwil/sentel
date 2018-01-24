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

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/cloustone/sentel/pkg/registry"
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
	RequestId string      `json:"requestId"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

func getAccessId(ctx echo.Context) string {
	return ctx.Get("AccessId").(string)
}

func getConfig(ctx echo.Context) config.Config {
	return ctx.Get("config").(config.Config)
}

func getRegistry(ctx echo.Context) *registry.Registry {
	return ctx.Get("registry").(*registry.Registry)
}

func syncProduceMessage(ctx echo.Context, topic string, value interface{}) error {
	c := getConfig(ctx)
	hosts := c.MustString("apiserver", "kafka")
	buf, err := message.Encode(value, message.JsonCodec)
	if err != nil {
		return err
	}
	return message.SendMessage(hosts, topic, buf)
}

func asyncProduceMessage(ctx echo.Context, topic string, value interface{}) error {
	c := getConfig(ctx)
	hosts := c.MustString("apiserver", "kafka")
	buf, err := message.Encode(value, message.JsonCodec)
	if err != nil {
		return err
	}
	return message.PostMessage(hosts, topic, buf)
}
