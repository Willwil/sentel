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
	"errors"
	"net/http"
	"strconv"

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
	msg := mns.Message{}
	if err := ctx.Bind(&msg); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if err := queue.SendMessage(msg); err != nil {
			return ctx.JSON(ServerError, err)
		}
	}
	return ctx.JSON(OK, nil)
}

func batchSendQueueMessage(ctx echo.Context) error {
	msgs := []mns.Message{}
	if err := ctx.Bind(&msgs); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if err := queue.BatchSendMessage(msgs); err != nil {
			return ctx.JSON(ServerError, err)
		}
	}
	return ctx.JSON(OK, nil)

}
func receiveQueueMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	ws, err := strconv.Atoi(ctx.QueryParam("ws"))
	if err != nil {
		return ctx.JSON(BadRequest, err)
	}
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if msgs, err := queue.ReceiveMessage(ws); err != nil {
			return ctx.JSON(ServerError, err)
		} else {
			return ctx.JSON(OK, msgs)
		}
	}
}

func batchReceiveQueueMessage(ctx echo.Context) error {
	numberOfMessages, err1 := strconv.Atoi(ctx.QueryParam("numberOfMessages"))
	ws, err2 := strconv.Atoi(ctx.QueryParam("ws"))
	if err1 != nil || err2 != nil {
		return ctx.JSON(BadRequest, errors.New("invalid parameter"))
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if msgs, err := queue.BatchReceiveMessages(ws, numberOfMessages); err != nil {
			return ctx.JSON(ServerError, err)
		} else {
			return ctx.JSON(OK, msgs)
		}
	}
}

func deleteQueueMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	msgId := ctx.QueryParam("msgId")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if err := queue.DeleteMessage(msgId); err != nil {
			return ctx.JSON(BadRequest, err)
		}
	}
	return ctx.JSON(OK, nil)
}

func batchDeleteQueueMessages(ctx echo.Context) error {
	req := struct {
		MessageIds []string `json:"messageIds"`
	}{}
	if err := ctx.Bind(req); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		for _, msgId := range req.MessageIds {
			queue.DeleteMessage(msgId)
		}
	}
	return ctx.JSON(OK, nil)
}

func peekQueueMessages(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	ws, err1 := strconv.Atoi(ctx.QueryParam("ws"))
	queue, err2 := mns.GetQueue(accountId, queueName)

	if err1 != nil || err2 != nil {
		return ctx.JSON(BadRequest, errors.New("invalid parameter"))
	}
	if msg, err := queue.PeekMessage(ws); err != nil {
		return ctx.JSON(ServerError, err)
	} else {
		return ctx.JSON(OK, msg)
	}
}

func batchPeekQueueMessages(ctx echo.Context) error {
	numberOfMessages, err1 := strconv.Atoi(ctx.QueryParam("numberOfMessages"))
	ws, err2 := strconv.Atoi(ctx.QueryParam("ws"))
	if err1 != nil || err2 != nil {
		return ctx.JSON(BadRequest, errors.New("invalid parameter"))
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if queue, err := mns.GetQueue(accountId, queueName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if msgs, err := queue.BatchPeekMessages(ws, numberOfMessages); err != nil {
			return ctx.JSON(ServerError, err)
		} else {
			return ctx.JSON(OK, msgs)
		}
	}
}

// Topic API
func createTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	if _, err := mns.CreateTopic(accountId, topicName); err != nil {
		return ctx.JSON(ServerError, err)
	}
	return ctx.JSON(OK, nil)
}

func setTopicAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	topicAttr := mns.TopicAttribute{}
	if err := ctx.Bind(&topicAttr); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	if topic, err := mns.GetTopic(accountId, topicName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		topic.SetAttribute(topicAttr)
	}
	return ctx.JSON(OK, nil)
}

func getTopicAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	if topic, err := mns.GetTopic(accountId, topicName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		if attr, err := topic.GetAttribute(); err != nil {
			return ctx.JSON(ServerError, err)
		} else {
			return ctx.JSON(OK, attr)
		}
	}
}

func deleteTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	if err := mns.DeleteTopic(accountId, topicName); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	return ctx.JSON(OK, nil)
}

func listTopics(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topics := mns.ListTopics(accountId)
	return ctx.JSON(OK, topics)
}

// Subscription API
func subscribe(ctx echo.Context) error {
	req := struct {
		Endpoint            string `json:"endpoint" bson:"Endpoint"`
		FilterTag           string `json:"filterTag" bson:"FilterTag"`
		NotifyStrategy      string `json:"notifyStrategy" bson:"NotifyStrategy"`
		NotifyContentFormat string `json:"notifyContentFormat" bson:"NotifyContentFormat"`
	}{}
	accountId := getAccount(ctx)
	subscriptionName := ctx.Param("subscriptionName")
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	if err := mns.Subscribe(accountId, subscriptionName, req.Endpoint, req.FilterTag, req.NotifyStrategy, req.NotifyContentFormat); err != nil {
		return ctx.JSON(ServerError, err)
	} else {
		return ctx.JSON(ServerError, map[string]string{"subscriptionId": subscriptionName})
	}
}

func unsubscribe(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	subscriptionName := ctx.Param("subscriptionName")
	if err := mns.Unsubscribe(accountId, topicName, subscriptionName); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	return ctx.JSON(OK, nil)
}

func getSubscriptionAttr(ctx echo.Context) error {
	accountId := getAccount(ctx)
	subscriptionName := ctx.Param("subscriptionName")
	if subscription, err := mns.GetSubscription(accountId, subscriptionName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		attr := subscription.GetAttribute()
		return ctx.JSON(OK, attr)
	}
}

func setSubscriptionAttr(ctx echo.Context) error {
	accountId := getAccount(ctx)
	subscriptionName := ctx.Param("subscriptionName")
	attr := mns.SubscriptionAttr{}
	if err := ctx.Bind(&attr); err != nil {
		return ctx.JSON(BadRequest, err)
	}
	if subscription, err := mns.GetSubscription(accountId, subscriptionName); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		subscription.SetAttribute(attr)
		return ctx.JSON(OK, nil)
	}
}

func listTopicSubscriptions(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	pageSize, err1 := strconv.Atoi(ctx.QueryParam("pageSize"))
	pageNo, err2 := strconv.Atoi(ctx.QueryParam("pageNo"))
	startIndex, err3 := strconv.Atoi(ctx.QueryParam("startIndex"))
	if err1 != nil || err2 != nil || err3 != nil {
		return ctx.JSON(BadRequest, errors.New("invalid parameter"))
	}
	if attrs, err := mns.ListTopicSubscriptions(accountId, topicName, pageNo, pageSize, startIndex); err != nil {
		return ctx.JSON(BadRequest, err)
	} else {
		return ctx.JSON(OK, attrs)
	}
}

func publishMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	return mns.ErrNotImplemented
}

func publishNotification(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
