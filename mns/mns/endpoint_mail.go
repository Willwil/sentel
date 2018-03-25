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
)

type mailEndpoint struct {
	attribute EndpointAttribute
	mailaddr  string
	email     *mail.Email
}

func newMailEndpoint(c config.Config, uri string) (endpoint Endpoint, err error) {
	mailaddr, err := parseMailScheme(uri)
	endpoint = &mailEndpoint{
		mailaddr: mailaddr,
		attribute: EndpointAttribute{
			Name:           uri,
			Type:           "mail",
			URI:            uri,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	}
	return
}

func parseMailScheme(uri string) (string, error) {
	return uri, nil
}

func (m mailEndpoint) GetAttribute() EndpointAttribute { return m.attribute }
func (m mailEndpoint) PushMessage(msg Message) error   { return nil }
