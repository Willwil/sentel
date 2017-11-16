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

package metric

import (
	"os"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/ceilometer/collector"
	"github.com/cloustone/sentel/core"
)

type MetricService struct {
	core.ServiceBase
	keepalive *time.Ticker
	stat      *time.Ticker
	name      string
	createdAt string
	ip        string
}

// MetricServiceFactory
type MetricServiceFactory struct{}

// New create apiService service factory
func (this *MetricServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// Get node ip, name and created time
	return &MetricService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
	}, nil
}

// Info
func (this *MetricService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "report-service",
	}
}

// Name
func (this *MetricService) Name() string {
	return "metric"
}

// Start
func (this *MetricService) Start() error {
	// Launch timer scheduler
	duration, err := this.Config.Int("mqttbroker", "report_duration")
	if err != nil {
		duration = 2
	}
	this.keepalive = time.NewTicker(1 * time.Second)
	this.stat = time.NewTicker(time.Duration(duration) * time.Second)
	go func(*MetricService) {
		for {
			select {
			case <-this.keepalive.C:
				this.reportKeepalive()
			case <-this.stat.C:
				this.reportHubStats()
			}
		}
	}(this)
	return nil
}

// Stop
func (this *MetricService) Stop() {
	this.keepalive.Stop()
	this.stat.Stop()
}

// reportHubStats report current iothub stats
func (this *MetricService) reportHubStats() {
	broker := base.GetBroker()
	// Stats
	stats := broker.GetStats("mqtt")
	collector.AsyncReport(this.Config, collector.TopicNameStats,
		&collector.Stats{
			NodeName: this.name,
			Service:  "mqtt",
			Values:   stats,
		})

	// Metrics
	metrics := broker.GetMetrics("mqtt")
	collector.AsyncReport(this.Config, collector.TopicNameMetric,
		&collector.Metric{
			NodeName: this.name,
			Service:  "mqtt",
			Values:   metrics,
		})
}

// reportKeepalive report node information to cluster manager
func (this *MetricService) reportKeepalive() {
	broker := base.GetBroker()
	// Node
	node := broker.GetNodeInfo()
	collector.AsyncReport(this.Config, collector.TopicNameNode,
		&collector.Node{
			NodeName:  node.NodeName,
			NodeIp:    node.NodeIp,
			CreatedAt: node.CreatedAt,
		})
}
