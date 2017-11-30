//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package mqtt

import "errors"

const (
	ServiceName = "mqtt"
)

const (
	mqttNetworkTcp       = "tcp"
	mqttNetworkTls       = "tls"
	mqttNetworkWebsocket = "ws"
	mqttNetworkHttps     = "https"
)

// Mqtt session state
const (
	mqttStateInvalid        = 0
	mqttStateNew            = 1
	mqttStateConnected      = 2
	mqttStateDisconnecting  = 3
	mqttStateConnectAsync   = 4
	mqttStateConnectPending = 5
	mqttStateConnectSrv     = 6
	mqttStateDisconnectWs   = 7
	mqttStateDisconnected   = 8
	mqttStateExpiring       = 9
)

// mqtt protocol
const (
	mqttProtocolInvalid = 0
	mqttProtocol31      = 1
	mqttProtocol311     = 2
	mqttProtocolS       = 3
)

var (
	mqttErrorInvalidProtocol = errors.New("Invalid protocol")
	mqttErrorInvalidVersion  = errors.New("Invalid protocol version")
	mqttErrorConnectPending  = errors.New("Connec pending")
	mqttErrorNoConnection    = errors.New("No connection")
	mqttErrorConnectRefused  = errors.New("Connection Refused")
	mqttErrorNotFound        = errors.New("Not found")
	mqttErrorNotSupported    = errors.New("Not supported")
	mqttErrorAutoFailed      = errors.New("Auth failed")
	mqttErrorUnkown          = errors.New("Unknown error")
)

// Message state
const (
	mqttMsgStateInvalid       = 0
	mqttMsgStatePublishQos0   = 1
	mqttMsgStatePublishQos1   = 2
	mqttMsgStateWaitPubAck    = 3
	mqttMsgStatePublishQos2   = 4
	mqttMsgStateWaitPubRec    = 5
	mqttMsgStateResendPubRel  = 6
	mqttMsgStateWaitPubRel    = 7
	mqttMsgStateResendPubComp = 8
	mqttMsgStateWaitPubComp   = 9
	mqttMsgStateSendPubRec    = 10
	mqttMsgStateQueued        = 11
)

// Message direction
const (
	mqttMessageDirectionIn  = 0
	mqttMessageDirectionOut = 1
)

/*
type mqttMessage struct {
	mid       uint16
	direction int
	topic     string
	payload   []uint8
	qos       uint8
	retain    bool
}
*/

// nameOfSessionState return name of session state
func nameOfSessionState(state uint8) string {
	switch state {
	case mqttStateNew:
		return "mqttStateNew"
	case mqttStateConnected:
		return "mqttStateConnected"
	case mqttStateDisconnecting:
		return "mqttStateDisconnecting"
	case mqttStateConnectAsync:
		return "mqttStateConnectAsync"
	case mqttStateConnectPending:
		return "mqttStateConnectPending"
	case mqttStateConnectSrv:
		return "mqttStateConnectSrv"
	case mqttStateDisconnectWs:
		return "mqttStateDisconnectWs"
	case mqttStateDisconnected:
		return "mqttStateDisconnected"
	case mqttStateExpiring:
		return "mqttStateExpiring"
	}
	return "mqttStateInvalid"
}

// nameOfMessageState return message state
func nameOfMessageState(state uint8) string {
	switch state {
	case mqttMsgStatePublishQos0:
		return "mqttMessageStatePublishQos0"
	case mqttMsgStatePublishQos1:
		return "mqttMessageStatePublishQos1"
	case mqttMsgStateWaitPubAck:
		return "mqttMessageStateWaitPubAck"
	case mqttMsgStatePublishQos2:
		return "mqttMessageStatePublishQos2"
	case mqttMsgStateWaitPubRec:
		return "mqttMessageStateWaitPubRec"
	case mqttMsgStateResendPubRel:
		return "mqttMessageStateResendPubRel"
	case mqttMsgStateWaitPubRel:
		return "mqttMessageStatePubRel"
	case mqttMsgStateResendPubComp:
		return "mqttMessageStateResendPubComp"
	case mqttMsgStateWaitPubComp:
		return "mqttMessageStateWaitPubComp"
	case mqttMsgStateSendPubRec:
		return "mqttMessageStateSendPubRec"
	case mqttMsgStateQueued:
		return "mqttMessageStateQueued"
	}
	return "mqttMessageStateInvalid"
}
