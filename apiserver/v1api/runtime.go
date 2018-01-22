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

package v1api

import "github.com/labstack/echo"

// Http Runtime Api

type deviceMessage struct {
	ProductId string `json:"productId"`
	DeviceId  string `json:"deviceId"`
	Payload   []byte `json:"payload"`
	Qos       uint8  `json:"qos"`
	Retain    bool   `json:"retain"`
}

// Send a could-to-device message
func SendMessageToDevice(ctx echo.Context) error {
	return nil
}

func BroadcastProductMessage(ctx echo.Context) error {
	return nil
}
