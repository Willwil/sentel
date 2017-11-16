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

func (this *coapService) NewSession(conn net.Conn) (base.Session, error) {
	id := this.CreateSessionId()
	s, err := newCoapSession(m, conn, id)
	return s, err
}

// CreateSessionId create id for new session
func (this *coapService) CreateSessionId() string {
	return uuid.NewV4().String()
}

// GetSessionTotalCount get total session count
func (this *coapService) GetSessionTotalCount() int64 {
	return int64(len(this.sessions))
}

func (this *coapService) RemoveSession(s base.Session) {
	this.mutex.Lock()
	defer this.mutex.Unlock()
	this.sessions[s.Identifier()] = nil
}
func (this *coapService) RegisterSession(s base.Session) {
	this.mutex.Lock()
	this.sessions[s.Identifier()] = s
	this.mutex.Unlock()
}

// Start
func (this *coapService) Start() error {
	host, _ := this.config.String("coap", "host")

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
		session, err := this.NewSession(conn)
		if err != nil {
			glog.Error("Mqtt create session failed")
			return err
		}
		glog.Infof("Mqtt new connection:%s", session.Identifier())
		this.RegisterSession(session)
		go session.Handle()
	}
	// notify main
	// this.chn <- 1
	return nil
}

func (this *coapService) Stop() {}

// Name
func (this *coapService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "coap",
	}
}

func (this *coapService) GetMetrics() *base.Metrics            { return nil }
func (this *coapService) GetStats() *base.Stats                { return nil }
func (this *coapService) GetClients() []*base.ClientInfo       { return nil }
func (this *coapService) GetClient(id string) *base.ClientInfo { return nil }
func (this *coapService) KickoffClient(id string) error        { return nil }

func (this *coapService) GetSessions(conditions map[string]bool) []*base.SessionInfo { return nil }
func (this *coapService) GetSession(id string) *base.SessionInfo                     { return nil }

func (this *coapService) GetRoutes() []*base.RouteInfo { return nil }
func (this *coapService) GetRoute() *base.RouteInfo    { return nil }

// Topic info
func (this *coapService) GetTopics() []*base.TopicInfo       { return nil }
func (this *coapService) GetTopic(id string) *base.TopicInfo { return nil }

// SubscriptionInfo
func (this *coapService) GetSubscriptions() []*base.SubscriptionInfo       { return nil }
func (this *coapService) GetSubscription(id string) *base.SubscriptionInfo { return nil }

// Service Info
func (this *coapService) GetServiceInfo() *base.ServiceInfo { return nil }
