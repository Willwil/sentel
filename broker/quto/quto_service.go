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

package quto

import (
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/core"

	"gopkg.in/mgo.v2"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type QutoService struct {
	core.ServiceBase
}

const (
	ServiceName = "metadata"
)

// QutoServiceFactory
type QutoServiceFactory struct{}

// New create metadata service factory
func (p *QutoServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	return &QutoService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
	}, nil

}

// Name
func (p *QutoService) Name() string {
	return ServiceName
}

// Start
func (p *QutoService) Start() error {
	return nil
}

// Stop
func (p *QutoService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

// handleNotifications handle notification from kafka
func (p *QutoService) handleNotifications(topic string, value []byte) error {
	return nil
}
