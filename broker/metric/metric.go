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

package metric

const (
	BytesReceived         = "bytes/recevied"
	BytesSent             = "bytes/sent"
	MessageMaxInflight    = "messages/inflight/max"
	MessageInflight       = "messages/inflight"
	MessageInQueue        = "messages/inqueue"
	MessageAwaitingRel    = "messages/awaiting/rel"
	MessageAwaitingComp   = "messages/awaitng/comp"
	MessagewaitingAck     = "messages/aiwaiting/ack"
	MessageDroped         = "messages/droped"
	MessageQos0Recevied   = "messages/qos0/received"
	MessageQos0Sent       = "messages/qos0/sent"
	MessageQos1Received   = "messages/qos1/recevied"
	MessageQos1Sent       = "messages/qos1/sent"
	MessageOos2Recevied   = "messages/qos2/received"
	MessageOos2Sent       = "messages/qos2/sent"
	MessageRetained       = "messages/retained"
	MessageReceived       = "messages/received"
	MessageSent           = "messages/sent"
	PacketConnack         = "packets/connack"
	PacketConnect         = "packets/connect"
	PacketDisconnect      = "packets/disconnect"
	PacketPingreq         = "packets/pingreq"
	PacketPingresp        = "packets/pingresp"
	PacketPubackRecevied  = "packets/puback/received"
	PacketPubackSent      = "packets/puback/sent"
	PacketPubcompReceived = "packets/pubcomp/received"
	PacketPubcompSent     = "packets/pubcomp/sent"
	PacketPublishReceived = "packets/publish/received"
	PacketPublishSent     = "packets/publish/sent"
	PacketPubrecReceived  = "packets/pubrec/received"
	PacketPubrecSent      = "packets/pubrec/sent"
	PacketPubrelReceived  = "packets/pubrel/received"
	PacketPubrelSent      = "packets/pubrel/sent"
	PacketReceived        = "packes/received"
	PacketSent            = "packets/sent"
	PacketSuback          = "packets/subback"
	PacketSubscribe       = "packets/subscribe"
	PacketUnsuback        = "packets/unsuback"
	PacketUnsubscribe     = "packets/unsubscribe"
)

type Metric interface {
	Get() map[string]uint64
	Add(name string, value uint64)
	Sub(name string, value uint64)
	AddMetric(metric Metric)
}
