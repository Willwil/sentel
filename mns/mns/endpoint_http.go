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
)

type httpEndpoint struct {
	attribute EndpointAttribute
	uri       string
}

func newHttpEndpoint(c config.Config, uri string) (endpoint Endpoint, err error) {
	uri, err = parseHttpScheme(uri)
	endpoint = &httpEndpoint{
		uri: uri,
		attribute: EndpointAttribute{
			Name:           uri,
			Type:           "http",
			URI:            uri,
			CreatedAt:      time.Now(),
			LastModifiedAt: time.Now(),
		},
	}
	return
}

func parseHttpScheme(uri string) (string, error) {
	return uri, nil
}

func (h httpEndpoint) GetAttribute() EndpointAttribute { return h.attribute }
func (h httpEndpoint) PushMessage(msg Message) error   { return nil }
