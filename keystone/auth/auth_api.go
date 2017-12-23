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

func authenticate(opts interface{}) error {
	switch opts.(type) {
	case TenantAuthOption:
		opt := opts.(*TenantAuthOption)
		if signer, err := newSigner(opt.SignMethod); err == nil {
			return signer.signTenant(opt)
		}
	case DeviceAuthOption:
		opt := opts.(*DeviceAuthOption)
		if signer, err := newSigner(opt.SignMethod); err == nil {
			return signer.signDevice(opt)
		}
	}
	return ErrorAuthDenied
}

func authorize(opts interface{}) error {
	return nil
}
