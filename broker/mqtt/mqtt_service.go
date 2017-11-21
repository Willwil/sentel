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
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"
	uuid "github.com/satori/go.uuid"

	"github.com/Shopify/sarama"
	"github.com/golang/glog"
)

const (
	maxMqttConnections = 1000000
	protocolName       = "mqtt3"
)

const (
	mqttNetworkTcp       = "tcp"
	mqttNetworkTls       = "tls"
	mqttNetworkWebsocket = "ws"
	mqttNetworkHttps     = "https"
)

// MQTT service declaration
type mqttService struct {
	core.ServiceBase
	sessions   map[string]base.Session // All mqtt sessions
	mutex      sync.Mutex              // Mutex to protect sessions
	localAddrs []string                // Local address used to identifier wether notification come from local
	storage    Storage                 // Storage for sessions and metadata
	protocol   string                  // Supported protocol, such as tcp,websocket, tls
	consumer   sarama.Consumer         // Kafka client connection handle
}

const (
	MqttProtocolTcp = "tcp"
	MqttProtocolWs  = "ws"
	MqttProtocolTls = "tls"
)

// MqttFactory
type MqttFactory struct {
	Protocol string // Indicate which protocol the factory to support
}

// New create mqtt service factory
func (p *MqttFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// Get all local ip address
	localAddrs := []string{}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		glog.Errorf("Failed to get local interface:%s", err)
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			return nil, errors.New("Failed to get local address")
		}
	}

	// Create storage
	s, err := NewStorage(c)
	if err != nil {
		return nil, errors.New("Failed to create storage in mqtt")
	}
	// kafka
	khosts, _ := core.GetServiceEndpoint(c, "broker", "kafka")
	consumer, err := sarama.NewConsumer(strings.Split(khosts, ","), nil)
	if err != nil {
		return nil, fmt.Errorf("Connecting with kafka:%s failed", khosts)
	}

	t := &mqttService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
		sessions:   make(map[string]base.Session),
		protocol:   p.Protocol,
		localAddrs: localAddrs,
		storage:    s,
		consumer:   consumer,
	}
	return t, nil
}

// MQTT Service

// Name
func (p *mqttService) Name() string {
	return "mqtt:" + p.protocol
}

// removeSession remove specified session from mqtt service
func (p *mqttService) removeSession(s base.Session) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	delete(p.sessions, s.Identifier())
}

// addSession add newly created session into mqtt service
func (p *mqttService) addSession(s base.Session) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	id := s.Identifier()
	if _, ok := p.sessions[id]; ok {
		return fmt.Errorf("Mqtt session '%s' is already regisitered", id)
	}
	p.sessions[id] = s
	return nil
}

// Info
func (p *mqttService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "mqtt",
	}
}

// Client
func (p *mqttService) GetClients() []*base.ClientInfo       { return nil }
func (p *mqttService) GetClient(id string) *base.ClientInfo { return nil }
func (p *mqttService) KickoffClient(id string) error        { return nil }

// Session Info
func (p *mqttService) GetSessions(conditions map[string]bool) []*base.SessionInfo { return nil }
func (p *mqttService) GetSession(id string) *base.SessionInfo                     { return nil }

// Route Info
func (p *mqttService) GetRoutes() []*base.RouteInfo { return nil }
func (p *mqttService) GetRoute() *base.RouteInfo    { return nil }

// Topic info
func (p *mqttService) GetTopics() []*base.TopicInfo       { return nil }
func (p *mqttService) GetTopic(id string) *base.TopicInfo { return nil }

// SubscriptionInfo
func (p *mqttService) GetSubscriptions() []*base.SubscriptionInfo       { return nil }
func (p *mqttService) GetSubscription(id string) *base.SubscriptionInfo { return nil }

// Service Info
func (p *mqttService) GetServiceInfo() *base.ServiceInfo { return nil }

// Start is mainloop for mqtt service
func (p *mqttService) Start() error {
	if err := p.subscribeTopic("session"); err != nil {
		return err
	}
	// Read protocol configuration for supported protocol
	protocolConfig, err := p.Config.String("mqtt", "protocols")
	if err != nil || protocolConfig == "" {
		return errors.New("Invalid mqtt protocol configuration")
	}
	protocols := strings.Split(protocolConfig, ",")
	if len(protocols) == 0 {
		return errors.New("No protocol service for mqtt broker")
	}
	for _, protocol := range protocols {
		host, err := p.Config.String("mqtt", protocol)
		if err != nil {
			return err
		}
		go p.startProtocolService(protocol, host)
	}
	return nil
}

// startProtocolService start mqtt protocol on different port
func (p *mqttService) startProtocolService(protocol string, host string) error {
	listen, err := listen(protocol, host, p.Config)
	if err != nil {
		glog.Errorf("Mqtt listen failed:%s", err)
		return err
	}
	glog.Infof("Mqtt service '%s' is listening on '%s'...", protocol, host)
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}

		id := uuid.NewV4().String()
		session, err := newMqttSession(p, conn, id)
		if err != nil {
			glog.Errorf("Mqtt create session failed:%s", err)
			return err
		}
		p.addSession(session)
		go func(s base.Session) {
			err := s.Handle()
			if err != nil {
				conn.Close()
				glog.Error(err)
			}
		}(session)
	}
}

// Stop
func (p *mqttService) Stop() {
}

// subscribeTopc subscribe topics from apiserver
func (p *mqttService) subscribeTopic(topic string) error {
	partitionList, err := p.consumer.Partitions(topic)
	if err != nil {
		return fmt.Errorf("Failed to get list of partions:%v", err)
		return err
	}

	for partition := range partitionList {
		pc, err := p.consumer.ConsumePartition(topic, int32(partition), sarama.OffsetNewest)
		if err != nil {
			glog.Errorf("Failed  to start consumer for partion %d:%s", partition, err)
			continue
		}
		defer pc.AsyncClose()
		p.WaitGroup.Add(1)

		go func(sarama.PartitionConsumer) {
			defer p.WaitGroup.Done()
			for msg := range pc.Messages() {
				p.handleNotifications(string(msg.Topic), msg.Value)
			}
		}(pc)
	}
	return nil
}

// handleNotifications handle notification from kafka
func (p *mqttService) handleNotifications(topic string, value []byte) error {
	switch topic {
	case TopicNameSession:
		return p.handleSessionNotifications(value)
	}
	return nil
}

// handleSessionNotifications handle session notification  from kafka
func (p *mqttService) handleSessionNotifications(value []byte) error {
	// Decode value received form other mqtt node
	var topics []SessionTopic
	if err := json.Unmarshal(value, &topics); err != nil {
		glog.Errorf("Mqtt session notifications failure:%s", err)
		return err
	}
	// Get local ip address
	for _, topic := range topics {
		switch topic.Action {
		case core.TopicActionUpdate:
			// Only deal with notification that is not  launched by myself
			for _, addr := range p.localAddrs {
				if addr != topic.Launcher {
					s, err := p.storage.FindSession(topic.SessionId)
					if err != nil {
						s.state = topic.State
					}
					//p.storage.UpdateSession(&StorageSession{Id: topic.SessionId, State: topic.State})
				}
			}
		case core.TopicActionDelete:

		case core.TopicActionRegister:
		default:
		}
	}
	return nil
}

func listen(network, laddr string, c core.Config) (net.Listener, error) {
	switch network {
	case mqttNetworkTcp:
		return net.Listen("tcp", laddr)
	case mqttNetworkTls:
		if _, err := c.String("security", "crt_file"); err != nil {
			return nil, err
		}
		if _, err := c.String("security", "key_file"); err != nil {
			return nil, err
		}
		crt := c.MustString("security", "crt_file")
		key := c.MustString("security", "key_file")
		cer, err := tls.LoadX509KeyPair(crt, key)
		if err != nil {
			return nil, err
		}
		config := tls.Config{Certificates: []tls.Certificate{cer}}
		return tls.Listen("tcp", laddr, &config)

	case mqttNetworkWebsocket:
	case mqttNetworkHttps:

	}
	return nil, fmt.Errorf("Unsupported network protocol '%s'", network)
}
