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

package watcher

import (
	"errors"
	"os"
	"os/exec"
	"text/template"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	sd "github.com/cloustone/sentel/pkg/service-discovery"
	"github.com/golang/glog"
)

type watcherService struct {
	config    config.Config
	discovery sd.ServiceDiscovery
}

// ServiceFactory service factory
type ServiceFactory struct{}

// New create watcher service instance
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	return &watcherService{
		config: c,
	}, nil
}

// Name
func (p *watcherService) Name() string { return "iothub" }

func (p *watcherService) Initialize() error {
	khosts, err := p.config.String("iothub", "zookeeper")
	if err != nil || khosts == "" {
		return errors.New("invalid zookeeper hosts option")
	}
	discovery, derr := sd.New(sd.Option{
		Backend:      sd.BackendZookeeper,
		Hosts:        khosts,
		ServicesPath: "/iotservices",
	})
	if derr != nil {
		return errors.New("service backend initialization failed")
	}
	p.discovery = discovery

	return nil
}

// Start
func (p *watcherService) Start() error {
	var path = p.config.MustString("iothub", "root_path")
	p.discovery.StartWatcher(path, p.handle, nil)
	services := p.discovery.GetServices(path)
	p.handle(services, nil)

	return nil
}

// Stop
func (p *watcherService) Stop() {
	p.discovery.Close()
}

func (p *watcherService) handle(services []sd.Service, ctx interface{}) {
	if tmpl, err := template.ParseFiles("/sentel/template.conf"); err == nil {
		if f, err := os.OpenFile("/etc/nginx/stream_mqtt.conf", os.O_CREATE|os.O_WRONLY, 0666); err == nil {
			defer f.Close()
			if err := tmpl.Execute(f, services); err != nil {
				glog.Errorf("template execute error, %s", err)
				return
			}
			cmd := exec.Command("nginx", "-s", "reload")
			cmd.Start()
			glog.Infof("update nginx successfuly")
			return
		}
	}
	glog.Error("update nginx failed")
}
