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
	"net/http"

	"github.com/gogap/errors"
)

func send(client MnsClient, decoder MnsDecoder, method Method, headers map[string]string, message interface{}, resource string, v interface{}) (statusCode int, err error) {
	var resp *http.Response
	if resp, err = client.Send(method, headers, message, resource); err != nil {
		return
	}

	if resp != nil {
		defer resp.Body.Close()
		statusCode = resp.StatusCode

		if resp.StatusCode != http.StatusCreated &&
			resp.StatusCode != http.StatusOK &&
			resp.StatusCode != http.StatusNoContent {

			errResp := ErrorMessageResponse{}
			if e := decoder.Decode(resp.Body, &errResp); e != nil {
				err = ERR_UNMARSHAL_ERROR_RESPONSE_FAILED.New(errors.Params{"err": e})
				return
			}
			err = ParseError(errResp, resource)
			return
		}

		if v != nil {
			if e := decoder.Decode(resp.Body, v); e != nil {
				err = ERR_UNMARSHAL_RESPONSE_FAILED.New(errors.Params{"err": e})
				return
			}
		}
	}

	return
}
