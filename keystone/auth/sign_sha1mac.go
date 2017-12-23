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
)

type sha1macSigner struct{}

func (p *sha1macSigner) name() string { return "sha1mac" }

func (p *sha1macSigner) signTenant(opt *TenantAuthOption) error {
	return nil
}
func (p *sha1macSigner) signDevice(opt *DeviceAuthOption) error {
	content := "clientId" + opt.ClientId + "deviceName" + opt.DeviceName + "productKey" + opt.ProductKey + "timestamp" + opt.Timestamp
	mac := hmac.New(sha1.New, []byte(opt.DeviceSecret))
	mac.Write([]byte(content))
	result := mac.Sum(nil)
	if opt.Sign == string(result) {
		return nil
	}
	return ErrorAuthDenied
}
