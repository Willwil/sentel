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
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gogap/errors"
)

type MnsLocation string

const (
	Beijing   MnsLocation = "cn-beijing"
	Hangzhou  MnsLocation = "cn-hangzhou"
	Qingdao   MnsLocation = "cn-qingdao"
	Singapore MnsLocation = "ap-southeast-1"
)

type QueueManager interface {
	CreateQueue(location MnsLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error)
	SetQueueAttributes(location MnsLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error)
	GetQueueAttributes(location MnsLocation, queueName string) (attr QueueAttribute, err error)
	DeleteQueue(location MnsLocation, queueName string) (err error)
	ListQueue(location MnsLocation, nextMarker Base64Bytes, retNumber int32, prefix string) (queues Queues, err error)
}

type queueManager struct {
	ownerId         string
	credential      Credential
	accessKeyId     string
	accessKeySecret string

	decoder MnsDecoder
}

// makeMnsUrl generate url for queue service in mns
// ownerId and location will be used in future
func makeMnsUrl(ownerId string, location string) string {
	formatOfMnsServer := "http://localhost:9091/mns/api/v1"
	return formatOfMnsServer
}

func checkQueueName(queueName string) (err error) {
	if len(queueName) > 256 {
		err = ERR_MNS_QUEUE_NAME_IS_TOO_LONG.New()
		return
	}
	return
}

func checkDelaySeconds(seconds int32) (err error) {
	if seconds > 60480 || seconds < 0 {
		err = ERR_MNS_DELAY_SECONDS_RANGE_ERROR.New()
		return
	}
	return
}

func checkMaxMessageSize(maxSize int32) (err error) {
	if maxSize < 1024 || maxSize > 65536 {
		err = ERR_MNS_MAX_MESSAGE_SIZE_RANGE_ERROR.New()
		return
	}
	return
}

func checkMessageRetentionPeriod(retentionPeriod int32) (err error) {
	if retentionPeriod < 60 || retentionPeriod > 1296000 {
		err = ERR_MNS_MSG_RETENTION_PERIOD_RANGE_ERROR.New()
		return
	}
	return
}

func checkVisibilityTimeout(visibilityTimeout int32) (err error) {
	if visibilityTimeout < 1 || visibilityTimeout > 43200 {
		err = ERR_MNS_MSG_VISIBILITY_TIMEOUT_RANGE_ERROR.New()
		return
	}
	return
}

func checkPollingWaitSeconds(pollingWaitSeconds int32) (err error) {
	if pollingWaitSeconds < 0 || pollingWaitSeconds > 30 {
		err = ERR_MNS_MSG_POOLLING_WAIT_SECONDS_RANGE_ERROR.New()
		return
	}
	return
}

func NewqueueManager(ownerId, accessKeyId, accessKeySecret string) QueueManager {
	return &queueManager{
		ownerId:         ownerId,
		accessKeyId:     accessKeyId,
		accessKeySecret: accessKeySecret,
		decoder:         NewMnsDecoder(),
	}
}

func checkAttributes(delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error) {
	if err = checkDelaySeconds(delaySeconds); err != nil {
		return
	}
	if err = checkMaxMessageSize(maxMessageSize); err != nil {
		return
	}
	if err = checkMessageRetentionPeriod(messageRetentionPeriod); err != nil {
		return
	}
	if err = checkVisibilityTimeout(visibilityTimeout); err != nil {
		return
	}
	if err = checkPollingWaitSeconds(pollingWaitSeconds); err != nil {
		return
	}
	return
}

func (p *queueManager) CreateQueue(location MnsLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	if err = checkAttributes(delaySeconds,
		maxMessageSize,
		messageRetentionPeriod,
		visibilityTimeout,
		pollingWaitSeconds); err != nil {
		return
	}

	message := CreateQueueRequest{
		DelaySeconds:           delaySeconds,
		MaxMessageSize:         maxMessageSize,
		MessageRetentionPeriod: messageRetentionPeriod,
		VisibilityTimeout:      visibilityTimeout,
		PollingWaitSeconds:     pollingWaitSeconds,
	}

	url := makeMnsUrl(p.ownerId, string(location))
	cli := NewMnsClient(url, p.accessKeyId, p.accessKeySecret)

	var code int
	code, err = send(cli, p.decoder, PUT, nil, &message, "queues/"+queueName, nil)

	if code == http.StatusNoContent {
		err = ERR_MNS_QUEUE_ALREADY_EXIST_AND_HAVE_SAME_ATTR.New(errors.Params{"name": queueName})
		return
	}

	return
}

func (p *queueManager) SetQueueAttributes(location MnsLocation, queueName string, delaySeconds int32, maxMessageSize int32, messageRetentionPeriod int32, visibilityTimeout int32, pollingWaitSeconds int32) (err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	if err = checkAttributes(delaySeconds,
		maxMessageSize,
		messageRetentionPeriod,
		visibilityTimeout,
		pollingWaitSeconds); err != nil {
		return
	}

	message := CreateQueueRequest{
		DelaySeconds:           delaySeconds,
		MaxMessageSize:         maxMessageSize,
		MessageRetentionPeriod: messageRetentionPeriod,
		VisibilityTimeout:      visibilityTimeout,
		PollingWaitSeconds:     pollingWaitSeconds,
	}

	url := makeMnsUrl(p.ownerId, string(location))
	cli := NewMnsClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = send(cli, p.decoder, PUT, nil, &message, fmt.Sprintf("queues/%s?metaoverride=true", queueName), nil)
	return
}

func (p *queueManager) GetQueueAttributes(location MnsLocation, queueName string) (attr QueueAttribute, err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	url := makeMnsUrl(p.ownerId, string(location))
	cli := NewMnsClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = send(cli, p.decoder, GET, nil, nil, "queues/"+queueName, &attr)

	return
}

func (p *queueManager) DeleteQueue(location MnsLocation, queueName string) (err error) {
	queueName = strings.TrimSpace(queueName)

	if err = checkQueueName(queueName); err != nil {
		return
	}

	url := makeMnsUrl(p.ownerId, string(location))
	cli := NewMnsClient(url, p.accessKeyId, p.accessKeySecret)

	_, err = send(cli, p.decoder, DELETE, nil, nil, "queues/"+queueName, nil)

	return
}

func (p *queueManager) ListQueue(location MnsLocation, nextMarker Base64Bytes, retNumber int32, prefix string) (queues Queues, err error) {
	url := makeMnsUrl(p.ownerId, string(location))
	cli := NewMnsClient(url, p.accessKeyId, p.accessKeySecret)

	header := map[string]string{}

	marker := ""
	if nextMarker != nil && len(nextMarker) > 0 {
		marker = strings.TrimSpace(string(nextMarker))
		if marker != "" {
			header["x-mns-marker"] = marker
		}
	}

	if retNumber > 0 {
		if retNumber >= 1 && retNumber <= 1000 {
			header["x-mns-ret-number"] = strconv.Itoa(int(retNumber))
		} else {
			err = REE_MNS_GET_QUEUE_RET_NUMBER_RANGE_ERROR.New()
			return
		}
	}

	prefix = strings.TrimSpace(prefix)
	if prefix != "" {
		header["x-mns-prefix"] = prefix
	}

	_, err = send(cli, p.decoder, GET, header, nil, "queues", &queues)

	return
}
