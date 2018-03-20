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

import (
	"github.com/gogap/errors"
)

const (
	MNS_ERR_NS = "MNS"

	MNS_ERR_TEMPSTR = "client response status error,code: {{.resp.Code}}, message: {{.resp.Message}}, resource: {{.resource}} request id: {{.resp.RequestId}}, host id: {{.resp.HostId}}"
)

var (
	ERR_SIGN_MESSAGE_FAILED        = errors.TN(MNS_ERR_NS, 1, "sign message failed, {{.err}}")
	ERR_MARSHAL_MESSAGE_FAILED     = errors.TN(MNS_ERR_NS, 2, "marshal message filed, {{.err}}")
	ERR_GENERAL_AUTH_HEADER_FAILED = errors.TN(MNS_ERR_NS, 3, "general auth header failed, {{.err}}")

	ERR_CREATE_NEW_REQUEST_FAILED = errors.TN(MNS_ERR_NS, 4, "create new request failed, {{.err}}")
	ERR_SEND_REQUEST_FAILED       = errors.TN(MNS_ERR_NS, 5, "send request failed, {{.err}}")
	ERR_READ_RESPONSE_BODY_FAILED = errors.TN(MNS_ERR_NS, 6, "read response body failed, {{.err}}")

	ERR_UNMARSHAL_ERROR_RESPONSE_FAILED = errors.TN(MNS_ERR_NS, 7, "unmarshal error response failed, {{.err}}")
	ERR_UNMARSHAL_RESPONSE_FAILED       = errors.TN(MNS_ERR_NS, 8, "unmarshal response failed, {{.err}}")
	ERR_DECODE_BODY_FAILED              = errors.TN(MNS_ERR_NS, 9, "decode body failed, {{.err}}, body: \"{{.body}}\"")
	ERR_GET_BODY_DECODE_ELEMENT_ERROR   = errors.TN(MNS_ERR_NS, 10, "get body decode element error, local: {{.local}}, error: {{.err}}")

	ERR_MNS_ACCESS_DENIED                = errors.TN(MNS_ERR_NS, 100, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_ACCESS_KEY_ID        = errors.TN(MNS_ERR_NS, 101, MNS_ERR_TEMPSTR)
	ERR_MNS_INTERNAL_ERROR               = errors.TN(MNS_ERR_NS, 102, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_AUTHORIZATION_HEADER = errors.TN(MNS_ERR_NS, 103, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_DATE_HEADER          = errors.TN(MNS_ERR_NS, 104, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_ARGUMENT             = errors.TN(MNS_ERR_NS, 105, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_DEGIST               = errors.TN(MNS_ERR_NS, 106, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_REQUEST_URL          = errors.TN(MNS_ERR_NS, 107, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_QUERY_STRING         = errors.TN(MNS_ERR_NS, 108, MNS_ERR_TEMPSTR)
	ERR_MNS_MALFORMED_XML                = errors.TN(MNS_ERR_NS, 109, MNS_ERR_TEMPSTR)
	ERR_MNS_MISSING_AUTHORIZATION_HEADER = errors.TN(MNS_ERR_NS, 110, MNS_ERR_TEMPSTR)
	ERR_MNS_MISSING_DATE_HEADER          = errors.TN(MNS_ERR_NS, 111, MNS_ERR_TEMPSTR)
	ERR_MNS_MISSING_VERSION_HEADER       = errors.TN(MNS_ERR_NS, 112, MNS_ERR_TEMPSTR)
	ERR_MNS_MISSING_RECEIPT_HANDLE       = errors.TN(MNS_ERR_NS, 113, MNS_ERR_TEMPSTR)
	ERR_MNS_MISSING_VISIBILITY_TIMEOUT   = errors.TN(MNS_ERR_NS, 114, MNS_ERR_TEMPSTR)
	ERR_MNS_MESSAGE_NOT_EXIST            = errors.TN(MNS_ERR_NS, 115, MNS_ERR_TEMPSTR)
	ERR_MNS_QUEUE_ALREADY_EXIST          = errors.TN(MNS_ERR_NS, 116, MNS_ERR_TEMPSTR)
	ERR_MNS_QUEUE_DELETED_RECENTLY       = errors.TN(MNS_ERR_NS, 117, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_QUEUE_NAME           = errors.TN(MNS_ERR_NS, 118, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_VERSION_HEADER       = errors.TN(MNS_ERR_NS, 119, MNS_ERR_TEMPSTR)
	ERR_MNS_INVALID_CONTENT_TYPE         = errors.TN(MNS_ERR_NS, 120, MNS_ERR_TEMPSTR)
	ERR_MNS_QUEUE_NAME_LENGTH_ERROR      = errors.TN(MNS_ERR_NS, 121, MNS_ERR_TEMPSTR)
	ERR_MNS_QUEUE_NOT_EXIST              = errors.TN(MNS_ERR_NS, 122, MNS_ERR_TEMPSTR)
	ERR_MNS_RECEIPT_HANDLE_ERROR         = errors.TN(MNS_ERR_NS, 123, MNS_ERR_TEMPSTR)
	ERR_MNS_SIGNATURE_DOES_NOT_MATCH     = errors.TN(MNS_ERR_NS, 124, MNS_ERR_TEMPSTR)
	ERR_MNS_TIME_EXPIRED                 = errors.TN(MNS_ERR_NS, 125, MNS_ERR_TEMPSTR)
	ERR_MNS_QPS_LIMIT_EXCEEDED           = errors.TN(MNS_ERR_NS, 134, MNS_ERR_TEMPSTR)
	ERR_MNS_UNKNOWN_CODE                 = errors.TN(MNS_ERR_NS, 135, MNS_ERR_TEMPSTR)

	ERR_MNS_QUEUE_NAME_IS_TOO_LONG                 = errors.TN(MNS_ERR_NS, 126, "queue name is too long, the max length is 256")
	ERR_MNS_DELAY_SECONDS_RANGE_ERROR              = errors.TN(MNS_ERR_NS, 127, "queue delay seconds is not in range of (0~60480)")
	ERR_MNS_MAX_MESSAGE_SIZE_RANGE_ERROR           = errors.TN(MNS_ERR_NS, 128, "max message size is not in range of (1024~65536)")
	ERR_MNS_MSG_RETENTION_PERIOD_RANGE_ERROR       = errors.TN(MNS_ERR_NS, 129, "message retention period is not in range of (60~129600)")
	ERR_MNS_MSG_VISIBILITY_TIMEOUT_RANGE_ERROR     = errors.TN(MNS_ERR_NS, 130, "message visibility timeout is not in range of (1~43200)")
	ERR_MNS_MSG_POOLLING_WAIT_SECONDS_RANGE_ERROR  = errors.TN(MNS_ERR_NS, 131, "message poolling wait seconds is not in range of (0~30)")
	REE_MNS_GET_QUEUE_RET_NUMBER_RANGE_ERROR       = errors.TN(MNS_ERR_NS, 132, "get queue list param of ret number is not in range of (1~1000)")
	ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR = errors.TN(MNS_ERR_NS, 133, "mns queue already exist, and the attribute is the same, queue name: {{.name}}")
)
