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
	"crypto/hmac"
	"crypto/sha1"
	"fmt"
)

var signMethods = map[string]func(*Options) error{
	"sha1mac": sha1mac,
}

func sign(opt *Options) error {
	if _, ok := signMethods[opt.SignMethod]; !ok {
		return fmt.Errorf("auth: Invalid sign method '%s'", opt.SignMethod)
	}
	return signMethods[opt.SignMethod](opt)
}

func sha1mac(opts *Options) error {
	content := "clientId" + opts.ClientId + "deviceName" + opts.DeviceName + "productKey" + opts.ProductKey + "timestamp" + opts.Timestamp
	mac := hmac.New(sha1.New, []byte(opts.DeviceSecret))
	mac.Write([]byte(content))
	result := mac.Sum(nil)
	if opts.Password == string(result) {
		return nil
	}
	return ErrorAuthDenied
}
