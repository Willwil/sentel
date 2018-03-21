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

package client

type MessageResponse struct {
	Code      string `json:"code,omitempty"`
	Message   string `json:"message,omitempty"`
	RequestId string `json:"request_id,omitempty"`
}

type ErrorMessageResponse struct {
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

type QueueAttribute struct {
	QueueName              string `json:"queue_name,omitempty"`
	DelaySeconds           int32  `json:"delay_senconds,omitempty"`
	MaxMessageSize         int32  `json:"maximum_message_size,omitempty"`
	MessageRetentionPeriod int32  `json:"message_retention_period,omitempty"`
	VisibilityTimeout      int32  `json:"visibility_timeout,omitempty"`
	PollingWaitSeconds     int32  `json:"polling_wait_secods,omitempty"`
	ActiveMessages         int64  `json:"active_messages,omitempty"`
	InactiveMessages       int64  `json:"inactive_messages,omitempty"`
	DelayMessages          int64  `json:"delay_messages,omitempty"`
	CreateTime             int64  `json:"create_time,omitempty"`
	LastModifyTime         int64  `json:"last_modify_time,omitempty"`
}

type Queue struct {
	QueueURL string `json:"url"`
}

type Queues struct {
	Queues     []Queue     `json:"queues"`
	NextMarker Base64Bytes `json:"next_marker"`
}

type Base64Bytes []byte
