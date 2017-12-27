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

package http

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/pkg/config"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type httpService struct {
	base.ServiceBase
}

const (
	ServiceName = "http"
)

type ServiceFactory struct{}

// New create metadata service factory
func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (base.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return &httpService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
	}, nil

}

// Name
func (p *httpService) Name() string {
	return ServiceName
}

func (p *httpService) Initialize() error { return nil }

// Start
func (p *httpService) Start() error {
	return nil
}

// Stop
func (p *httpService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

// handleNotifications handle notification from kafka
func (p *httpService) handleNotifications(topic string, value []byte) error {
	return nil
}
