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

const (
	RoleTenant = "tenant"
	RoleDevice = "device"
)

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

type DeviceAuthParam struct {
	ClientId     string `json:"clientId"`
	DeviceName   string `json:"deviceName"`
	ProductKey   string `json:"productKey"`
	SecurityMode int    `json:"securityMode"`
	SignMethod   string `json:"signMethod"`
	Timestamp    string `json:"timestamp"`
	Sign         string `json:"sign"`
	DeviceSecret string `json:"deviceSecret"`
}

type DeviceAuthorizeParam struct {
}
