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
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/sms"

	"github.com/goinggo/mapstructure"
)

type smsEndpoint struct {
	attribute EndpointAttribute
	smsaddr   string
	sms       sms.SMS
}

func newSMSEndpoint(c config.Config, attr SubscriptionAttribute) (endpoint Endpoint, err error) {
	addr, err := parseSMSScheme(attr.Endpoint)
	endpoint = &smsEndpoint{
		smsaddr: addr,
		attribute: EndpointAttribute{
			Name:           attr.Endpoint,
			Type:           "sms",
			URI:            attr.Endpoint,
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
	sms := smsAttribute{}
	if err := mapstructure.Decode(attrs, &sms); err != nil {
		return ErrInvalidArgument.With(err.Error())
	}
	if sms.FreeSignName == "" ||
		sms.TemplateCode == "" ||
		(sms.Type != "singleContent" && sms.Type != "multiContent") ||
		sms.Params == "" {
		return ErrInvalidArgument
	}

	return nil

}
