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
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/queue"
	sm "github.com/cloustone/sentel/broker/sessionmgr"
	"github.com/cloustone/sentel/core"

	auth "github.com/cloustone/sentel/broker/auth"
	uuid "github.com/satori/go.uuid"

	"github.com/golang/glog"
)

// mqttSession handle mqtt protocol session
type mqttSession struct {
	config         core.Config        // Session Configuration
	conn           net.Conn           // Underlay network connection handle
	clientId       string             // Session's client identifier
	inpacket       mqttPacket         // Current inprocessing packet
	keepalive      uint16             // MQTT protocol kepalive time
	protocol       uint8              // MQTT protocol version
	username       string             // Session's user name and password
	password       string             // Session's password
	lastMessageIn  time.Time          // Last message's in time
	lastMessageOut time.Time          // last message out time
	bytesReceived  int64              // Byte recivied
	cleanSession   uint8              // Wether the session is clean session
	willMsg        *base.Message      // MQTT protocol will message
	waitgroup      sync.WaitGroup     // Waitgroup for all goroutine
	mountpoint     string             // Topic's mount point
	stateMutex     sync.Mutex         // Mutex to session state
	msgState       uint8              // Mqtt message state
	sessionState   uint8              // Mqtt session state
	authOptions    *auth.Options      // authentication options
	aliveTimer     *time.Timer        // Keepalive timer
	pingTime       *time.Time         // Ping timer
	availableChan  chan *base.Message // Event chanel for data availabel notification
	packetChan     chan *mqttPacket   // Mqtt packet channel
	errorChan      chan int           // Error channel
	queue          queue.Queue        // Session's queue
}

// newMqttSession create new session  for each client connection
func newMqttSession(m *mqttService, conn net.Conn) (*mqttSession, error) {
	// make a persudo alive timer for convenience, and stop it at first
	palivetimer := time.NewTimer(time.Duration(20 * time.Second))
	palivetimer.Stop()
	return &mqttSession{
		config:        m.Config,
		conn:          conn,
		bytesReceived: 0,
		sessionState:  mqttStateNew,
		inpacket:      newMqttPacket(),
		protocol:      mqttProtocolInvalid,
		mountpoint:    "",
		msgState:      mqttMsgStateInvalid,
		aliveTimer:    palivetimer,
		availableChan: make(chan *base.Message),
		packetChan:    make(chan *mqttPacket),
		errorChan:     make(chan int),
		queue:         nil,
	}, nil
}

// checkTopiValidity will check topic's validity
func checkTopicValidity(topic string) error {
	return nil
}

// setSessionState set session's state
func (p *mqttSession) setSessionState(state uint8) {
	old := p.sessionState
	p.sessionState = state
	glog.Infof("mqtt: client '%s' session state changed (%s->%s)", p.clientId,
		nameOfSessionState(old),
		nameOfSessionState(state))
}

// checkSessionState check wether the session state is same with expected state
func (p *mqttSession) checkSessionState(state uint8) error {
	if p.sessionState == state {
		return nil
	} else {
		return fmt.Errorf("mqtt: state miss occurred, expected:%s, now:%s",
			nameOfSessionState(state), nameOfSessionState(p.sessionState))
	}
}

// setMessageState set session's state
func (p *mqttSession) setMessageState(state uint8) {
	old := p.msgState
	p.msgState = state
	glog.Infof("mqtt: client '%s' message state changed (%s->%s)", p.clientId,
		nameOfMessageState(old),
		nameOfMessageState(state))
}

// Handle is mainprocessor for iot device client
func (p *mqttSession) Handle() error {
	glog.Infof("Handling session:%s", p.clientId)

	// Start goroutine to read packet from client
	go func(p *mqttSession) {
		if err := p.inpacket.DecodeFromReader(p.conn, base.NilDecodeFeedback{}); err != nil {
			glog.Error(err)
			p.errorChan <- 1
		}
	}(p)
	defer event.Notify(&event.Event{Type: event.SessionDestroy, ClientId: p.clientId})
	var err error
	for {
		select {
		case <-p.errorChan:
			err = errors.New("mqtt read error occurred")
		case msg := <-p.availableChan:
			err = p.handleDataAvailableNotification(msg)
		case <-p.aliveTimer.C:
			err = p.handleAliveTimeout()
		case <-p.packetChan:
			err = p.handleMqttInPacket()
		}
		if err != nil {
			glog.Error(err)
			return err
		}
		// Check sesstion state
		if err = p.checkSessionState(mqttStateDisconnected); err != nil {
			glog.Error(err)
			return err
		}
	}
	return err
}

// handleMqttInPacket handle all mqtt in packet
func (p *mqttSession) handleMqttInPacket() error {
	glog.Infof("%s,%s", nameOfSessionState(p.sessionState), nameOfMessageState(p.msgState))
	switch p.inpacket.command & 0xF0 {
	case PINGREQ:
		return p.handlePingReq()
	case CONNECT:
		return p.handleConnect()
	case DISCONNECT:
		return p.handleDisconnect()
	case PUBLISH:
		return p.handlePublish()
	case PUBREL:
		return p.handlePubRel()
	case SUBSCRIBE:
		return p.handleSubscribe()
	case UNSUBSCRIBE:
		return p.handleUnsubscribe()
	case PUBACK:
		return p.handlePubAck()
	default:
		return fmt.Errorf("Unrecognized protocol command:%d", int(p.inpacket.command&0xF0))
	}
	p.inpacket.Clear()
	return nil
}

// Destroy will destory the current session
func (p *mqttSession) Destroy() error {
	// Stop packet sender goroutine
	p.waitgroup.Wait()
	if p.conn != nil {
		p.conn.Close()
		p.conn = nil
	}
	if p.aliveTimer != nil {
		p.aliveTimer.Stop()
		p.aliveTimer = nil
	}
	return nil
}

// handleDataAvailableNotification read message from queue and send to client
func (p *mqttSession) handleDataAvailableNotification(s *base.Message) error {
	for {
		if msg := p.queue.Front(); msg != nil {
			packet := makePubPacketByMessage(msg)
			if err := p.writePacket(packet); err == nil {
				switch msg.Qos {
				case 0:
					// NO client response for qos0 pub packet
					p.queue.Pop()
					continue
				case 1:
					p.msgState = mqttMsgStateWaitPubAck
					break
				default:
				}
			}
		}
	}
	return nil
}

// handleAliveTimeout reply to client with ping rsp after timeout
func (p *mqttSession) handleAliveTimeout() error {
	return errors.New("mqtt: time out occured")
}

// DataAvailable called when queue have data to deal with
func (p *mqttSession) DataAvailable(q queue.Queue, msg *base.Message) {
	p.availableChan <- msg
}

// makePubPacketByMessage construct a mqtt packet by message
func makePubPacketByMessage(msg *base.Message) *mqttPacket {
	// Write data back to client
	length := len(msg.Topic)
	if msg.Qos > 0 {
		length += 2
	}
	packet := mqttPacket{
		command:         PUBLISH,
		qos:             msg.Qos,
		retain:          msg.Retain,
		dup:             msg.Dup,
		remainingLength: length,
	}
	packet.initializePacket()
	packet.writeString(msg.Topic)
	if msg.Qos > 0 {
		packet.writeUint16(msg.PacketId)
	}
	packet.writeBytes(msg.Payload)
	return &packet
}

func (p *mqttSession) Id() string            { return p.clientId }
func (p *mqttSession) BrokerId() string      { return base.GetBrokerId() }
func (p *mqttSession) Info() *sm.SessionInfo { return nil }
func (p *mqttSession) IsValid() bool         { return true }
func (p *mqttSession) IsPersistent() bool    { return (p.cleanSession == 0) }

// handlePingReq handle ping request packet
func (p *mqttSession) handlePingReq() error {
	glog.Infof("Received PINGREQ from %s", p.clientId)
	return p.sendPingRsp()
}

// handleConnect handle connect packet
func (p *mqttSession) handleConnect() error {
	glog.Infof("Handling CONNECT packet from %s", p.clientId)

	if err := p.checkSessionState(mqttStateNew); err != nil {
		return err
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
		p.protocol = mqttProtocol31

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
		return err
	}
	if p.protocol == mqttProtocol311 {
		if cflags&0x01 != 0x00 {
			return errors.New("Invalid protocol version in connect flags")
		}
	}
	cleanSession := (cflags & 0x02) >> 1
	will := cflags & 0x04
	willQos := (cflags & 0x18) >> 3
	if willQos >= 3 { // qos level3 is not supported
		return fmt.Errorf("Invalid Will Qos in CONNECT from %s", p.clientId)
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
	clientId, err := p.inpacket.readString()
	if err != nil {
		return err
	}
	if clientId == "" {
		if p.protocol == mqttProtocol31 {
			p.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
		} else {
			option, err := p.config.Bool("mqtt", "allow_zero_length_clientId")
			if err != nil && (!option || cleanSession == 0) {
				p.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
				return errors.New("Invalid mqtt packet with client id")
			}
			clientId = uuid.NewV4().String()
		}
	}

	// Deal with topc
	var willTopic string
	var willMsg *base.Message
	var payload []uint8

	if will > 0 {
		willMsg = new(base.Message)
		// Get topic
		topic, err := p.inpacket.readString()
		if err != nil || topic == "" {
			return nil
		}
		willTopic = p.mountpoint + topic
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
		// Parse mqtt client request options
		p.authOptions, err = p.parseRequestOptions(clientId, username, password)
		if err != nil {
			return err
		}

		err := auth.Authenticate(p.authOptions)
		switch err {
		case nil:
		case auth.ErrorAuthDenied:
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
		// Get anonymous allow configuration
		allowAnonymous, _ := p.config.Bool("mqtt", "allow_anonymous")
		if usernameFlag > 0 && allowAnonymous == false {
			// Dont allow anonymous client connection
			p.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	// Check wether username will be used as client id,
	// The connection request will be refused if the option is set
	if option, err := p.config.Bool("mqtt", "user_name_as_client_id"); err != nil && option {
		if p.username != "" {
			clientId = p.username
		} else {
			p.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			return mqttErrorInvalidProtocol
		}
	}
	conack := 0

	if cleanSession == 0 {
		// Find if the client already has an entry, p must be done after any security check
		if found, _ := sm.FindSession(clientId); found != nil {
			// Found old session
			if !found.IsValid() {
				glog.Errorf("Invalid session(%s) in store", found.Id)
			}
			info := found.Info()
			if p.protocol == mqttProtocol311 {
				if cleanSession == 0 {
					conack |= 0x01
				}
			}
			p.cleanSession = cleanSession
			if p.cleanSession == 0 && info.CleanSession == 0 {
				// Resume last session and notify other mqtt node to release resource
				event.Notify(&event.Event{
					Type:     event.SessionResume,
					ClientId: clientId,
				})

			}

		}
	}
	if willMsg != nil {
		p.willMsg = willMsg
		p.willMsg.Topic = willTopic
		if len(payload) > 0 {
			p.willMsg.Payload = payload
		} else {
			p.willMsg.Payload = nil
		}
		p.willMsg.Qos = willQos
		p.willMsg.Retain = willRetain
	}
	p.clientId = clientId
	p.cleanSession = cleanSession
	p.pingTime = nil
	// Remove any queued messages that are no longer allowd through ACL
	// Assuming a possible change of username
	sm.DeleteMessageWithValidator(
		clientId,
		func(msg *base.Message) bool {
			err := auth.Authorize(clientId, username, willTopic, auth.AclRead, nil)
			return err != nil
		})

	// Change session state and reply client
	p.setSessionState(mqttStateConnected)
	err = p.sendConnAck(uint8(conack), CONNACK_ACCEPTED)

	// Create queue for this sesion and otify event service that new session created
	queue, err := queue.NewQueue(clientId, (cleanSession == 0), p)
	if err != nil {
		glog.Fatalf("broker: Failed to create queue for mqtt client '%s'", clientId)
		return err
	}
	p.queue = queue
	sm.RegisterSession(p)
	event.Notify(&event.Event{
		Type:     event.SessionCreate,
		ClientId: clientId,
		Detail:   &event.SessionCreateDetail{Persistent: (cleanSession == 0)}})

	// reset and start keep alive timer
	p.aliveTimer.Reset(time.Second * time.Duration(p.keepalive))
	return nil
}

// handleDisconnect handle disconnect packet
func (p *mqttSession) handleDisconnect() error {
	glog.Infof("Received DISCONNECT from %s", p.clientId)

	if p.inpacket.remainingLength != 0 {
		return mqttErrorInvalidProtocol
	}
	if p.protocol == mqttProtocol311 && (p.inpacket.command&0x0F) != 0x00 {
		p.disconnect()
		return mqttErrorInvalidProtocol
	}
	p.setSessionState(mqttStateDisconnecting)
	p.disconnect()
	return nil
}

// disconnect will disconnect current connection because of protocol error
func (p *mqttSession) disconnect() {
	if err := p.checkSessionState(mqttStateDisconnected); err == nil {
		if p.cleanSession > 0 {
			event.Notify(&event.Event{Type: event.SessionDestroy, ClientId: p.clientId})
			p.clientId = ""
		}
		p.setSessionState(mqttStateDisconnected)
		if p.aliveTimer != nil {
			p.aliveTimer.Stop()
		}
		p.conn.Close()
		p.conn = nil
	}
}

// handleSubscribe handle subscribe packet
func (p *mqttSession) handleSubscribe() error {
	payload := make([]uint8, 0)

	glog.Infof("Received SUBSCRIBE from %s", p.clientId)
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
			glog.Errorf("Invalid subscription topic %s from %s, disconnecting", sub, p.clientId)
			return mqttErrorInvalidProtocol
		}
		if qos, err = p.inpacket.readByte(); err != nil {
			return err
		}

		if qos > 2 {
			glog.Errorf("Invalid Qos in subscription %s from %s", sub, p.clientId)
			return mqttErrorInvalidProtocol
		}

		sub = p.mountpoint + sub
		if qos != 0x80 {
			event.Notify(&event.Event{Type: event.TopicSubscribe,
				ClientId: p.clientId,
				Detail:   event.TopicSubscribeDetail{Topic: sub, Qos: qos, Retain: true}})
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
	glog.Infof("Received UNSUBSCRIBE from %s", p.clientId)

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
			return fmt.Errorf("Invalid unsubscription string from %s, disconnecting", p.clientId)
		}
		event.Notify(&event.Event{Type: event.TopicUnsubscribe, ClientId: p.clientId, Detail: event.TopicUnsubscribeDetail{Topic: sub}})
	}

	return p.sendCommandWithMid(UNSUBACK, mid, false)
}

// handlePublish handle publish packet
func (p *mqttSession) handlePublish() error {
	glog.Infof("Received PUBLISH from %s", p.clientId)

	var topic string
	var mid uint16
	var err error
	var payload []uint8

	dup := (p.inpacket.command & 0x08) >> 3
	qos := (p.inpacket.command & 0x06) >> 1
	if qos == 3 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing", p.clientId)
	}
	retain := (p.inpacket.command & 0x01)

	// Topic
	if topic, err = p.inpacket.readString(); err != nil {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", p.clientId)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, p.clientId)
	}
	topic = p.mountpoint + topic

	if qos > 0 {
		mid, err = p.inpacket.readUint16()
		if err != nil {
			return err
		}
	}

	// Payload
	payloadlen := p.inpacket.remainingLength - p.inpacket.pos
	if payloadlen > 0 {
		limitSize, _ := p.config.Int("mqtt", "message_size_limit")
		if payloadlen > limitSize {
			return mqttErrorInvalidProtocol
		}
		payload, err = p.inpacket.readBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	err = auth.Authorize(p.clientId, p.username, topic, auth.AclWrite, nil)
	switch err {
	case auth.ErrorAuthDenied:
		return mqttErrorInvalidProtocol
	default:
		return err
	}
	glog.Infof("Received PUBLISH from %s(d:%d, q:%d r:%d, m:%d, '%s',..(%d)bytes",
		p.clientId, dup, qos, retain, mid, topic, payloadlen)

	// Check wether the message has been stored :TODO
	dup = 0
	if qos > 0 && sm.FindMessage(p.clientId, mid, 1) != nil {
		dup = 1
	}
	detail := event.TopicPublishDetail{
		Topic:   topic,
		Qos:     qos,
		Retain:  (retain > 0),
		Payload: payload,
		Dup:     dup == 1,
	}

	switch qos {
	case 0:
		event.Notify(&event.Event{Type: event.TopicPublish, ClientId: p.clientId, Detail: detail})
	case 1:
		event.Notify(&event.Event{Type: event.TopicPublish, ClientId: p.clientId, Detail: detail})
		err = p.sendPubAck(mid) // TODO
	case 2:
		err = errors.New("MQTT qos 2 is not supported now")
	default:
		err = mqttErrorInvalidProtocol
	}

	return err
}

// handlePubAck handle pubrel packet
func (p *mqttSession) handlePubAck() error {
	if p.msgState == mqttMsgStateWaitPubAck {
		q := queue.GetQueue(p.clientId)
		q.Pop() // TODO: we should check packet identifier
	}
	return nil
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

	//sm.DeleteMessage(p.clientId, mid, sm.MessageDirectionIn) // TODO:Qos2
	return p.sendPubComp(mid)
}

// sendSimpleCommand send a simple command
func (p *mqttSession) sendSimpleCommand(cmd uint8) error {
	packet := &mqttPacket{
		command:        cmd,
		remainingCount: 0,
	}
	return p.writePacket(packet)
}

// sendPingRsp send ping response to client
func (p *mqttSession) sendPingRsp() error {
	glog.Infof("Sending PINGRESP to %s", p.clientId)
	return p.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (p *mqttSession) sendConnAck(ack uint8, result uint8) error {
	glog.Infof("Sending CONNACK from %s", p.clientId)

	packet := &mqttPacket{
		command:         CONNACK,
		remainingLength: 2,
	}
	packet.initializePacket()
	packet.payload[packet.pos+0] = ack
	packet.payload[packet.pos+1] = result

	return p.writePacket(packet)
}

// sendSubAck send subscription acknowledge to client
func (p *mqttSession) sendSubAck(mid uint16, payload []uint8) error {
	glog.Infof("Sending SUBACK on %s", p.clientId)
	packet := &mqttPacket{
		command:         SUBACK,
		remainingLength: 2 + int(len(payload)),
	}

	packet.initializePacket()
	packet.writeUint16(mid)
	if len(payload) > 0 {
		packet.writeBytes(payload)
	}
	return p.writePacket(packet)
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
	return p.writePacket(packet)
}

// sendPubAck
func (p *mqttSession) sendPubAck(mid uint16) error {
	glog.Infof("Sending PUBACK to %s with MID:%d", p.clientId, mid)
	return p.sendCommandWithMid(PUBACK, mid, false)
}

// sendPubRec
func (p *mqttSession) sendPubRec(mid uint16) error {
	glog.Infof("Sending PUBRREC to %s with MID:%d", p.clientId, mid)
	return p.sendCommandWithMid(PUBREC, mid, false)
}

func (p *mqttSession) sendPubComp(mid uint16) error {
	glog.Infof("Sending PUBCOMP to %s with MID:%d", p.clientId, mid)
	return p.sendCommandWithMid(PUBCOMP, mid, false)
}

func (p *mqttSession) updateOutMessage(mid uint16, state int) error {
	return nil
}

func (p *mqttSession) generateMid() uint16 {
	// TODO: generate message ID
	return 0
}

func (p *mqttSession) sendPublish(subQos uint8, srcQos uint8, topic string, payload []uint8) error {
	/* Check for ACL topic accesp. */
	// TODO

	var qos uint8
	if option, err := p.config.Bool("mqtt", "upgrade_outgoing_qos"); err != nil && option {
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

	return p.writePacket(&packet)
}

func (p *mqttSession) writePacket(packet *mqttPacket) error {
	packet.pos = 0
	packet.toprocess = packet.length
	for packet.toprocess > 0 {
		len, err := p.conn.Write(packet.payload[packet.pos:packet.toprocess])
		if err != nil {
			return err
		}
		if len > 0 {
			packet.toprocess -= len
			packet.pos += len
		} else {
			return fmt.Errorf("Failed to send packet to '%s'", p.clientId)
		}
	}
	return nil
}

// getAuthOptions return authentication options from user's request
// mqttClientId:clientId +"|securemode=3,signmethod=hmacsha1,timestampe=xxxxx|"
// mqttUserName:deviceName+"&"+productKey
func (p *mqttSession) parseRequestOptions(clientId, userName, password string) (*auth.Options, error) {
	options := &auth.Options{}

	names := strings.Split(userName, "&")
	if len(names) != 2 {
		return nil, fmt.Errorf("Invalid authentication user name options:'%s'", userName)
	}
	options.DeviceName = names[0]
	options.ProductKey = names[1]

	names = strings.Split(clientId, "|")
	if len(names) != 2 {
		return nil, fmt.Errorf("Invalid authentication clientId options:'%s'", clientId)
	}
	options.ClientId = names[0]
	names = strings.Split(names[1], ",")
	for _, pair := range names {
		values := strings.Split(pair, "=")
		if len(values) != 2 {
			return nil, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
		}
		switch values[0] {
		case auth.SecurityMode:
			val, err := strconv.Atoi(values[1])
			if err != nil {
				return nil, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
			}
			options.SecurityMode = val
		case auth.SignMethod:
			options.SignMethod = values[1]
		case auth.Timestamp:
			if _, err := strconv.ParseUint(values[1], 10, 64); err != nil {
				return nil, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
			}
			options.Timestamp = values[1]
		default:
			return nil, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
		}
	}
	return options, nil
}

func (p *mqttSession) getMountPoint() string {
	return p.mountpoint // TODO
}
