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

package sd

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"

	"github.com/samuel/go-zookeeper/zk"
)

type serviceDisZK struct {
	config    config.Config
	conn      *zk.Conn
	waitgroup sync.WaitGroup
	quitChan  chan interface{}
	handler   WatcherHandlerFunc
}

func newServiceDiscoveryZK(c config.Config) (ServiceDiscovery, error) {
	hosts, err := c.String("service-discovery", "hosts")
	if err != nil {
		return nil, errors.New("no zookeeper hosts")
	}
	conn, _, err := zk.Connect(strings.Split(hosts, ","), time.Second*2)
	if err != nil {
		return nil, fmt.Errorf("service discovery can not connect with zk:%s", hosts)
	}
	return &serviceDisZK{
		config:    c,
		waitgroup: sync.WaitGroup{},
		quitChan:  make(chan interface{}),
		conn:      conn,
	}, nil
}

func (p *serviceDisZK) RegisterService(s Service) error {
	path := fmt.Sprintf("/iotservices/%s", s.Name)
	buf, err := json.Marshal(&s)
	if err != nil {
		return fmt.Errorf("service '%s' data marshal failed", s.Name)
	}
	acls := zk.WorldACL(zk.PermAll)
	_, err = p.conn.Create(path, buf, 0, acls)
	return err
}

func (p *serviceDisZK) RemoveService(s Service) {
	path := fmt.Sprintf("/iotservices/%s", s.Name)
	p.conn.Delete(path, 0)

}
func (p *serviceDisZK) UpdateService(s Service) error {
	return nil
}

func (p *serviceDisZK) StartWatcher(rootPath string, handler WatcherHandlerFunc, ctx interface{}) error {
	if found, _, err := p.conn.Exists(rootPath); err != nil {
		return err
	} else if !found {
		if path, err := p.conn.Create(rootPath, nil, 0, zk.WorldACL(zk.PermAll)); err != nil {
			return err
		} else if rootPath != path {
			return fmt.Errorf("different path found, '%s' != '%s'", rootPath, path)
		}
	}
	p.waitgroup.Add(1)
	p.handler = handler
	go func(p *serviceDisZK) {
		defer p.waitgroup.Done()
		for {
			_, _, childCh, err := p.conn.ChildrenW(rootPath)
			if err != nil {
				glog.Errorf("service-discovery children error, %+v", err)
				continue
			}
			select {
			case <-childCh:
				services := p.GetServices(rootPath)
				if len(services) != 0 {
					p.handler(services, ctx)
				}
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

func (p *serviceDisZK) GetServices(rootPath string) []Service {
	services := []Service{}
	children, _, err := p.conn.Children(rootPath)
	if err != nil {
		return services
	}
	for _, child := range children {
		if data, _, err := p.conn.Get(rootPath + "/" + child); err == nil {
			var s Service
			err := json.Unmarshal(data, &s)
			if err != nil {
				glog.Errorf("json unmarshal error: data %+v, err %+v", data, err)
			}
			services = append(services, s)
		}
	}
	return services
}

func (p *serviceDisZK) Close() {
	if p.handler != nil {
		p.quitChan <- true
		p.waitgroup.Wait()
	}
	p.conn.Close()
}
