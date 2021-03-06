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
	"fmt"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/message"
	"github.com/labstack/echo"
)

// Http Runtime Api

type deviceMessage struct {
	ProductId string `json:"productId"`
	DeviceId  string `json:"deviceId"`
	Topic     string `json:"topic"`
	Payload   []byte `json:"payload"`
	Qos       uint8  `json:"qos"`
	Retain    bool   `json:"retain"`
}

// Send a could-to-device message
func SendMessageToDevice(ctx echo.Context) error {
	req := deviceMessage{}

	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, apiResponse{Message: err.Error()})
	}
	if req.ProductId == "" || req.DeviceId == "" || req.Topic == "" || len(req.Payload) == 0 {
		return ctx.JSON(BadRequest, apiResponse{Message: "invalid parameter"})
	}
	c := getConfig(ctx)
	e := event.TopicPublishEvent{
		Type:    event.TopicPublish,
		Topic:   req.Topic,
		Payload: req.Payload,
		Qos:     req.Qos,
		Retain:  req.Retain,
	}
	topic := fmt.Sprintf("%s-%s-%s", req.ProductId, req.DeviceId, req.Topic)
	value, _ := event.Encode(&e, event.JSONCodec)
	msg := message.Broker{TopicName: topic, Payload: value}
	if err := message.PostMessage(c, &msg); err != nil {
		return ctx.JSON(ServerError, apiResponse{Message: err.Error()})
	}
	return ctx.JSON(OK, apiResponse{})
}

func BroadcastProductMessage(ctx echo.Context) error {
	return ctx.JSON(NotImplemented, apiResponse{})
}
