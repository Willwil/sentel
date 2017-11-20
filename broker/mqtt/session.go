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

import (
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
		sendMsgChannel:    make(chan *mqttMessage, msgqsize),
		msgs:              make([]*mqttMessage, msgqsize),
		storedMsgs:        make([]*mqttMessage, msgqsize),
	}

	return s, nil
}

func (p *mqttSession) Identifier() string    { return p.id }
func (p *mqttSession) Service() core.Service { return p.mgr }
func (p *mqttSession) RegisterObserver(o base.SessionObserver) {
	if p.observer != nil {
		glog.Error("MqttSession register multiple observer")
	}
	p.observer = o
}
func (p *mqttSession) Info() *base.SessionInfo { return nil }

// launchPacketSendHandler launch goroutine to send packet queued for client
func (p *mqttSession) launchPacketSendHandler() {
	go func(stopChannel chan int, packetChannel chan *mqttPacket, msgChannel chan *mqttMessage) {
		defer p.waitgroup.Add(1)

		for {
			select {
			case <-stopChannel:
				return
			case packet := <-packetChannel:
				for packet.toprocess > 0 {
					len, err := p.conn.Write(packet.payload[packet.pos:packet.toprocess])
					if err != nil {
						glog.Fatal("Failed to send packet to '%s:%s'", p.id, err)
						return
					}
					if len > 0 {
						packet.toprocess -= len
						packet.pos += len
					} else {
						glog.Fatal("Failed to send packet to '%s'", p.id)
						return
					}
				}
			case msg := <-msgChannel:
				p.processMessage(msg)
			case <-time.After(1 * time.Second):
				p.processTimeout()
			}
		}
	}(p.sendStopChannel, p.sendPacketChannel, p.sendMsgChannel)
}

// processMessage proceess messages
func (p *mqttSession) processMessage(msg *mqttMessage) error {

	return nil
}

// processTimeout proceess timeout
func (p *mqttSession) processTimeout() error {

	return nil
}

// Handle is mainprocessor for iot device client
func (p *mqttSession) Handle() error {

	glog.Infof("Handling session:%s", p.id)
	defer p.Destroy()

	p.launchPacketSendHandler()
	for {
		var err error
		if err = p.inpacket.DecodeFromReader(p.conn, base.NilDecodeFeedback{}); err != nil {
			glog.Error(err)
			return err
		}
		switch p.inpacket.command & 0xF0 {
		case PINGREQ:
			err = p.handlePingReq()
		case CONNECT:
			err = p.handleConnect()
		case DISCONNECT:
			err = p.handleDisconnect()
		case PUBLISH:
			err = p.handlePublish()
		case PUBREL:
			err = p.handlePubRel()
		case SUBSCRIBE:
			err = p.handleSubscribe()
		case UNSUBSCRIBE:
			err = p.handleUnsubscribe()
		default:
			err = fmt.Errorf("Unrecognized protocol command:%d", int(p.inpacket.command&0xF0))
		}
		if err != nil {
			glog.Error(err)
			return err
		}
		// Check sesstion state
		if p.state == mqttStateDisconnected {
			break
		}
		p.inpacket.Clear()
	}
	return nil
}

// Destroy will destory the current session
func (p *mqttSession) Destroy() error {
	// Stop packet sender goroutine
	p.sendStopChannel <- 1
	p.waitgroup.Wait()
	if p.conn != nil {
		p.conn.Close()
	}
	p.mgr.removeSession(p)
	return nil
}

// generateId generate id fro session or client
func (p *mqttSession) generateId() string {
	return uuid.NewV4().String()
}

// handlePingReq handle ping request packet
func (p *mqttSession) handlePingReq() error {
	glog.Infof("Received PINGREQ from %s", p.Identifier())
	return p.sendPingRsp()
}

// handleConnect handle connect packet
func (p *mqttSession) handleConnect() error {
	glog.Infof("Handling CONNECT packet from %s", p.id)

	if p.state != mqttStateNew {
		return errors.New("Invalid session state")
	}
	// Check protocol name and version
	protocolName, err := p.inpacket.readString()
	if err != nil {
		return err
	}
	protocolVersion, err := p.inpacket.readByte()
	if err != nil {
		return err
	}
	switch protocolName {
	case PROTOCOL_NAME_V31:
		if protocolVersion&0x7F != PROTOCOL_VERSION_V31 {
			p.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		p.protocol = mqttProtocol311

	case PROTOCOL_NAME_V311:
		if protocolVersion&0x7F != PROTOCOL_VERSION_V311 {
			p.sendConnAck(0, CONNACK_REFUSED_PROTOCOL_VERSION)
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		// Reserved flags is not set to 0, must disconnect
		if p.inpacket.command&0x0F != 0x00 {
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		p.protocol = mqttProtocol311
	default:
		return fmt.Errorf("Invalid protocol name '%s' in CONNECT packet", protocolName)
	}

	// Check connect flags
	cflags, err := p.inpacket.readByte()
	if err != nil {
		return nil
	}
	/*
		if p.mgr.protocol == mqttProtocol311 {
			if cflags&0x01 != 0x00 {
				return errorp.New("Invalid protocol version in connect flags")
			}
		}
	*/
	cleanSession := (cflags & 0x02) >> 1
	will := cflags & 0x04
	willQos := (cflags & 0x18) >> 3
	if willQos >= 3 { // qos level3 is not supported
		return fmt.Errorf("Invalid Will Qos in CONNECT from %s", p.id)
	}

	willRetain := (cflags & 0x20) == 0x20
	passwordFlag := cflags & 0x40
	usernameFlag := cflags & 0x80
	keepalive, err := p.inpacket.readUint16()
	if err != nil {
		return err
	}
	p.keepalive = keepalive

	// Deal with client identifier
	clientid, err := p.inpacket.readString()
	if err != nil {
		return err
	}
	if clientid == "" {
		if p.protocol == mqttProtocol31 {
			p.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
		} else {
			if option, err := p.config.Bool(p.mgr.protocol, "allow_zero_length_clientid"); err != nil && (!option || cleanSession == 0) {
				p.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
				return errors.New("Invalid mqtt packet with client id")
			}

			clientid = p.generateId()
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
		topic, err := p.inpacket.readString()
		if err != nil || topic == "" {
			return nil
		}
		willTopic = topic
		if p.observer != nil {
			willTopic = p.observer.OnGetMountPoint() + topic
		}
		if err := checkTopicValidity(willTopic); err != nil {
			return err
		}
		// Get willtopic's payload
		willPayloadLength, err := p.inpacket.readUint16()
		if err != nil {
			return err
		}
		if willPayloadLength > 0 {
			payload, err = p.inpacket.readBytes(int(willPayloadLength))
			if err != nil {
				return err
			}
		}
	} else {
		if p.protocol == mqttProtocol311 {
			if willQos != 0 || willRetain {
				return mqttErrorInvalidProtocol
			}
		}
	} // else will

	var username string
	var password string
	if usernameFlag > 0 {
		username, err = p.inpacket.readString()
		if err == nil {
			if passwordFlag > 0 {
				password, err = p.inpacket.readString()
				if err == mqttErrorInvalidProtocol {
					if p.protocol == mqttProtocol31 {
						passwordFlag = 0
					} else if p.protocol == mqttProtocol311 {
						return err
					} else {
						return err
					}
				}
			}
		} else {
			if p.protocol == mqttProtocol31 {
				usernameFlag = 0
			} else {
				return err
			}
		}
	} else { // username flag
		if p.protocol == mqttProtocol311 {
			if passwordFlag > 0 {
				return mqttErrorInvalidProtocol
			}
		}
	}

	if usernameFlag > 0 {
		if p.observer != nil {
			err := p.observer.OnAuthenticate(p, username, password)
			switch err {
			case nil:
			case base.IotErrorAuthFailed:
				p.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
				p.disconnect()
				return err
			default:
				p.disconnect()
				return err

			}
			// Get username and passowrd sucessfuly
			p.username = username
			p.password = password
		}
		// Get anonymous allow configuration
		allowAnonymous, _ := p.config.Bool(p.mgr.protocol, "allow_anonymous")
		if usernameFlag > 0 && allowAnonymous == false {
			// Dont allow anonymous client connection
			p.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	// Check wether username will be used as client id,
	// The connection request will be refused if the option is set
	if option, err := p.config.Bool(p.mgr.protocol, "user_name_as_client_id"); err != nil && option {
		if p.username != "" {
			clientid = p.username
		} else {
			p.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	conack := 0
	// Find if the client already has an entry, p must be done after any security check
	if found, _ := p.storage.FindSession(clientid); found != nil {
		// Found old session
		if found.state == mqttStateInvalid {
			glog.Errorf("Invalid session(%s) in store", found.id)
		}
		if p.protocol == mqttProtocol311 {
			if cleanSession == 0 {
				conack |= 0x01
			}
		}
		p.cleanSession = cleanSession

		if p.cleanSession == 0 && found.cleanSession == 0 {
			// Resume last session   // fix me ssddn
			p.storage.UpdateSession(p)
			// Notify other mqtt node to release resource
			core.AsyncProduceMessage(p.config,
				"session",
				TopicNameSession,
				&SessionTopic{
					TopicBase: core.TopicBase{Action: core.TopicActionUpdate},
					Launcher:  p.conn.LocalAddr().String(),
					SessionId: clientid,
					State:     mqttStateDisconnecting,
				})
		}

	} else {
		// Register the session in storage
		p.storage.RegisterSession(p)
	}

	if willMsg != nil {
		p.willMsg = willMsg
		p.willMsg.topic = willTopic
		if len(payload) > 0 {
			p.willMsg.payload = payload
		} else {
			p.willMsg.payload = nil
		}
		p.willMsg.qos = willQos
		p.willMsg.retain = willRetain
	}
	p.clientID = clientid
	p.cleanSession = cleanSession
	p.pingTime = nil
	p.isDroping = false

	// Remove any queued messages that are no longer allowd through ACL
	// Assuming a possible change of username
	p.storage.DeleteMessageWithValidator(
		clientid,
		func(msg StorageMessage) bool {
			err := auth.CheckAcl(clientid, username, willTopic, auth.AclActionRead)
			if err != nil {
				return false
			}
			return true
		})

	p.state = mqttStateConnected
	err = p.sendConnAck(uint8(conack), CONNACK_ACCEPTED)
	return err
}

// handleDisconnect handle disconnect packet
func (p *mqttSession) handleDisconnect() error {
	glog.Infof("Received DISCONNECT from %s", p.id)

	if p.inpacket.remainingLength != 0 {
		return mqttErrorInvalidProtocol
	}
	if p.protocol == mqttProtocol311 && (p.inpacket.command&0x0F) != 0x00 {
		p.disconnect()
		return mqttErrorInvalidProtocol
	}
	p.state = mqttStateDisconnecting
	p.disconnect()
	return nil
}

// disconnect will disconnect current connection because of protocol error
func (p *mqttSession) disconnect() {
	if p.state == mqttStateDisconnected {
		return
	}
	if p.cleanSession > 0 {
		p.storage.DeleteSession(p.id)
		p.id = ""
	}
	p.state = mqttStateDisconnected
	p.conn.Close()
	p.conn = nil
}

// handleSubscribe handle subscribe packet
func (p *mqttSession) handleSubscribe() error {
	payload := make([]uint8, 0)

	glog.Infof("Received SUBSCRIBE from %s", p.id)
	if p.protocol == mqttProtocol311 {
		if (p.inpacket.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	mid, err := p.inpacket.readUint16()
	if err != nil {
		return err
	}
	// Deal each subscription
	for p.inpacket.pos < p.inpacket.remainingLength {
		sub := ""
		qos := uint8(0)
		if sub, err = p.inpacket.readString(); err != nil {
			return err
		}
		if checkTopicValidity(sub) != nil {
			glog.Errorf("Invalid subscription topic %s from %s, disconnecting", sub, p.id)
			return mqttErrorInvalidProtocol
		}
		if qos, err = p.inpacket.readByte(); err != nil {
			return err
		}

		if qos > 2 {
			glog.Errorf("Invalid Qos in subscription %s from %s", sub, p.id)
			return mqttErrorInvalidProtocol
		}

		if p.observer != nil {
			mp := p.observer.OnGetMountPoint()
			sub = mp + sub
		}
		if qos != 0x80 {
			if err := p.storage.AddSubscription(p.id, sub, qos); err != nil {
				return err
			}
			if err := p.storage.RetainSubscription(p.id, sub, qos); err != nil {
				return err
			}
		}
		payload = append(payload, qos)
	}

	if p.protocol == mqttProtocol311 && len(payload) == 0 {
		return mqttErrorInvalidProtocol
	}
	return p.sendSubAck(mid, payload)
}

// handleUnsubscribe handle unsubscribe packet
func (p *mqttSession) handleUnsubscribe() error {
	glog.Infof("Received UNSUBSCRIBE from %s", p.id)

	if p.protocol == mqttProtocol311 && (p.inpacket.command&0x0f) != 0x02 {
		return mqttErrorInvalidProtocol
	}
	mid, err := p.inpacket.readUint16()
	if err != nil {
		return err
	}
	// Iterate all subscription
	for p.inpacket.pos < p.inpacket.remainingLength {
		sub, err := p.inpacket.readString()
		if err != nil {
			return mqttErrorInvalidProtocol
		}
		if err := checkTopicValidity(sub); err != nil {
			return fmt.Errorf("Invalid unsubscription string from %s, disconnecting", p.id)
		}
		p.storage.RemoveSubscription(p.id, sub)
	}

	return p.sendCommandWithMid(UNSUBACK, mid, false)
}

// handlePublish handle publish packet
func (p *mqttSession) handlePublish() error {
	glog.Infof("Received PUBLISH from %s", p.id)

	var topic string
	var mid uint16
	var err error
	var payload []uint8

	dup := (p.inpacket.command & 0x08) >> 3
	qos := (p.inpacket.command & 0x06) >> 1
	if qos == 3 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing", p.id)
	}
	retain := (p.inpacket.command & 0x01)

	// Topic
	if topic, err = p.inpacket.readString(); err != nil {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", p.id)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, p.id)
	}
	if p.observer != nil && p.observer.OnGetMountPoint() != "" {
		topic = p.observer.OnGetMountPoint() + topic
	}

	if qos > 0 {
		mid, err = p.inpacket.readUint16()
		if err != nil {
			return err
		}
	}

	// Payload
	payloadlen := p.inpacket.remainingLength - p.inpacket.pos
	if payloadlen > 0 {
		limitSize, _ := p.config.Int(p.mgr.protocol, "message_size_limit")
		if payloadlen > limitSize {
			return mqttErrorInvalidProtocol
		}
		payload, err = p.inpacket.readBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	if p.observer != nil {
		err := auth.CheckAcl(p.id, p.username, topic, auth.AclActionWrite)
		switch err {
		case auth.ErrorAclDenied:
			return mqttErrorInvalidProtocol
		default:
			return err
		}
	}
	glog.Infof("Received PUBLISH from %s(d:%d, q:%d r:%d, m:%d, '%s',..(%d)bytes",
		p.id, dup, qos, retain, mid, topic, payloadlen)

	// Check wether the message has been stored
	dup = 0
	var storedMsg *mqttMessage
	if qos > 0 {
		for _, storedMsg = range p.storedMsgs {
			if storedMsg.mid == mid && storedMsg.direction == mqttMessageDirectionIn {
				dup = 1
				break
			}
		}
	}
	msg := StorageMessage{
		ID:        uint(mid),
		SourceID:  p.id,
		Topic:     topic,
		Direction: MessageDirectionIn,
		State:     0,
		Qos:       qos,
		Retain:    (retain > 0),
		Payload:   payload,
	}

	switch qos {
	case 0:
		err = p.storage.QueueMessage(p.id, msg)
	case 1:
		err = p.storage.QueueMessage(p.id, msg)
		err = p.sendPubAck(mid)
	case 2:
		err = nil
		if dup > 0 {
			err = p.storage.InsertMessage(p.id, mid, MessageDirectionIn, msg)
		}
		if err == nil {
			err = p.sendPubRec(mid)
		}
	default:
		err = mqttErrorInvalidProtocol
	}

	return err
}

// handlePubRel handle pubrel packet
func (p *mqttSession) handlePubRel() error {
	// Check protocol specifal requriement
	if p.protocol == mqttProtocol311 {
		if (p.inpacket.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	mid, err := p.inpacket.readUint16()
	if err != nil {
		return err
	}

	p.storage.DeleteMessage(p.id, mid, MessageDirectionIn)
	return p.sendPubComp(mid)
}

// sendSimpleCommand send a simple command
func (p *mqttSession) sendSimpleCommand(cmd uint8) error {
	packet := &mqttPacket{
		command:        cmd,
		remainingCount: 0,
	}
	return p.queuePacket(packet)
}

// sendPingRsp send ping response to client
func (p *mqttSession) sendPingRsp() error {
	glog.Infof("Sending PINGRESP to %s", p.id)
	return p.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (p *mqttSession) sendConnAck(ack uint8, result uint8) error {
	glog.Infof("Sending CONNACK from %s", p.id)

	packet := &mqttPacket{
		command:         CONNACK,
		remainingLength: 2,
	}
	packet.initializePacket()
	packet.payload[packet.pos+0] = ack
	packet.payload[packet.pos+1] = result

	return p.queuePacket(packet)
}

// sendSubAck send subscription acknowledge to client
func (p *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	glog.Infof("Sending SUBACK on %s", p.id)
	packet := &mqttPacket{
		command:         SUBACK,
		remainingLength: 2 + int(len(payload)),
	}

	packet.initializePacket()
	packet.writeUint16(mid)
	if len(payload) > 0 {
		packet.writeBytes(payload)
	}
	return p.queuePacket(packet)
}

// sendCommandWithMid send command with message identifier
func (p *mqttSession) sendCommandWithMid(command uint8, mid uint16, dup bool) error {
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
	return p.queuePacket(packet)
}

// sendPubAck
func (p *mqttSession) sendPubAck(mid uint16) error {
	glog.Info("Sending PUBACK to %s with MID:%d", p.id, mid)
	return p.sendCommandWithMid(PUBACK, mid, false)
}

// sendPubRec
func (p *mqttSession) sendPubRec(mid uint16) error {
	glog.Info("Sending PUBRREC to %s with MID:%d", p.id, mid)
	return p.sendCommandWithMid(PUBREC, mid, false)
}

func (p *mqttSession) sendPubComp(mid uint16) error {
	glog.Info("Sending PUBCOMP to %s with MID:%d", p.id, mid)
	return p.sendCommandWithMid(PUBCOMP, mid, false)
}

func (p *mqttSession) queuePacket(packet *mqttPacket) error {
	packet.pos = 0
	packet.toprocess = packet.length
	p.sendPacketChannel <- packet
	return nil
}

func (p *mqttSession) QueueMessage(msg *mqttMessage) error {
	p.sendMsgChannel <- msg
	return nil
}

func (p *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}

func (p *mqttSession) generateMid() uint16 {
	return 0
}

func (p *mqttSession) sendPublish(subQos uint8, srcQos uint8, topic string, payload []uint8) error {
	/* Check for ACL topic accesp. */
	// TODO

	var qos uint8
	if option, err := p.config.Bool(p.mgr.protocol, "upgrade_outgoing_qos"); err != nil && option {
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
		mid = p.generateMid()
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

	return p.queuePacket(&packet)
}
