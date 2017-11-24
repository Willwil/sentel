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
	"errors"
	"fmt"
	"net"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

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

// MQTT service declaration
type mqttService struct {
	core.ServiceBase
	sessions   map[string]base.Session // All mqtt sessions
	mutex      sync.Mutex              // Mutex to protect sessions
	localAddrs []string                // Local address used to identifier wether notification come from local
	protocol   string                  // Supported protocol, such as tcp,websocket, tls
	consumer   sarama.Consumer         // Kafka client connection handle
}

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
		consumer:   consumer,
	}
	return t, nil
}

// MQTT Service

// Name
func (p *mqttService) Name() string {
	return ServiceName
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
		ServiceName: ServiceName,
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
	p.WaitGroup.Add(1)
	defer p.WaitGroup.Done()
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
		go func(s *mqttSession) {
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
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
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
