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
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/docker-service"

	"github.com/samuel/go-zookeeper/zk"
)

type serviceDisZK struct {
	conn *zk.Conn
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
	return &serviceDisZK{conn: conn}, nil
}

func (p *serviceDisZK) RegisterService(s ds.Service) error {
	path := fmt.Sprintf("/iotservices/%s", s.Name)
	buf, err := json.Marshal(&s)
	if err != nil {
		return fmt.Errorf("service '%s' data marshal failed", s.Name)
	}
	acls := zk.WorldACL(zk.PermAll)
	_, err = p.conn.Create(path, buf, 0, acls)
	return err
}

func (p *serviceDisZK) RemoveService(s ds.Service) {
	path := fmt.Sprintf("/iotservices/%s", s.Name)
	p.conn.Delete(path, 0)

}
func (p *serviceDisZK) UpdateService(s ds.Service) error {
	return nil
}

func (p *serviceDisZK) Close() { p.conn.Close() }
