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

package coap

import (
	"errors"
	"net"
	"sync"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/core"

	uuid "src/github.com/satori/go.uuid"

	"github.com/golang/glog"
)

const (
	protocolName = "coap"
)

type coapService struct {
	core.ServiceBase
	index      int64
	sessions   map[string]base.Session
	mutex      sync.Mutex // Maybe not so good
	protocol   uint8
	localAddrs []string
}

// CoapFactory
type CoapFactory struct{}

// New create coap service factory
func (m *CoapFactory) New(protocol string, c core.Config, ch chan base.ServiceCommand) (base.Service, error) {
	var localAddrs []string = []string{}
	// Get all local ip address
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		glog.Errorf("Failed to get local address:%s", err)
		return nil, err
	}
	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && ipnet.IP.To4() != nil {
			localAddrs = append(localAddrs, ipnet.IP.String())
		}
	}
	if len(localAddrs) == 0 {
		return nil, errors.New("Failed to get local address")
	}
	t := &coapService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			Quit:      ch,
			WaitGroup: sync.WaitGroup{},
		},
		index:      -1,
		sessions:   make(map[string]base.Session),
		protocol:   2,
		localAddrs: localAddrs,
	}
	return t, nil
}

// MQTT Service

func (p *coapService) NewSession(conn net.Conn) (base.Session, error) {
	id := p.CreateSessionId()
	s, err := newCoapSession(m, conn, id)
	return s, err
}

// CreateSessionId create id for new session
func (p *coapService) CreateSessionId() string {
	return uuid.NewV4().String()
}

// GetSessionTotalCount get total session count
func (p *coapService) GetSessionTotalCount() int64 {
	return int64(len(p.sessions))
}

func (p *coapService) RemoveSession(s base.Session) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.sessions[s.Identifier()] = nil
}
func (p *coapService) RegisterSession(s base.Session) {
	p.mutex.Lock()
	p.sessions[s.Identifier()] = s
	p.mutex.Unlock()
}

// Start
func (p *coapService) Start() error {
	host, _ := p.config.String("coap", "host")

	listen, err := net.Listen("tcp", host)
	if err != nil {
		glog.Errorf("Coap listen failed:%s", err)
		return err
	}
	glog.Infof("Coap server is listening on '%s'...", host)
	for {
		conn, err := listen.Accept()
		if err != nil {
			continue
		}
		session, err := p.NewSession(conn)
		if err != nil {
			glog.Error("Mqtt create session failed")
			return err
		}
		glog.Infof("Mqtt new connection:%s", session.Identifier())
		p.RegisterSession(session)
		go session.Handle()
	}
	// notify main
	// p.chn <- 1
	return nil
}

func (p *coapService) Stop() {}

// Name
func (p *coapService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "coap",
	}
}

func (p *coapService) GetMetrics() *base.Metrics            { return nil }
func (p *coapService) GetStats() *base.Stats                { return nil }
func (p *coapService) GetClients() []*base.ClientInfo       { return nil }
func (p *coapService) GetClient(id string) *base.ClientInfo { return nil }
func (p *coapService) KickoffClient(id string) error        { return nil }

func (p *coapService) GetSessions(conditions map[string]bool) []*base.SessionInfo { return nil }
func (p *coapService) GetSession(id string) *base.SessionInfo                     { return nil }

func (p *coapService) GetRoutes() []*base.RouteInfo { return nil }
func (p *coapService) GetRoute() *base.RouteInfo    { return nil }

// Topic info
func (p *coapService) GetTopics() []*base.TopicInfo       { return nil }
func (p *coapService) GetTopic(id string) *base.TopicInfo { return nil }

// SubscriptionInfo
func (p *coapService) GetSubscriptions() []*base.SubscriptionInfo       { return nil }
func (p *coapService) GetSubscription(id string) *base.SubscriptionInfo { return nil }

// Service Info
func (p *coapService) GetServiceInfo() *base.ServiceInfo { return nil }
