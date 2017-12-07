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

package base

import (
	"net"
	"time"

	"github.com/golang/glog"
)

type BrokerStartupInfo struct {
	Id        string    `json:"nodeId"`
	Ip        string    `json:"nodeIp"`
	Version   string    `json:"version"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"nodeStatus"`
	UpdatedAt time.Time `json:"updatedAt"`
}

var brokerStartupInfo BrokerStartupInfo

func GetBrokerId() string                      { return brokerStartupInfo.Id }
func GetBrokerStartupInfo() *BrokerStartupInfo { return &brokerStartupInfo }

func SetBrokerStartupInfo(info *BrokerStartupInfo) {
	brokerStartupInfo.Id = info.Id
	brokerStartupInfo.CreatedAt = time.Now()
	brokerStartupInfo.Status = "started"
	// Get ip address
	if addrs, err := net.InterfaceAddrs(); err == nil {
		for _, address := range addrs {
			if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
				if ipnet.IP.To4() != nil {
					brokerStartupInfo.Ip = ipnet.IP.String()
					break
				}
			}
		}
	}
	if brokerStartupInfo.Ip == "" {
		glog.Error("Failed to get local broker address")
	}
}
