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
	"github.com/cloustone/sentel/pkg/mail"

	"github.com/goinggo/mapstructure"
)

type mailEndpoint struct {
	config    config.Config
	attribute EndpointAttribute
	fromaddr  string
	toaddr    string
	mailpwd   string
	host      string
}

func newMailEndpoint(c config.Config, attr SubscriptionAttribute) (endpoint Endpoint, err error) {
	fromaddr, err1 := c.StringWithSection("mail", "from_address")
	mailpwd, err2 := c.StringWithSection("mail", "password")
	host, err3 := c.StringWithSection("mail", "host")
	if err1 != nil || err2 != nil || err3 != nil {
		return nil, ErrInvalidArgument
	}
	toaddr, err := parseMailScheme(attr.Endpoint)
	endpoint = &mailEndpoint{
		config:   c,
		fromaddr: fromaddr,
		toaddr:   toaddr,
		mailpwd:  mailpwd,
		host:     host,
		attribute: EndpointAttribute{
			Name:           attr.Endpoint,
			Type:           "mail",
			URI:            attr.Endpoint,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	}
	return
}

func parseMailScheme(uri string) (string, error) {
	return uri, nil
}

type mailAttribute struct {
	AccountName    string `json:"AccountName"`
	Subject        string `json:"Subject"`
	AddressType    uint8  `json:"AddressType"`
	IsHtml         uint8  `json:"IsHtml"`
	ReplyToAddress uint8  `json:"ReplyToAddress"`
}

func (m mailEndpoint) GetAttribute() EndpointAttribute { return m.attribute }
func (m mailEndpoint) PushMessage(body []byte, tag string, attrs map[string]interface{}) error {
	attr := mailAttribute{}
	if err := mapstructure.Decode(attrs, &attr); err != nil {
		return ErrInvalidArgument.With(err.Error())
	}
	if attr.AccountName == "" ||
		attr.Subject == "" ||
		attr.AddressType > 1 ||
		attr.IsHtml > 1 ||
		attr.ReplyToAddress > 1 {
		return ErrInvalidArgument
	}
	return mail.SendMail(m.fromaddr, m.toaddr, m.mailpwd, attr.Subject, string(body), "", "", m.host)
}
