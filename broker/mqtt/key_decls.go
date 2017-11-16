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

package mqtt

// MQTT specified session infor
const (
	InfoCleanSession       = "clean_session"
	InfoMessageMaxInflight = "inflight_max"
	InfoMessageInflight    = "message_in_flight"
	InfoMessageInQueue     = "message_in_queue"
	InfoMessageDropped     = "message_dropped"
	InfoAwaitingRel        = "message_awaiting_rel"
	InfoAwaitingComp       = "message_awaitng_comp"
	InfoAwaitingAck        = "message_aiwaiting _ack"
	InfoCreatedAt          = "created_at"
)

// Stats declarations
const (
	StatClientsMax         = "clients/max"
	StatClientsCount       = "client/count"
	StatQueuesMax          = "queues/max"
	StatQueuesCount        = "queues/count"
	StatRetainedMax        = "retained/max"
	StatRetainedCount      = "retained/count"
	StatSessionsMax        = "sessions/max"
	StatSessionsCount      = "sessions/count"
	StatSubscriptionsMax   = "subscriptions/max"
	StatSubscriptionsCount = "subscriptions/count"
	StatTopicsMax          = "topics/max"
	StatTopicsCount        = "topic/count"
)

// Metrics declarations
const (
	MetricBytesReceived         = "bytes/recevied"
	MetricBytesSent             = "bytes/sent"
	MetricMessageDroped         = "messages/droped"
	MetricMessageQos0Recevied   = "messages/qos0/received"
	MetricMessageQos0Sent       = "messages/qos0/sent"
	MetricMessageQos1Received   = "messages/qos1/recevied"
	MetricMessageQos1Sent       = "messages/qos1/sent"
	MetricMessageOos2Recevied   = "messages/qos2/received"
	MetricMessageOos2Sent       = "messages/qos2/sent"
	MetricMessageRetained       = "messages/retained"
	MetricMessageReceived       = "messages/received"
	MetricMessageSent           = "messages/sent"
	MetricPacketConnack         = "packets/connack"
	MetricPacketConnect         = "packets/connect"
	MetricPacketDisconnect      = "packets/disconnect"
	MetricPacketPingreq         = "packets/pingreq"
	MetricPacketPingresp        = "packets/pingresp"
	MetricPacketPubackRecevied  = "packets/puback/received"
	MetricPacketPubackSent      = "packets/puback/sent"
	MetricPacketPubcompReceived = "packets/pubcomp/received"
	MetricPacketPubcompSent     = "packets/pubcomp/sent"
	MetricPacketPublishReceived = "packets/publish/received"
	MetricPacketPublishSent     = "packets/publish/sent"
	MetricPacketPubrecReceived  = "packets/pubrec/received"
	MetricPacketPubrecSent      = "packets/pubrec/sent"
	MetricPacketPubrelReceived  = "packets/pubrel/received"
	MetricPacketPubrelSent      = "packets/pubrel/sent"
	MetricPacketReceived        = "packes/received"
	MetricPacketSent            = "packets/sent"
	MetricPacketSuback          = "packets/subback"
	MetricPacketSubscribe       = "packets/subscribe"
	MetricPacketUnsuback        = "packets/unsuback"
	MetricPacketUnsubscribe     = "packets/unsubscribe"
)
