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

package apiservice

import (
	"strings"

	"github.com/labstack/echo"
)

const (
	AUTHORIZATION = "Authorization"
	CONTENT_TYPE  = "Content-Type"
	CONTENT_MD5   = "Content-MD5"
	MQ_VERSION    = "x-mns-version"
	DATE          = "Date"
	KEEP_ALIVE    = "Keep-Alive"
)

type apiToken struct {
	method        string
	accessId      string
	contentMD5    string
	contentType   string
	date          string
	headers       []string
	resource      string
	signature     string
	authenticated bool
}

func newApiToken(ctx echo.Context) apiToken {
	req := ctx.Request()
	token := apiToken{
		method:   req.Method,
		headers:  []string{},
		resource: req.URL.Path,
	}
	if v := req.Header.Get(AUTHORIZATION); v != "" {
		items := strings.Split(v, " ")
		if len(items) == 2 {
			values := strings.Split(items[1], ":")
			if len(values) == 2 {
				token.accessId = values[0]
				token.signature = values[1]
			}
		}
	}
	token.contentMD5 = req.Header.Get(CONTENT_MD5)
	token.contentType = req.Header.Get(CONTENT_TYPE)
	token.date = req.Header.Get(DATE)
	for k, v := range req.Header {
		if strings.HasPrefix(k, "x-mns-") {
			token.headers = append(token.headers, k+":"+strings.Join(v, ""))
		}
	}
	return token
}

func (p apiToken) GetPrincipal() interface{} {
	return p.accessId
}
func (p apiToken) GetCrenditals() interface{} {
	return p.signature
}

func (p apiToken) IsAuthenticated() bool {
	return p.authenticated
}
