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
	"github.com/cloustone/sentel/broker/quto"
	"github.com/cloustone/sentel/core"

	"github.com/golang/glog"
)

// MQTT service declaration
type mqttService struct {
	base.ServiceBase
}

// New create mqtt service factory
func New(c core.Config, quit chan os.Signal) (base.Service, error) {
	t := &mqttService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
	}
	return t, nil
}

// MQTT Service

// Name
func (p *mqttService) Name() string {
	return ServiceName
}

func (p *mqttService) Initialize() error { return nil }

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
	go func(p *mqttService) {
		for {
			select {
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// startProtocolService start mqtt protocol on different port
func (p *mqttService) startProtocolService(protocol string, host string) error {
	listen, err := listen(protocol, host, p.Config)
	if err != nil {
		glog.Errorf("Mqtt listen '%s', '%s' failed:%s", protocol, host, err)
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
		// Check wether connection over quto
		if quto.CheckQuto(quto.MaxConnections, 1) != true {
			glog.Error("broker: over quto, closing connection...")
			conn.Close()
			continue
		}
		session, err := newMqttSession(p, conn)
		if err != nil {
			glog.Errorf("Mqtt create session failed:%s", err)
			return err
		}
		go session.Handle()
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
