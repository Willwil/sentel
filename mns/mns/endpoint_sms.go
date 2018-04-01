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

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/sms"

	"github.com/goinggo/mapstructure"
)

type smsEndpoint struct {
	config    config.Config
	attribute EndpointAttribute
	to        string
}

func newSMSEndpoint(c config.Config, subscription Subscription) (endpoint Endpoint, err error) {
	to, err := parseSMSScheme(subscription.Endpoint)
	endpoint = &smsEndpoint{
		config: c,
		to:     to,
		attribute: EndpointAttribute{
			Name:           subscription.Endpoint,
			Type:           "sms",
			URI:            subscription.Endpoint,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	}
	return
}

func parseSMSScheme(uri string) (string, error) {
	return uri, nil
}

type smsAttribute struct {
	FreeSignName string `json:"FreeSignName" bson:"FreeSignName"`
	TemplateCode string `json:"TemplateCode" bson:"TemplateCode"`
	Type         string `json:"Type" bson:"Type"`
	Receiver     string `json:"Receiver" bson:"Receiver"`
	Params       string `json:"SmsParams, omitempty" bson:"SmsParams"`
}

func (s smsEndpoint) GetAttribute() EndpointAttribute { return s.attribute }

func (s smsEndpoint) PushMessage(body []byte, tag string, attrs map[string]interface{}) error {
	attr := smsAttribute{}
	if err := mapstructure.Decode(attrs, &attr); err != nil {
		return ErrInvalidArgument.With(err.Error())
	}
	if attr.FreeSignName == "" ||
		attr.TemplateCode == "" ||
		(attr.Type != "singleContent" && attr.Type != "multiContent") ||
		attr.Params == "" {
		return ErrInvalidArgument
	}

	return sendSMS(s.config, s.to, attr)
}

func sendSMS(c config.Config, to string, attr smsAttribute) error {
	appid, _ := c.StringWithSection("sms", "appid")
	appkey, _ := c.StringWithSection("sms", "appkey")
	signtype, _ := c.StringWithSection("sms", "signtype")
	project, _ := c.StringWithSection("sms", "project")
	if appid == "" || appkey == "" || signtype == "" {
		return errors.New("invalid submail configs setting")
	}
	configs := make(map[string]string)
	configs["appid"] = appid
	configs["appkey"] = appkey
	configs["signtype"] = signtype
	dialer := sms.NewDialer(configs)
	msg := sms.NewMessage()
	msg.AddTo(to)
	msg.SetProject(project)
	values := make(map[string]string)
	if err := json.Unmarshal([]byte(attr.Params), values); err != nil {
		return ErrInvalidArgument.With(err.Error())
	}
	for k, v := range values {
		msg.AddVariable(k, v)
	}
	return dialer.DialAndSend(msg)
}
