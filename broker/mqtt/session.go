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

import (
	"context"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"

	auth "github.com/cloustone/sentel/broker/auth"
	uuid "github.com/satori/go.uuid"

	"github.com/golang/glog"
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

type mqttSession struct {
	mgr               *mqttService
	config            core.Config
	storage           Storage
	authapi           auth.IAuthAPI
	conn              net.Conn
	id                string
	clientID          string
	state             uint8
	inpacket          mqttPacket
	bytesReceived     int64
	pingTime          *time.Time
	address           string
	keepalive         uint16
	protocol          uint8
	observer          base.SessionObserver
	username          string
	password          string
	lastMessageIn     time.Time
	lastMessageOut    time.Time
	cleanSession      uint8
	isDroping         bool
	willMsg           *mqttMessage
	stateMutex        sync.Mutex
	sendStopChannel   chan int
	sendPacketChannel chan *mqttPacket
	sendMsgChannel    chan *mqttMessage
	waitgroup         sync.WaitGroup
	stats             *base.Stats
	metrics           *base.Metrics

	// resume field
	msgs       []*mqttMessage
	storedMsgs []*mqttMessage
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqttService, conn net.Conn, id string) (*mqttSession, error) {
	// Get session message queue size, if it is not set, default is 10
	qsize, err := m.Config.Int(m.protocol, "session_queue_size")
	if err != nil {
		qsize = 10
	}
	authapi, err := auth.NewAuthApi(m.Config)
	if err != nil {
		return nil, err
	}
	var msgqsize int
	msgqsize, err = m.Config.Int(m.protocol, "session_msg_queue_size")
	if err != nil {
		msgqsize = 10
	}

	s := &mqttSession{
		mgr:               m,
		config:            m.Config,
		storage:           m.storage,
		conn:              conn,
		id:                id,
		bytesReceived:     0,
		state:             mqttStateNew,
		inpacket:          newMqttPacket(),
		protocol:          mqttProtocolInvalid,
		observer:          nil,
		sendStopChannel:   make(chan int),
		sendPacketChannel: make(chan *mqttPacket, qsize),
		authapi:           authapi,
		stats:             base.NewStats(true),
		metrics:           base.NewMetrics(true),
		sendMsgChannel:    make(chan *mqttMessage, msgqsize),
		msgs:              make([]*mqttMessage, msgqsize),
		storedMsgs:        make([]*mqttMessage, msgqsize),
	}

	return s, nil
}

func (this *mqttSession) Identifier() string    { return this.id }
func (this *mqttSession) Service() core.Service { return this.mgr }
func (this *mqttSession) RegisterObserver(o base.SessionObserver) {
	if this.observer != nil {
		glog.Error("MqttSession register multiple observer")
	}
	this.observer = o
}
func (this *mqttSession) GetStats() *base.Stats     { return this.stats }
func (this *mqttSession) GetMetrics() *base.Metrics { return this.metrics }
func (this *mqttSession) Info() *base.SessionInfo   { return nil }

// launchPacketSendHandler launch goroutine to send packet queued for client
func (this *mqttSession) launchPacketSendHandler() {
	go func(stopChannel chan int, packetChannel chan *mqttPacket, msgChannel chan *mqttMessage) {
		defer this.waitgroup.Add(1)

		for {
			select {
			case <-stopChannel:
				return
			case p := <-packetChannel:
				for p.toprocess > 0 {
					len, err := this.conn.Write(p.payload[p.pos:p.toprocess])
					if err != nil {
						glog.Fatal("Failed to send packet to '%s:%s'", this.id, err)
						return
					}
					if len > 0 {
						p.toprocess -= len
						p.pos += len
					} else {
						glog.Fatal("Failed to send packet to '%s'", this.id)
						return
					}
				}
			case msg := <-msgChannel:
				this.processMessage(msg)
			case <-time.After(1 * time.Second):
				this.processTimeout()
			}
		}
	}(this.sendStopChannel, this.sendPacketChannel, this.sendMsgChannel)
}

// processMessage proceess messages
func (this *mqttSession) processMessage(msg *mqttMessage) error {

	return nil
}

// processTimeout proceess timeout
func (this *mqttSession) processTimeout() error {

	return nil
}

// Handle is mainprocessor for iot device client
func (this *mqttSession) Handle() error {

	glog.Infof("Handling session:%s", this.id)
	defer this.Destroy()

	this.launchPacketSendHandler()
	for {
		var err error
		if err = this.inpacket.DecodeFromReader(this.conn, base.NilDecodeFeedback{}); err != nil {
			glog.Error(err)
			return err
		}
		switch this.inpacket.command & 0xF0 {
		case PINGREQ:
			err = this.handlePingReq()
		case CONNECT:
			err = this.handleConnect()
		case DISCONNECT:
			err = this.handleDisconnect()
		case PUBLISH:
			err = this.handlePublish()
		case PUBREL:
			err = this.handlePubRel()
		case SUBSCRIBE:
			err = this.handleSubscribe()
		case UNSUBSCRIBE:
			err = this.handleUnsubscribe()
		default:
			err = fmt.Errorf("Unrecognized protocol command:%d", int(this.inpacket.command&0xF0))
		}
		if err != nil {
			glog.Error(err)
			return err
		}
		// Check sesstion state
		if this.state == mqttStateDisconnected {
			break
		}
		this.inpacket.Clear()
	}
	return nil
}

// Destroy will destory the current session
func (this *mqttSession) Destroy() error {
	// Stop packet sender goroutine
	this.sendStopChannel <- 1
	this.waitgroup.Wait()
	if this.conn != nil {
		this.conn.Close()
	}
	this.mgr.removeSession(this)
	return nil
}

// generateId generate id fro session or client
func (this *mqttSession) generateId() string {
	return uuid.NewV4().String()
}

// handlePingReq handle ping request packet
func (this *mqttSession) handlePingReq() error {
	glog.Infof("Received PINGREQ from %s", this.Identifier())
	return this.sendPingRsp()
}

// handleConnect handle connect packet
func (this *mqttSession) handleConnect() error {
	glog.Infof("Handling CONNECT packet from %s", this.id)

	if this.state != mqttStateNew {
		return errors.New("Invalid session state")
	}
	// Check protocol name and version
	protocolName, err := this.inpacket.readString()
	if err != nil {
		return err
	}
	protocolVersion, err := this.inpacket.readByte()
	if err != nil {
		return err
	}
	switch protocolName {
	case PROTOCOL_NAME_V31:
		if protocolVersion&0x7F != PROTOCOL_VERSION_V31 {
			this.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		this.protocol = mqttProtocol311

	case PROTOCOL_NAME_V311:
		if protocolVersion&0x7F != PROTOCOL_VERSION_V311 {
			this.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		// Reserved flags is not set to 0, must disconnect
		if this.inpacket.command&0x0F != 0x00 {
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		this.protocol = mqttProtocol311
	default:
		return fmt.Errorf("Invalid protocol name '%s' in CONNECT packet", protocolName)
	}

	// Check connect flags
	cflags, err := this.inpacket.readByte()
	if err != nil {
		return nil
	}
	/*
		if this.mgr.protocol == mqttProtocol311 {
			if cflags&0x01 != 0x00 {
				return errorthis.New("Invalid protocol version in connect flags")
			}
		}
	*/
	cleanSession := (cflags & 0x02) >> 1
	will := cflags & 0x04
	willQos := (cflags & 0x18) >> 3
	if willQos >= 3 { // qos level3 is not supported
		return fmt.Errorf("Invalid Will Qos in CONNECT from %s", this.id)
	}

	willRetain := (cflags & 0x20) == 0x20
	passwordFlag := cflags & 0x40
	usernameFlag := cflags & 0x80
	keepalive, err := this.inpacket.readUint16()
	if err != nil {
		return err
	}
	this.keepalive = keepalive

	// Deal with client identifier
	clientid, err := this.inpacket.readString()
	if err != nil {
		return err
	}
	if clientid == "" {
		if this.protocol == mqttProtocol31 {
			this.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
		} else {
			if option, err := this.config.Bool(this.mgr.protocol, "allow_zero_length_clientid"); err != nil && (!option || cleanSession == 0) {
				this.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
				return errors.New("Invalid mqtt packet with client id")
			}

			clientid = this.generateId()
		}
	}

	// TODO: clientid_prefixes check

	// Deal with topc
	var willTopic string
	var willMsg *mqttMessage
	var payload []uint8

	if will > 0 {
		willMsg = new(mqttMessage)
		// Get topic
		topic, err := this.inpacket.readString()
		if err != nil || topic == "" {
			return nil
		}
		willTopic = topic
		if this.observer != nil {
			willTopic = this.observer.OnGetMountPoint() + topic
		}
		if err := checkTopicValidity(willTopic); err != nil {
			return err
		}
		// Get willtopic's payload
		willPayloadLength, err := this.inpacket.readUint16()
		if err != nil {
			return err
		}
		if willPayloadLength > 0 {
			payload, err = this.inpacket.readBytes(int(willPayloadLength))
			if err != nil {
				return err
			}
		}
	} else {
		if this.protocol == mqttProtocol311 {
			if willQos != 0 || willRetain {
				return mqttErrorInvalidProtocol
			}
		}
	} // else will

	var username string
	var password string
	if usernameFlag > 0 {
		username, err = this.inpacket.readString()
		if err == nil {
			if passwordFlag > 0 {
				password, err = this.inpacket.readString()
				if err == mqttErrorInvalidProtocol {
					if this.protocol == mqttProtocol31 {
						passwordFlag = 0
					} else if this.protocol == mqttProtocol311 {
						return err
					} else {
						return err
					}
				}
			}
		} else {
			if this.protocol == mqttProtocol31 {
				usernameFlag = 0
			} else {
				return err
			}
		}
	} else { // username flag
		if this.protocol == mqttProtocol311 {
			if passwordFlag > 0 {
				return mqttErrorInvalidProtocol
			}
		}
	}

	if usernameFlag > 0 {
		if this.observer != nil {
			err := this.observer.OnAuthenticate(this, username, password)
			switch err {
			case nil:
			case base.IotErrorAuthFailed:
				this.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
				this.disconnect()
				return err
			default:
				this.disconnect()
				return err

			}
			// Get username and passowrd sucessfuly
			this.username = username
			this.password = password
		}
		// Get anonymous allow configuration
		allowAnonymous, _ := this.config.Bool(this.mgr.protocol, "allow_anonymous")
		if usernameFlag > 0 && allowAnonymous == false {
			// Dont allow anonymous client connection
			this.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	// Check wether username will be used as client id,
	// The connection request will be refused if the option is set
	if option, err := this.config.Bool(this.mgr.protocol, "user_name_as_client_id"); err != nil && option {
		if this.username != "" {
			clientid = this.username
		} else {
			this.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	conack := 0
	// Find if the client already has an entry, this must be done after any security check
	if found, _ := this.storage.FindSession(clientid); found != nil {
		// Found old session
		if found.state == mqttStateInvalid {
			glog.Errorf("Invalid session(%s) in store", found.id)
		}
		if this.protocol == mqttProtocol311 {
			if cleanSession == 0 {
				conack |= 0x01
			}
		}
		this.cleanSession = cleanSession

		if this.cleanSession == 0 && found.cleanSession == 0 {
			// Resume last session   // fix me ssddn
			this.storage.UpdateSession(this)
			// Notify other mqtt node to release resource
			core.AsyncProduceMessage(this.config,
				"session",
				TopicNameSession,
				&SessionTopic{
					Launcher:  this.conn.LocalAddr().String(),
					SessionId: clientid,
					Action:    ObjectActionUpdate,
					State:     mqttStateDisconnecting,
				})
		}

	} else {
		// Register the session in storage
		this.storage.RegisterSession(this)
	}

	if willMsg != nil {
		this.willMsg = willMsg
		this.willMsg.topic = willTopic
		if len(payload) > 0 {
			this.willMsg.payload = payload
		} else {
			this.willMsg.payload = nil
		}
		this.willMsg.qos = willQos
		this.willMsg.retain = willRetain
	}
	this.clientID = clientid
	this.cleanSession = cleanSession
	this.pingTime = nil
	this.isDroping = false

	// Remove any queued messages that are no longer allowd through ACL
	// Assuming a possible change of username
	this.storage.DeleteMessageWithValidator(
		clientid,
		func(msg StorageMessage) bool {
			err := this.authapi.CheckAcl(context.Background(), clientid, username, willTopic, auth.AclActionRead)
			if err != nil {
				return false
			}
			return true
		})

	this.state = mqttStateConnected
	err = this.sendConnAck(uint8(conack), CONNACK_ACCEPTED)
	return err
}

// handleDisconnect handle disconnect packet
func (this *mqttSession) handleDisconnect() error {
	glog.Infof("Received DISCONNECT from %s", this.id)

	if this.inpacket.remainingLength != 0 {
		return mqttErrorInvalidProtocol
	}
	if this.protocol == mqttProtocol311 && (this.inpacket.command&0x0F) != 0x00 {
		this.disconnect()
		return mqttErrorInvalidProtocol
	}
	this.state = mqttStateDisconnecting
	this.disconnect()
	return nil
}

// disconnect will disconnect current connection because of protocol error
func (this *mqttSession) disconnect() {
	if this.state == mqttStateDisconnected {
		return
	}
	if this.cleanSession > 0 {
		this.storage.DeleteSession(this.id)
		this.id = ""
	}
	this.state = mqttStateDisconnected
	this.conn.Close()
	this.conn = nil
}

// handleSubscribe handle subscribe packet
func (this *mqttSession) handleSubscribe() error {
	payload := make([]uint8, 0)

	glog.Infof("Received SUBSCRIBE from %s", this.id)
	if this.protocol == mqttProtocol311 {
		if (this.inpacket.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	mid, err := this.inpacket.readUint16()
	if err != nil {
		return err
	}
	// Deal each subscription
	for this.inpacket.pos < this.inpacket.remainingLength {
		sub := ""
		qos := uint8(0)
		if sub, err = this.inpacket.readString(); err != nil {
			return err
		}
		if checkTopicValidity(sub) != nil {
			glog.Errorf("Invalid subscription topic %s from %s, disconnecting", sub, this.id)
			return mqttErrorInvalidProtocol
		}
		if qos, err = this.inpacket.readByte(); err != nil {
			return err
		}

		if qos > 2 {
			glog.Errorf("Invalid Qos in subscription %s from %s", sub, this.id)
			return mqttErrorInvalidProtocol
		}

		if this.observer != nil {
			mp := this.observer.OnGetMountPoint()
			sub = mp + sub
		}
		if qos != 0x80 {
			if err := this.storage.AddSubscription(this.id, sub, qos); err != nil {
				return err
			}
			if err := this.storage.RetainSubscription(this.id, sub, qos); err != nil {
				return err
			}
		}
		payload = append(payload, qos)
	}

	if this.protocol == mqttProtocol311 && len(payload) == 0 {
		return mqttErrorInvalidProtocol
	}
	return this.sendSubAck(mid, payload)
}

// handleUnsubscribe handle unsubscribe packet
func (this *mqttSession) handleUnsubscribe() error {
	glog.Infof("Received UNSUBSCRIBE from %s", this.id)

	if this.protocol == mqttProtocol311 && (this.inpacket.command&0x0f) != 0x02 {
		return mqttErrorInvalidProtocol
	}
	mid, err := this.inpacket.readUint16()
	if err != nil {
		return err
	}
	// Iterate all subscription
	for this.inpacket.pos < this.inpacket.remainingLength {
		sub, err := this.inpacket.readString()
		if err != nil {
			return mqttErrorInvalidProtocol
		}
		if err := checkTopicValidity(sub); err != nil {
			return fmt.Errorf("Invalid unsubscription string from %s, disconnecting", this.id)
		}
		this.storage.RemoveSubscription(this.id, sub)
	}

	return this.sendCommandWithMid(UNSUBACK, mid, false)
}

// handlePublish handle publish packet
func (this *mqttSession) handlePublish() error {
	glog.Infof("Received PUBLISH from %s", this.id)

	var topic string
	var mid uint16
	var err error
	var payload []uint8

	dup := (this.inpacket.command & 0x08) >> 3
	qos := (this.inpacket.command & 0x06) >> 1
	if qos == 3 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing", this.id)
	}
	retain := (this.inpacket.command & 0x01)

	// Topic
	if topic, err = this.inpacket.readString(); err != nil {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", this.id)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, this.id)
	}
	if this.observer != nil && this.observer.OnGetMountPoint() != "" {
		topic = this.observer.OnGetMountPoint() + topic
	}

	if qos > 0 {
		mid, err = this.inpacket.readUint16()
		if err != nil {
			return err
		}
	}

	// Payload
	payloadlen := this.inpacket.remainingLength - this.inpacket.pos
	if payloadlen > 0 {
		limitSize, _ := this.config.Int(this.mgr.protocol, "message_size_limit")
		if payloadlen > limitSize {
			return mqttErrorInvalidProtocol
		}
		payload, err = this.inpacket.readBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	if this.observer != nil {
		err := this.authapi.CheckAcl(context.Background(), this.id, this.username, topic, auth.AclActionWrite)
		switch err {
		case auth.ErrorAclDenied:
			return mqttErrorInvalidProtocol
		default:
			return err
		}
	}
	glog.Infof("Received PUBLISH from %s(d:%d, q:%d r:%d, m:%d, '%s',..(%d)bytes",
		this.id, dup, qos, retain, mid, topic, payloadlen)

	// Check wether the message has been stored
	dup = 0
	var storedMsg *mqttMessage
	if qos > 0 {
		for _, storedMsg = range this.storedMsgs {
			if storedMsg.mid == mid && storedMsg.direction == mqttMessageDirectionIn {
				dup = 1
				break
			}
		}
	}
	msg := StorageMessage{
		ID:        uint(mid),
		SourceID:  this.id,
		Topic:     topic,
		Direction: MessageDirectionIn,
		State:     0,
		Qos:       qos,
		Retain:    (retain > 0),
		Payload:   payload,
	}

	switch qos {
	case 0:
		err = this.storage.QueueMessage(this.id, msg)
	case 1:
		err = this.storage.QueueMessage(this.id, msg)
		err = this.sendPubAck(mid)
	case 2:
		err = nil
		if dup > 0 {
			err = this.storage.InsertMessage(this.id, mid, MessageDirectionIn, msg)
		}
		if err == nil {
			err = this.sendPubRec(mid)
		}
	default:
		err = mqttErrorInvalidProtocol
	}

	return err
}

// handlePubRel handle pubrel packet
func (this *mqttSession) handlePubRel() error {
	// Check protocol specifal requriement
	if this.protocol == mqttProtocol311 {
		if (this.inpacket.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	mid, err := this.inpacket.readUint16()
	if err != nil {
		return err
	}

	this.storage.DeleteMessage(this.id, mid, MessageDirectionIn)
	return this.sendPubComp(mid)
}

// sendSimpleCommand send a simple command
func (this *mqttSession) sendSimpleCommand(cmd uint8) error {
	p := &mqttPacket{
		command:        cmd,
		remainingCount: 0,
	}
	return this.queuePacket(p)
}

// sendPingRsp send ping response to client
func (this *mqttSession) sendPingRsp() error {
	glog.Infof("Sending PINGRESP to %s", this.id)
	return this.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (this *mqttSession) sendConnAck(ack uint8, result uint8) error {
	glog.Infof("Sending CONNACK from %s", this.id)

	packet := &mqttPacket{
		command:         CONNACK,
		remainingLength: 2,
	}
	packet.initializePacket()
	packet.payload[packet.pos+0] = ack
	packet.payload[packet.pos+1] = result

	return this.queuePacket(packet)
}

// sendSubAck send subscription acknowledge to client
func (this *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	glog.Infof("Sending SUBACK on %s", this.id)
	packet := &mqttPacket{
		command:         SUBACK,
		remainingLength: 2 + int(len(payload)),
	}

	packet.initializePacket()
	packet.writeUint16(mid)
	if len(payload) > 0 {
		packet.writeBytes(payload)
	}
	return this.queuePacket(packet)
}

// sendCommandWithMid send command with message identifier
func (this *mqttSession) sendCommandWithMid(command uint8, mid uint16, dup bool) error {
	packet := &mqttPacket{
		command:         command,
		remainingLength: 2,
	}
	if dup {
		packet.command |= 8
	}
	packet.initializePacket()
	packet.payload[packet.pos+0] = uint8((mid & 0xFF00) >> 8)
	packet.payload[packet.pos+1] = uint8(mid & 0xff)
	return this.queuePacket(packet)
}

// sendPubAck
func (this *mqttSession) sendPubAck(mid uint16) error {
	glog.Info("Sending PUBACK to %s with MID:%d", this.id, mid)
	return this.sendCommandWithMid(PUBACK, mid, false)
}

// sendPubRec
func (this *mqttSession) sendPubRec(mid uint16) error {
	glog.Info("Sending PUBRREC to %s with MID:%d", this.id, mid)
	return this.sendCommandWithMid(PUBREC, mid, false)
}

func (this *mqttSession) sendPubComp(mid uint16) error {
	glog.Info("Sending PUBCOMP to %s with MID:%d", this.id, mid)
	return this.sendCommandWithMid(PUBCOMP, mid, false)
}

func (this *mqttSession) queuePacket(p *mqttPacket) error {
	p.pos = 0
	p.toprocess = p.length
	this.sendPacketChannel <- p
	return nil
}

func (this *mqttSession) QueueMessage(msg *mqttMessage) error {
	this.sendMsgChannel <- msg
	return nil
}

func (this *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}

func (this *mqttSession) generateMid() uint16 {
	return 0
}

func (this *mqttSession) sendPublish(subQos uint8, srcQos uint8, topic string, payload []uint8) error {
	/* Check for ACL topic accesthis. */
	// TODO

	var qos uint8
	if option, err := this.config.Bool(this.mgr.protocol, "upgrade_outgoing_qos"); err != nil && option {
		qos = subQos
	} else {
		if srcQos > subQos {
			qos = subQos
		} else {
			qos = srcQos
		}
	}

	var mid uint16
	if qos > 0 {
		mid = this.generateMid()
	} else {
		mid = 0
	}

	packet := newMqttPacket()
	packet.command = PUBLISH
	packet.remainingLength = 2 + len(topic) + len(payload)
	if qos > 0 {
		packet.remainingLength += 2
	}
	packet.initializePacket()
	packet.writeString(topic)
	if qos > 0 {
		packet.writeUint16(mid)
	}
	packet.writeBytes(payload)

	return this.queuePacket(&packet)
}
