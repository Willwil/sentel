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
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gogap/errors"
)

const (
	AUTHORIZATION = "Authorization"
	CONTENT_TYPE  = "Content-Type"
	CONTENT_MD5   = "Content-MD5"
	MQ_VERSION    = "x-mns-version"
	HOST          = "Host"
	DATE          = "Date"
	KEEP_ALIVE    = "Keep-Alive"
)

type Credential interface {
	Signature(method Method, headers map[string]string, resource string) (signature string, err error)
	SetSecretKey(accessKeySecret string)
}

type mnsCredential struct {
	accessKeySecret string
}

func NewCredential(accessKeySecret string) *mnsCredential {
	credential := new(mnsCredential)
	credential.accessKeySecret = accessKeySecret
	return credential
}

func (p *mnsCredential) SetSecretKey(accessKeySecret string) {
	p.accessKeySecret = accessKeySecret
}

func (p *mnsCredential) Signature(method Method, headers map[string]string, resource string) (signature string, err error) {
	signItems := []string{}
	signItems = append(signItems, string(method))

	contentMD5 := ""
	contentType := ""
	date := time.Now().UTC().Format(http.TimeFormat)

	if v, exist := headers[CONTENT_MD5]; exist {
		contentMD5 = v
	}

	if v, exist := headers[CONTENT_TYPE]; exist {
		contentType = v
	}

	if v, exist := headers[DATE]; exist {
		date = v
	}

	mnsHeaders := []string{}

	for k, v := range headers {
		if strings.HasPrefix(k, "x-mns-") {
			mnsHeaders = append(mnsHeaders, k+":"+strings.TrimSpace(v))
		}
	}

	sort.Sort(sort.StringSlice(mnsHeaders))

	stringToSign := string(method) + "\n" +
		contentMD5 + "\n" +
		contentType + "\n" +
		date + "\n" +
		strings.Join(mnsHeaders, "\n") + "\n" +
		resource

	sha1Hash := hmac.New(sha1.New, []byte(p.accessKeySecret))
	if _, e := sha1Hash.Write([]byte(stringToSign)); e != nil {
		err = ERR_SIGN_MESSAGE_FAILED.New(errors.Params{"err": e})
		return
	}

	signature = base64.StdEncoding.EncodeToString(sha1Hash.Sum(nil))

	return
}
