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

import "errors"

var (
	ErrorInvalidArgument = errors.New("invalid argument")
	ErrorAuthDenied      = errors.New("authentication denied")
)

type ApiAuthParam struct {
	Format      string `json:"format"`
	Version     string `json:"version"`
	AccessId    string `json:"accessId"`
	Signature   string `json:"sign"`
	SignMethod  string `json:"signMethod"`
	SignVersion string `json:"signVersion"`
	SignNonce   string `json:"singNonce"`
	Timestamp   string `json:"timestamp"`
}

func Authenticate(opts interface{}) error {
	switch opts.(type) {
	case ApiAuthParam:
		opt := opts.(*ApiAuthParam)
		if signer, err := newSigner(opt.SignMethod); err == nil {
			return signer.signApi(opt)
		}
	}
	return ErrorAuthDenied
}
