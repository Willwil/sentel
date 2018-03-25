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

package mns

type MnsError struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (e MnsError) Error() string { return e.Message }

func NewError(status int, message string) *MnsError {
	return &MnsError{StatusCode: status, Message: message}
}

var (
	ErrNotImplemented              = NewError(501, "NotImplemented")
	ErrAccesDenied                 = NewError(403, "AccessDenied")
	ErrInvalidAccessId             = NewError(403, "InvalidAccessKeyId")
	ErrInternalError               = NewError(500, "InternalError")
	ErrInvalidAuthorizationHeader  = NewError(400, "InvalidAuthorizationHeader")
	ErrInvalidDateHeader           = NewError(400, "InvalidDateHeader")
	ErrInvalidArgument             = NewError(400, "InvalidArgument")
	ErrInvalidDegist               = NewError(400, "InvalidDegist")
	ErrInvalidRequestUrl           = NewError(400, "InvalidRequestURL")
	ErrInvalidQueryString          = NewError(400, "InvalidQueryString")
	ErrMissingAuthorizationHeader  = NewError(400, "MissingAuthorizationHeader")
	ErrMissingDateHeader           = NewError(400, "MissingDateHeader")
	ErrMissingVersionHeader        = NewError(400, "MissingVersionHeader")
	ErrMissingReceiptHandle        = NewError(400, "MissingReceiptHandle")
	ErrMissingVisiibilityTimeout   = NewError(400, "MissingVisibilityTimeout")
	ErrMessageNotExist             = NewError(409, "MessageNotExist")
	ErrQueueAlreadyExist           = NewError(400, "QueueAlreadyExist")
	ErrQueueDeletedRecently        = NewError(400, "QueueDeletedRecently")
	ErrInvalidQueueName            = NewError(400, "InvalidQueueName")
	ErrQueueNameLengthError        = NewError(400, "QueueNameLengthError")
	ErrQueueNotExist               = NewError(404, "QueueNotExist")
	ErrReceiptHandlerError         = NewError(400, "ReceiptHandleError")
	ErrSignatureDoesNotMatch       = NewError(400, "SignatureDoesNotMatch")
	ErrTimeExpired                 = NewError(408, "TimeExpired")
	ErrQpsLimitExceeded            = NewError(400, "QpsLimitExceeded")
	ErrTopicAlreadyExit            = NewError(400, "TopicAlreadyExist")
	ErrTopicNameInvalid            = NewError(400, "TopicNameInvalid")
	ErrTopicNameLengthError        = NewError(400, "TopicNameLengthError")
	ErrTopicNotExist               = NewError(404, "TopicNotExist")
	ErrSubscriptionNameInvalid     = NewError(400, "SubscriptionNameInvalid")
	ErrSubscriptionNameLengthError = NewError(400, "SubscriptionNameLengthError")
	ErrSubscriptionNotExist        = NewError(404, "SubscriptionNotExist")
	ErrSubscriptionAlreadyExist    = NewError(409, "SubscriptionAlreadyExist")
	ErrEndpointInvalid             = NewError(400, "EndpointInvalid")
)
