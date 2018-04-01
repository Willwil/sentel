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
	"strconv"

	"github.com/cloustone/sentel/mns/mns"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
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

type Base64Bytes []byte

type ErrorMessageResponse struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}
type MessageResponse struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}

type MessageSendRequest struct {
	MessageBody  Base64Bytes `json:"message_body"`
	DelaySeconds int64       `json:"delay_seconds"`
	Priority     int64       `json:"priority"`
}

type BatchMessageSendRequest struct {
	Messages []MessageSendRequest `json:"message"`
}

type ReceiptHandles struct {
	ReceiptHandles []string `json:"receipt_handles"`
}

type MessageSendResponse struct {
	MessageResponse
	MessageId      string `json:"message_id"`
	MessageBodyMD5 string `json:"message_body_md5"`
}

type BatchMessageSendResponse struct {
	Messages []MessageSendResponse `json:"messages"`
}

type CreateQueueRequest struct {
	DelaySeconds           int32 `json:"delay_senconds,omitempty"`
	MaxMessageSize         int32 `json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32 `json:"message_retention_period,omitempty"`
	VisibilityTimeout      int32 `json:"visibility_timeout,omitempty"`
	PollingWaitSeconds     int32 `json:"polling_wait_secods,omitempty"`
}

type MessageReceiveResponse struct {
	MessageResponse
	ReceiptHandle    string      `json:"receipt_handle"`
	MessageBodyMD5   string      `json:"message_body_md5"`
	MessageBody      Base64Bytes `json:"message_body"`
	EnqueueTime      int64       `json:"enqueue_time"`
	NextVisibleTime  int64       `json:"next_visible_time"`
	FirstDequeueTime int64       `json:"first_dequeue_time"`
	DequeueCount     int64       `json:"dequeue_count"`
	Priority         int64       `json:"priority"`
}

type BatchMessageReceiveResponse struct {
	Messages []MessageReceiveResponse `json:"messages"`
}

type MessageVisibilityChangeResponse struct {
	ReceiptHandle   string `json:"receipt_handle"`
	NextVisibleTime int64  `json:"next_visible_time"`
}

func getAccount(ctx echo.Context) string {
	principal := ctx.Get("Principal").(shiro.Principal)
	return principal.Name()
}

func reply(ctx echo.Context, val ...interface{}) error {
	if len(val) > 0 {
		switch val[0].(type) {
		case *mns.MnsError:
			err := val[0].(*mns.MnsError)
			resp := ErrorMessageResponse{
				Code:      err.Message,
				RequestId: ctx.Request().Header.Get(echo.HeaderXRequestID),
			}
			if len(val) > 1 {
				err := val[1].(error)
				resp.Message = err.Error()
			}
			return ctx.JSON(err.StatusCode, resp)
		default:
			return ctx.JSON(http.StatusOK, val)
		}
	}
	return ctx.JSON(http.StatusOK, nil)
}

func replyObject(ctx echo.Context, obj interface{}, err error) error {
	if err == nil {
		return reply(ctx, obj)
	} else {
		return reply(ctx, err)
	}
}

// Queue APIs
func createQueue(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	if _, err := getManager().CreateQueue(accountId, queueName); err != nil {
		return reply(ctx, err)
	}
	return reply(ctx)
}

func setQueueAttribute(ctx echo.Context) error {
	attr := mns.QueueAttribute{}
	if err := ctx.Bind(&attr); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	err := getManager().SetQueueAttribute(accountId, queueName, attr)
	return reply(ctx, err)
}

func getQueueAttribute(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	attr, err := getManager().GetQueueAttribute(accountId, queueName)
	return replyObject(ctx, attr, err)
}

func deleteQueue(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	err := getManager().DeleteQueue(accountId, queueName)
	return reply(ctx, err)
}

func getQueueList(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queues := getManager().GetQueues(accountId)
	return replyObject(ctx, queues, nil)
}

func sendQueueMessage(ctx echo.Context) error {
	msg := mns.QueueMessage{}
	if err := ctx.Bind(&msg); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	err := getManager().SendQueueMessage(accountId, queueName, msg)
	return reply(ctx, err)
}

func batchSendQueueMessage(ctx echo.Context) error {
	msgs := []mns.QueueMessage{}
	if err := ctx.Bind(&msgs); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	err := getManager().BatchSendQueueMessage(accountId, queueName, msgs)
	return reply(ctx, err)

}
func receiveQueueMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	ws, err := strconv.Atoi(ctx.QueryParam("ws"))
	if err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	msgs, err := getManager().ReceiveQueueMessage(accountId, queueName, ws)
	return replyObject(ctx, msgs, err)
}

func batchReceiveQueueMessage(ctx echo.Context) error {
	numberOfMessages, err1 := strconv.Atoi(ctx.QueryParam("numberOfMessages"))
	ws, err2 := strconv.Atoi(ctx.QueryParam("ws"))
	if err1 != nil || err2 != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	msgs, err := getManager().BatchReceiveQueueMessages(accountId, queueName, ws, numberOfMessages)
	return replyObject(ctx, msgs, err)
}

func deleteQueueMessage(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	msgId := ctx.QueryParam("msgId")
	err := getManager().DeleteQueueMessage(accountId, queueName, msgId)
	return reply(ctx, err)
}

func batchDeleteQueueMessages(ctx echo.Context) error {
	req := struct {
		MessageIds []string `json:"messageIds"`
	}{}
	if err := ctx.Bind(req); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	err := getManager().BatchDeleteQueueMessages(accountId, queueName, req.MessageIds)
	return reply(ctx, err)
}

func peekQueueMessages(ctx echo.Context) error {
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	ws, err := strconv.Atoi(ctx.QueryParam("ws"))
	if err != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	msg, err := getManager().PeekQueueMessage(accountId, queueName, ws)
	return reply(ctx, msg, err)
}

func batchPeekQueueMessages(ctx echo.Context) error {
	numberOfMessages, err1 := strconv.Atoi(ctx.QueryParam("numberOfMessages"))
	ws, err2 := strconv.Atoi(ctx.QueryParam("ws"))
	if err1 != nil || err2 != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	accountId := getAccount(ctx)
	queueName := ctx.Param("queueName")
	msgs, err := getManager().BatchPeekQueueMessages(accountId, queueName, ws, numberOfMessages)
	return replyObject(ctx, msgs, err)
}

// Topic API
func createTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	attr, err := getManager().CreateTopic(accountId, topicName)
	return replyObject(ctx, attr, err)
}

func setTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	topicAttr := mns.Topic{}
	if err := ctx.Bind(&topicAttr); err != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	err := getManager().SetTopic(accountId, topicName, topicAttr)
	return reply(ctx, err)
}

func getTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	attr, err := getManager().GetTopic(accountId, topicName)
	return reply(ctx, attr, err)
}

func deleteTopic(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	err := getManager().DeleteTopic(accountId, topicName)
	return reply(ctx, err)
}

func listTopics(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topics := getManager().ListTopics(accountId)
	return reply(ctx, topics)
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
	topicName := ctx.Param("topicName")
	subscriptionName := ctx.Param("subscriptionName")
	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	err := getManager().Subscribe(accountId, topicName, subscriptionName, req.Endpoint, req.FilterTag, req.NotifyStrategy, req.NotifyContentFormat)
	return replyObject(ctx, map[string]string{"subscriptionId": subscriptionName}, err)
}

func unsubscribe(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	subscriptionName := ctx.Param("subscriptionName")
	err := getManager().Unsubscribe(accountId, topicName, subscriptionName)
	return reply(ctx, err)
}

func getSubscription(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	subscriptionName := ctx.Param("subscriptionName")
	attr, err := getManager().GetSubscription(accountId, topicName, subscriptionName)
	return replyObject(ctx, attr, err)
}

func setSubscription(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	subscriptionName := ctx.Param("subscriptionName")
	attr := mns.Subscription{}
	if err := ctx.Bind(&attr); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	err := getManager().SetSubscription(accountId, topicName, subscriptionName, attr)
	return reply(ctx, err)
}

func listTopicSubscriptions(ctx echo.Context) error {
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")
	pageSize, err1 := strconv.Atoi(ctx.QueryParam("pageSize"))
	pageNo, err2 := strconv.Atoi(ctx.QueryParam("pageNo"))
	if err1 != nil || err2 != nil {
		return reply(ctx, mns.ErrInvalidArgument)
	}
	attrs, err := getManager().ListTopicSubscriptions(accountId, topicName, pageNo, pageSize)
	return replyObject(ctx, attrs, err)
}

func publishMessage(ctx echo.Context) error {
	req := struct {
		Body       []byte                 `json:"body" bson:"Body"`
		Tag        string                 `json:"tag" bson:"Tag"`
		Attributes map[string]interface{} `json:"attributes" bson:"Attributes"`
	}{}
	accountId := getAccount(ctx)
	topicName := ctx.Param("topicName")

	if err := ctx.Bind(&req); err != nil {
		return reply(ctx, mns.ErrInvalidArgument, err)
	}
	err := getManager().PublishMessage(accountId, topicName, req.Body, req.Tag, req.Attributes)
	return reply(ctx, err)
}

func publishNotification(ctx echo.Context) error {
	return mns.ErrNotImplemented
}
