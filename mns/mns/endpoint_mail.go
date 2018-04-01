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

	"github.com/goinggo/mapstructure"
	"gopkg.in/gomail.v2"
)

type mailEndpoint struct {
	config    config.Config
	attribute EndpointAttribute
	toaddr    string
}

func newMailEndpoint(c config.Config, subscription Subscription) (endpoint Endpoint, err error) {
	toaddr, err := parseMailScheme(subscription.Endpoint)
	endpoint = &mailEndpoint{
		config: c,
		toaddr: toaddr,
		attribute: EndpointAttribute{
			Name:           subscription.Endpoint,
			Type:           "mail",
			URI:            subscription.Endpoint,
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
	return sendMail(m.config, m.toaddr, attr.Subject, string(body), "html")
}

func sendMail(c config.Config, to, subject, content, contentType string, attachments ...string) error {
	host, _ := c.StringWithSection("email", "host")
	port, _ := c.IntWithSection("email", "port")
	account, _ := c.StringWithSection("email", "account")
	pwd, _ := c.StringWithSection("email", "pwd")
	from, _ := c.StringWithSection("email", "from")
	if host == "" || port < 0 || account == "" || pwd == "" || from == "" {
		return ErrInvalidArgument
	}
	dialer := gomail.NewDialer(host, port, account, pwd)
	msg := gomail.NewMessage()
	msg.SetHeader("From", from)
	msg.SetHeader("To", to)
	msg.SetHeader("Subject", subject)
	msg.SetBody(contentType, content)
	for _, attachment := range attachments {
		msg.Attach(attachment)
	}
	return dialer.DialAndSend(msg)
}
