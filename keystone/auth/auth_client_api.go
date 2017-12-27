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

package auth

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cloustone/sentel/pkg/config"
)

type authResponse struct {
	Success bool `json:"success"`
}

func Authenticate(c config.Config, opts interface{}) error {
	hosts := c.MustString("keystone", "hosts")
	format := "application/json;charset=utf-8"

	switch opts.(type) {
	case ApiAuthParam:
		if buf, err := json.Marshal(opts); err == nil {
			url := "http://" + hosts + "/keystone/api/v1/auth/api"
			req := bytes.NewBuffer([]byte(buf))
			if resp, err := http.Post(url, format, req); err == nil {
				body, _ := ioutil.ReadAll(resp.Body)
				result := authResponse{}
				if err := json.Unmarshal(body, &result); err == nil && result.Success == true {
					return nil
				}
			}
		}
	case DeviceAuthParam:
		if buf, err := json.Marshal(opts); err == nil {
			url := "http://" + hosts + "/keystone/api/v1/auth/device"
			req := bytes.NewBuffer([]byte(buf))
			if resp, err := http.Post(url, format, req); err == nil {
				body, _ := ioutil.ReadAll(resp.Body)
				result := authResponse{}
				if err := json.Unmarshal(body, &result); err == nil && result.Success == true {
					return nil
				}
			}
		}
	}
	return ErrorAuthDenied
}

func Authorize(c config.Config, opts interface{}) error {
	return nil
}
