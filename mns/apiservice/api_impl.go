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

package apiservice

import (
	"net/http"

	"github.com/cloustone/sentel/mns/mns"
	"github.com/labstack/echo"
)

const (
	OK             = http.StatusOK
	ServerError    = http.StatusInternalServerError
	BadRequest     = http.StatusBadRequest
	NotFound       = http.StatusNotFound
	Unauthorized   = http.StatusUnauthorized
	NotImplemented = http.StatusNotImplemented
)

type response struct {
	RequestId string      `json:"requestId"`
	Message   string      `json:"message"`
	Result    interface{} `json:"result"`
}

func getAccount(ctx echo.Context) string {
	return "" // AccountId will be extracted from commaon request body
}

// Queue APIs
func createQueue(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if _, err := mns.CreateQueue(accountId, queueName); err != nil {
		return ctx.JSON(ServerError, err)
	}
	return ctx.JSON(OK, nil)
}

func setQueueAttribute(ctx echo.Context) error {
	attr := mns.QueueAttribute{}
	if err := ctx.Bind(&attr); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(ServerError, err)
	} else {
		queue.SetAttribute(attr)
		return ctx.JSON(OK, nil)
	}
}

func getQueueAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(ServerError, err)
	} else {
		attr := queue.GetAttribute()
		return ctx.JSON(OK, attr)
	}

}

func deleteQueue(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if err := mns.DeleteQueue(accountId, queueName); err != nil {
		return ctx.JSON(ServerError, err)
	}
	return ctx.JSON(OK, nil)
}

func getQueueList(ctx echo.Context) error {
	accountId := getAccount(ctx)
	if queues, err := mns.GetQueueList(accountId); err != nil {
		return ctx.JSON(ServerError, err)
	} else {
		return ctx.JSON(OK, queues)
	}
}

func sendQueueMessage(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
func batchSendQueueMessage(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
func receiveQueueMessage(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func batchReceiveQueueMessage(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
func deleteQueueMessage(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func batchDeleteQueueMessages(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
func peekQueueMessages(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
func batchPeekQueueMessages(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

// Topic API
func createTopic(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
func setTopicAttribute(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func getTopicAttribute(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func deleteTopic(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func listTopics(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

// Subscription API
func getSubscription(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func subscribe(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func getSubscriptionAttr(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func setSubscriptionAttr(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func unsubscribe(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func listTopicSubscriptions(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func publishMessage(ctx echo.Context) error {
	return mns.ErrNotImplemented
}

func publishNotification(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
