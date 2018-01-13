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
	"math"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	"github.com/cloustone/sentel/broker/queue"
	sm "github.com/cloustone/sentel/broker/sessionmgr"
	"github.com/cloustone/sentel/pkg/config"

	auth "github.com/cloustone/sentel/broker/auth"

	"github.com/golang/glog"
)

// mqttSession handle mqtt protocol session
type mqttSession struct {
	config         config.Config    // Session Configuration
	conn           net.Conn         // Underlay network connection handle
	clientId       string           // Session's client identifier
	keepalive      uint16           // MQTT protocol kepalive time
	protocol       uint8            // MQTT protocol version
	cleanSession   uint8            // Wether the session is clean session
	willMsg        *base.Message    // MQTT protocol will message
	mountpoint     string           // Topic's mount point
	msgState       uint8            // Mqtt message state
	sessionState   uint8            // Mqtt session state
	authctx        auth.Context     // authentication options
	aliveTimer     *time.Timer      // Keepalive timer
	pingTime       *time.Time       // Ping timer
	availableChan  chan int         // Event chanel for data availabel notification
	packetChan     chan *mqttPacket // Mqtt packet channel
	errorChan      chan error       // Error channel
	queue          queue.Queue      // Session's queue
	nextPacketId   uint16           // Next Packet id index
	lastMessageIn  time.Time        // Last message's in time
	lastMessageOut time.Time        // last message out time
	bytesReceived  int64            // Byte recivied
	authNeed       bool             // Indicate wether authentication is needed
	metrics        metric.Metric    // Metrics
}

// newMqttSession create new session  for each client connection
func newMqttSession(mqtt *mqttService, conn net.Conn) (*mqttSession, error) {
	// Retrieve authentication option
	authNeed := true
	if n, err := mqtt.config.Bool("broker", "auth"); err == nil {
		authNeed = n
	}

	// Create session instance
	return &mqttSession{
		config:        mqtt.config,
		conn:          conn,
		mountpoint:    "",
		bytesReceived: 0,
		sessionState:  mqttStateNew,
		protocol:      mqttProtocolInvalid,
		msgState:      mqttMsgStateQueued,
		aliveTimer:    nil,
		availableChan: make(chan int),
		packetChan:    make(chan *mqttPacket),
		errorChan:     make(chan error),
		queue:         nil,
		pingTime:      nil,
		nextPacketId:  0,
		clientId:      "",
		willMsg:       nil,
		authNeed:      authNeed,
		metrics:       metric.NewMetric(ServiceName),
	}, nil
}

// checkTopiValidity will check topic's validity
func checkTopicValidity(topic string) error {
	return nil
}

// setSessionState set session's state
func (p *mqttSession) setSessionState(state uint8) {
	glog.Infof("mqtt: client '%s' session state changed (%s->%s)", p.clientId,
		nameOfSessionState(p.sessionState),
		nameOfSessionState(state))
	p.sessionState = state
}

// checkSessionState check wether the session state is same with expected state
func (p *mqttSession) checkSessionState(state uint8) error {
	if p.sessionState != state {
		return fmt.Errorf("mqtt: state miss occurred, expected:%s, now:%s",
			nameOfSessionState(state), nameOfSessionState(p.sessionState))
	}
	return nil
}

// setMessageState set session's state
func (p *mqttSession) setMessageState(state uint8) {
	glog.Infof("mqtt: client '%s' message state changed (%s->%s)", p.clientId,
		nameOfMessageState(p.msgState),
		nameOfMessageState(state))
	p.msgState = state
}

// Handle is mainprocessor for iot device client
func (p *mqttSession) Handle() error {
	// Start goroutine to read packet from client
	go func(p *mqttSession) {
		for {
			packet := newMqttPacket()
			if err := packet.DecodeFromReader(p.conn, base.NilDecodeFeedback{}); err != nil {
				p.errorChan <- fmt.Errorf("mqtt packet decoder failed:%s", err.Error())
				return
			}
			glog.Infof("Received mqtt packet '%s'", nameOfPacket(packet))
			if p.aliveTimer != nil {
				p.aliveTimer.Reset(time.Duration(int(float32(p.keepalive)*1.5)) * time.Second)
			}
			p.packetChan <- packet
		}
	}(p)

	for {
		var err error
		select {
		case e := <-p.errorChan:
			err = e
		case <-p.availableChan:
			err = p.handleDataAvailableNotification()
		case packet := <-p.packetChan:
			err = p.handleMqttPacket(packet)
		}
		// Disconnect the connection when any errors occured
		if err != nil {
			glog.Error(err)
			p.disconnect(err)
			return err
		}
		// Return from session main routine when session is disconnected
		if err = p.checkSessionState(mqttStateDisconnected); err == nil {
			break
		}
	}
	return nil
}

// handleMqttInPacket handle all mqtt in packet
func (p *mqttSession) handleMqttPacket(packet *mqttPacket) error {
	glog.Infof("mqtt packet type:%s, clientId:%s, session state:%s, message state:%s",
		nameOfPacket(packet),
		p.clientId,
		nameOfSessionState(p.sessionState),
		nameOfMessageState(p.msgState))
	p.updatePacketMetrics(packet)
	switch packet.command & 0xF0 {
	case PINGREQ:
		return p.handlePingReq(packet)
	case CONNECT:
		return p.handleConnect(packet)
	case DISCONNECT:
		return p.handleDisconnect(packet)
	case PUBLISH:
		return p.handlePublish(packet)
	case PUBREL:
		return p.handlePubRel(packet)
	case SUBSCRIBE:
		return p.handleSubscribe(packet)
	case UNSUBSCRIBE:
		return p.handleUnsubscribe(packet)
	case PUBACK:
		return p.handlePubAck(packet)
	default:
		return fmt.Errorf("Unrecognized protocol command:%d", int(packet.command&0xF0))
	}
}

// handleDataAvailableNotification read message from queue and send to client
func (p *mqttSession) handleDataAvailableNotification() error {
	if p.msgState != mqttMsgStateQueued {
		return fmt.Errorf("mqtt invalid message state:%s", nameOfMessageState(p.msgState))
	}
	for p.queue.Length() > 0 {
		if msg := p.queue.Front(); msg != nil {
			if err := p.sendPublish(msg); err == nil {
				p.metrics.Add(metric.MessageSent, 1)
				switch msg.Qos {
				case 0:
					// No client response for qos0 pub packet
					p.queue.Pop()
					p.metrics.Add(metric.PacketPublishSent, 1)
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

// updateMtrics update session metrics
func (p *mqttSession) updatePacketMetrics(packet *mqttPacket) {
	switch packet.command & 0xF0 {
	case PINGREQ:
		p.metrics.Add(metric.PacketPingreq, 1)
	case CONNECT:
		p.metrics.Add(metric.PacketConnect, 1)
	case DISCONNECT:
		p.metrics.Add(metric.PacketDisconnect, 1)
	case PUBLISH:
		p.metrics.Add(metric.PacketPublishReceived, 1)
	case PUBREL:
		p.metrics.Add(metric.PacketPubrelReceived, 1)
	case SUBSCRIBE:
		p.metrics.Add(metric.PacketSubscribe, 1)
	case UNSUBSCRIBE:
		p.metrics.Add(metric.PacketUnsubscribe, 1)
	case PUBACK:
		p.metrics.Add(metric.PacketPubackReceived, 1)
	}
	p.metrics.Add(metric.PacketReceived, 1)
}

// handleAliveTimeout reply to client with ping rsp after timeout
func (p *mqttSession) handleAliveTimeout() error {
	return errors.New("mqtt: time out occured")
}

// DataAvailable called when queue have data to deal with
func (p *mqttSession) DataAvailable(q queue.Queue, msg *base.Message) {
	p.availableChan <- 1
}

func (p *mqttSession) Id() string            { return p.clientId }
func (p *mqttSession) BrokerId() string      { return base.GetBrokerId() }
func (p *mqttSession) Info() *sm.SessionInfo { return nil }
func (p *mqttSession) IsValid() bool         { return true }
func (p *mqttSession) IsPersistent() bool    { return (p.cleanSession == 0) }

// handleConnect handle connect packet
func (p *mqttSession) handleConnect(packet *mqttPacket) error {
	if err := p.checkSessionState(mqttStateNew); err != nil {
		return err
	}
	// Check protocol name and version
	protocolName, err := packet.readString()
	if err != nil {
		return err
	}
	protocolVersion, err := packet.readByte()
	if err != nil {
		return err
	}
	// Check connect flags
	cflags, err := packet.readByte()
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
		if packet.command&0x0F != 0x00 {
			return fmt.Errorf("Invalid protocol version '%d' in CONNECT packet", protocolVersion)
		}
		p.protocol = mqttProtocol311
		if cflags&0x01 != 0x00 {
			return errors.New("Invalid protocol version in connect flags")
		}
	default:
		return fmt.Errorf("Invalid protocol name '%s' in CONNECT packet", protocolName)
	}

	// Analyze connecton flag
	cleanSession := (cflags & 0x02) >> 1
	will := cflags & 0x04
	willQos := (cflags & 0x18) >> 3
	if willQos >= 3 {
		return fmt.Errorf("Invalid Will Qos in CONNECT from %s", p.clientId)
	}
	willRetain := (cflags & 0x20) == 0x20

	// Keepalive and client identifier
	keepalive, err := packet.readUint16()
	if err != nil {
		return err
	}
	p.keepalive = keepalive

	// Deal with client identifier
	clientId, err := packet.readString()
	if err != nil || clientId == "" {
		p.sendConnAck(0, CONNACK_REFUSED_IDENTIFIER_REJECTED)
		return errors.New("Invalid mqtt packet with client id")
	}

	// Deal with will message
	var willMsg *base.Message = nil
	if will > 0 {
		topic, err := packet.readString()
		if err != nil || topic == "" {
			return mqttErrorInvalidProtocol
		}
		willTopic := p.mountpoint + topic
		if err := checkTopicValidity(willTopic); err != nil {
			return err
		}
		// Get willtopic's payload
		willPayloadLength, err := packet.readUint16()
		if err != nil {
			return mqttErrorInvalidProtocol
		}
		var payload []uint8
		if willPayloadLength > 0 {
			payload, err = packet.readBytes(int(willPayloadLength))
			if err != nil {
				return mqttErrorInvalidProtocol
			}
		}
		willMsg = &base.Message{
			Topic:   willTopic,
			Qos:     willQos,
			Retain:  willRetain,
			Payload: payload,
		}
	} else {
		if p.protocol == mqttProtocol311 {
			if willQos != 0 || willRetain {
				return mqttErrorInvalidProtocol
			}
		}
	}

	// Get user name and password and parse mqtt client request options to authentication
	if cflags&0x40 == 0 || cflags&0x80 == 0 {
		return mqttErrorInvalidProtocol
	}
	username, usrErr := packet.readString()
	password, pwdErr := packet.readString()
	if usrErr != nil || pwdErr != nil {
		return mqttErrorInvalidProtocol
	}
	p.authctx, err = p.parseRequestOptions(clientId, username, password)
	if err != nil {
		return err
	}
	if p.authNeed {
		err = auth.Authenticate(p.authctx)
		switch err {
		case nil:
			// Successfuly authenticated
		case auth.ErrUnauthorized:
			p.sendConnAck(0, CONNACK_REFUSED_NOT_AUTHORIZED)
			p.disconnect(err)
			return err
		default:
			p.disconnect(err)
			return err
		}
	}
	p.clientId = p.authctx.ClientId

	// Find if the client already has an entry, p must be done after any security check
	conack := 0
	if cleanSession == 0 {
		if found, _ := sm.FindSession(clientId); found != nil {
			// Found old session
			if !found.IsValid() {
				glog.Errorf("Invalid session(%s) in store", found.Id())
			}
			info := found.Info()
			if p.protocol == mqttProtocol311 {
				if cleanSession == 0 {
					conack |= 0x01
				}
			}
			if p.cleanSession == 0 && info.CleanSession == 0 {
				// Resume last session and notify other mqtt node to release resource
				event.Notify(event.SessionResume, clientId, nil)
			}
		}
	}
	// Remove any queued messages that are no longer allowd through ACL
	// Assuming a possible change of username
	if willMsg != nil {
		metadata.DeleteMessageWithValidator(
			clientId,
			func(msg *base.Message) bool {
				err := auth.Authorize(p.authctx, clientId, msg.Topic, auth.AclRead)
				return err != nil
			})
		metadata.AddMessage(clientId, willMsg)
	}
	p.willMsg = willMsg
	p.cleanSession = cleanSession

	// Create queue for this sesion and otify event service that new session created
	if q, err := queue.NewQueue(p.clientId, (cleanSession == 0), p); err != nil {
		glog.Error(err)
		return err
	} else {
		p.queue = q
	}
	// Change session state and reply client
	p.setSessionState(mqttStateConnected)
	err = p.sendConnAck(uint8(conack), CONNACK_ACCEPTED)
	sm.RegisterSession(p)
	event.Notify(event.SessionCreate, p.clientId,
		&event.SessionCreateDetail{
			Persistent: (cleanSession == 0),
		})

	// start keep alive timer
	if p.keepalive > 0 {
		p.aliveTimer = time.NewTimer(time.Duration(int(float32(p.keepalive)*1.5)) * time.Second)
		go func(p *mqttSession) {
			for {
				select {
				case <-p.aliveTimer.C:
					p.errorChan <- fmt.Errorf("mqtt time out for client '%s'", p.clientId)
					return
				}
			}
		}(p)
	}

	return nil
}

// handleDisconnect handle disconnect packet
func (p *mqttSession) handleDisconnect(packet *mqttPacket) error {
	if packet.remainingLength != 0 {
		return mqttErrorInvalidProtocol
	}
	if p.protocol == mqttProtocol311 && (packet.command&0x0F) != 0x00 {
		return mqttErrorInvalidProtocol
	}
	// WillMessage must be deleted when received disconnect packet
	if p.willMsg != nil {
		metadata.DeleteMessage(p.clientId, p.willMsg.PacketId, 0)
		p.willMsg = nil
	}
	p.setSessionState(mqttStateDisconnecting)
	p.disconnect(nil)
	return nil
}

// disconnect will disconnect current connection because of protocol error
func (p *mqttSession) disconnect(reason error) {
	if err := p.checkSessionState(mqttStateDisconnected); err != nil {
		glog.Infof("mqtt session '%s' is disconnecting...", p.clientId)
		// Publish will message if session is not normoally disconnected
		if reason != nil && p.willMsg != nil {
			event.Notify(event.TopicPublish, p.clientId,
				&event.TopicPublishDetail{
					Topic:   p.willMsg.Topic,
					Payload: p.willMsg.Payload,
					Qos:     p.willMsg.Qos,
					Retain:  p.willMsg.Retain,
				})
		}
		event.Notify(event.SessionDestroy, p.clientId, nil)
		p.setSessionState(mqttStateDisconnected)
		if p.aliveTimer != nil {
			p.aliveTimer.Stop()
		}
		p.conn.Close()
		p.conn = nil
		// Close channel
		close(p.availableChan)
		close(p.packetChan)
		// Releae metrics
		metric.FreeMetric(ServiceName, p.metrics)
	}
}

// handleSubscribe handle subscribe packet
func (p *mqttSession) handleSubscribe(packet *mqttPacket) error {
	payload := make([]uint8, 0)
	if p.protocol == mqttProtocol311 && packet.command&0x0F != 0x02 {
		return mqttErrorInvalidProtocol
	}
	// Get packet identifier
	pid, err := packet.readUint16()
	if err != nil {
		return err
	}
	// Deal each subscription
	for packet.pos < packet.remainingLength {
		topic := ""
		qos := uint8(0)
		if topic, err = packet.readString(); err != nil {
			return err
		}
		if checkTopicValidity(topic) != nil {
			glog.Errorf("Invalid subscription topic %s from %s, disconnecting", topic, p.clientId)
			return mqttErrorInvalidProtocol
		}
		if qos, err = packet.readByte(); err != nil {
			return err
		}

		if qos > 2 {
			glog.Errorf("Invalid Qos in subscription %s from %s", topic, p.clientId)
			return mqttErrorInvalidProtocol
		}

		topic = p.mountpoint + topic
		if qos != 0x80 {
			event.Notify(event.TopicSubscribe, p.clientId,
				&event.TopicSubscribeDetail{Topic: topic, Qos: qos, Retain: true})
		}
		payload = append(payload, qos)
	}

	if p.protocol == mqttProtocol311 && len(payload) == 0 {
		return mqttErrorInvalidProtocol
	}
	return p.sendSubAck(pid, payload)
}

// handleUnsubscribe handle unsubscribe packet
func (p *mqttSession) handleUnsubscribe(packet *mqttPacket) error {
	if p.protocol == mqttProtocol311 && (packet.command&0x0f) != 0x02 {
		return mqttErrorInvalidProtocol
	}
	pid, err := packet.readUint16()
	if err != nil {
		return err
	}
	// Iterate all subscription
	for packet.pos < packet.remainingLength {
		topic, err := packet.readString()
		if err != nil {
			return mqttErrorInvalidProtocol
		}
		if err := checkTopicValidity(topic); err != nil {
			return fmt.Errorf("Invalid unsubscription string from %s, disconnecting", p.clientId)
		}
		event.Notify(event.TopicUnsubscribe, p.clientId, &event.TopicUnsubscribeDetail{Topic: topic})
	}

	return p.sendCommandWithPacketId(UNSUBACK, pid, false)
}

// handlePublish handle publish packet
func (p *mqttSession) handlePublish(packet *mqttPacket) error {
	var topic string
	var pid uint16
	var err error
	var payload []uint8

	dup := (packet.command & 0x08) >> 3
	qos := (packet.command & 0x06) >> 1
	if qos > 2 {
		return fmt.Errorf("Invalid Qos in PUBLISH from %s, disconnectiing", p.clientId)
	}
	retain := (packet.command & 0x01)

	// Topic
	if topic, err = packet.readString(); err != nil {
		return fmt.Errorf("Invalid topic in PUBLISH from %s", p.clientId)
	}
	if checkTopicValidity(topic) != nil {
		return fmt.Errorf("Invalid topic in PUBLISH(%s) from %s", topic, p.clientId)
	}
	topic = p.mountpoint + topic
	if qos > 0 {
		pid, err = packet.readUint16()
		if err != nil {
			return err
		}
	}

	// Payload
	payloadlen := packet.remainingLength - packet.pos
	if payloadlen > 0 {
		limitSize, _ := p.config.Int("mqtt", "message_size_limit")
		if payloadlen > limitSize {
			return mqttErrorInvalidProtocol
		}
		payload, err = packet.readBytes(payloadlen)
		if err != nil {
			return err
		}
	}
	// Check for topic access
	if err := auth.Authorize(p.authctx, p.clientId, topic, auth.AclWrite); err != nil {
		return mqttErrorConnectRefused
	}
	glog.Infof("Received PUBLISH from %s(d:%d, q:%d r:%d, m:%d, '%s',..(%d)bytes",
		p.clientId, dup, qos, retain, pid, topic, payloadlen)

	// Check wether the message has been stored :TODO
	dup = 0
	if qos > 0 && sm.FindMessage(p.clientId, pid, 1) != nil {
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
		event.Notify(event.TopicPublish, p.clientId, &detail)
	case 1:
		event.Notify(event.TopicPublish, p.clientId, &detail)
		err = p.sendPubAck(pid)
	case 2:
		err = errors.New("MQTT qos 2 is not supported now")
	default:
		err = mqttErrorInvalidProtocol
	}

	return err
}

// handlePubAck handle pubrel packet
func (p *mqttSession) handlePubAck(packet *mqttPacket) error {
	if p.msgState == mqttMsgStateWaitPubAck {
		q := queue.GetQueue(p.clientId)
		q.Pop()
		p.setMessageState(mqttMsgStateQueued)
		p.availableChan <- 1
	}
	return nil
}

// handlePubRel handle pubrel packet
func (p *mqttSession) handlePubRel(packet *mqttPacket) error {
	// Check protocol specifal requriement
	if p.protocol == mqttProtocol311 {
		if (packet.command & 0x0F) != 0x02 {
			return mqttErrorInvalidProtocol
		}
	}
	// Get message identifier
	pid, err := packet.readUint16()
	if err != nil {
		return err
	}

	//sm.DeleteMessage(p.clientId, pid, sm.MessageDirectionIn) // TODO:Qos2
	return p.sendPubComp(pid)
}

// sendSimpleCommand send a simple command
func (p *mqttSession) sendSimpleCommand(cmd uint8) error {
	packet := &mqttPacket{
		command:         cmd,
		remainingLength: 1,
	}
	packet.initializePacket()
	return p.writePacket(packet)
}

// handlePingReq handle ping request packet
func (p *mqttSession) handlePingReq(packet *mqttPacket) error {
	return p.sendPingRsp()
}

// sendPingRsp send ping response to client
func (p *mqttSession) sendPingRsp() error {
	return p.sendSimpleCommand(PINGRESP)
}

// sendConnAck send connection response to client
func (p *mqttSession) sendConnAck(ack uint8, result uint8) error {
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
func (p *mqttSession) sendSubAck(pid uint16, payload []uint8) error {
	packet := &mqttPacket{
		command:         SUBACK,
		remainingLength: 2 + int(len(payload)),
	}

	packet.initializePacket()
	packet.writeUint16(pid)
	if len(payload) > 0 {
		packet.writeBytes(payload)
	}
	return p.writePacket(packet)
}

// sendCommandWithPacketId send command with message identifier
func (p *mqttSession) sendCommandWithPacketId(command uint8, pid uint16, dup bool) error {
	packet := &mqttPacket{
		command:         command,
		remainingLength: 2,
	}
	if dup {
		packet.command |= 8
	}
	packet.initializePacket()
	packet.payload[packet.pos+0] = uint8((pid & 0xFF00) >> 8)
	packet.payload[packet.pos+1] = uint8(pid & 0xff)
	return p.writePacket(packet)
}

// sendPubAck
func (p *mqttSession) sendPubAck(pid uint16) error {
	glog.Infof("Sending PUBACK to %s with MID:%d", p.clientId, pid)
	return p.sendCommandWithPacketId(PUBACK, pid, false)
}

// sendPubRec
func (p *mqttSession) sendPubRec(pid uint16) error {
	glog.Infof("Sending PUBRREC to %s with MID:%d", p.clientId, pid)
	return p.sendCommandWithPacketId(PUBREC, pid, false)
}

func (p *mqttSession) sendPubComp(pid uint16) error {
	glog.Infof("Sending PUBCOMP to %s with MID:%d", p.clientId, pid)
	return p.sendCommandWithPacketId(PUBCOMP, pid, false)
}

func (p *mqttSession) makePacketId() uint16 {
	if p.nextPacketId < math.MaxUint16 {
		pid := p.nextPacketId
		p.nextPacketId += 1
		return pid
	}
	p.nextPacketId = 0
	return p.nextPacketId
}

func (p *mqttSession) sendPublish(msg *base.Message) error {
	qos := msg.Qos
	pid := uint16(0)

	// TODO:upgrade qos
	if qos > 0 {
		pid = p.makePacketId()
	}

	packet := newMqttPacket()
	packet.command = PUBLISH
	packet.remainingLength = 2 + len(msg.Topic) + len(msg.Payload)
	if qos > 0 {
		packet.remainingLength += 2
	}
	packet.initializePacket()
	packet.writeString(msg.Topic)
	if qos > 0 {
		packet.writeUint16(pid)
	}
	packet.writeBytes(msg.Payload)
	return p.writePacket(packet)
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
func (p *mqttSession) parseRequestOptions(clientId, userName, password string) (auth.Context, error) {
	ctx := auth.Context{}

	names := strings.Split(userName, "&")
	if len(names) != 2 {
		return ctx, fmt.Errorf("Invalid authentication user name options:'%s'", userName)
	}
	ctx.DeviceName = names[0]
	ctx.ProductId = names[1]

	names = strings.Split(clientId, "|")
	if len(names) != 2 {
		return ctx, fmt.Errorf("Invalid authentication clientId options:'%s'", clientId)
	}
	ctx.ClientId = names[0]
	names = strings.Split(names[1], ",")
	for _, pair := range names {
		values := strings.Split(pair, "=")
		if len(values) != 2 {
			return ctx, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
		}
		switch values[0] {
		case auth.SecurityMode:
			val, err := strconv.Atoi(values[1])
			if err != nil {
				return ctx, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
			}
			ctx.SecurityMode = val
		case auth.SignMethod:
			ctx.SignMethod = values[1]
		case auth.Timestamp:
			if _, err := strconv.ParseUint(values[1], 10, 64); err != nil {
				return ctx, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
			}
			ctx.Timestamp = values[1]
		default:
			return ctx, fmt.Errorf("Invalid authentication clientId options:'%s'", pair)
		}
	}
	return ctx, nil
}
