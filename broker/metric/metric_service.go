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
func (p *MetricServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
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
func (p *MetricService) Info() *base.ServiceInfo {
	return &base.ServiceInfo{
		ServiceName: "report-service",
	}
}

// Name
func (p *MetricService) Name() string {
	return "metric"
}

// Start
func (p *MetricService) Start() error {
	// Launch timer scheduler
	duration, err := p.Config.Int("mqttbroker", "report_duration")
	if err != nil {
		duration = 2
	}
	p.keepalive = time.NewTicker(1 * time.Second)
	p.stat = time.NewTicker(time.Duration(duration) * time.Second)
	go func(*MetricService) {
		for {
			select {
			case <-p.keepalive.C:
				p.reportKeepalive()
			case <-p.stat.C:
				p.reportHubStats()
			}
		}
	}(p)
	return nil
}

// Stop
func (p *MetricService) Stop() {
	p.keepalive.Stop()
	p.stat.Stop()
}

// reportHubStats report current iothub stats
func (p *MetricService) reportHubStats() {
	broker := base.GetBroker()
	// Stats
	stats := broker.GetStats("mqtt")
	collector.AsyncReport(p.Config, collector.TopicNameStats,
		&collector.Stats{
			NodeName: p.name,
			Service:  "mqtt",
			Values:   stats,
		})

	// Metrics
	metrics := broker.GetMetrics("mqtt")
	collector.AsyncReport(p.Config, collector.TopicNameMetric,
		&collector.Metric{
			NodeName: p.name,
			Service:  "mqtt",
			Values:   metrics,
		})
}

// reportKeepalive report node information to cluster manager
func (p *MetricService) reportKeepalive() {
	broker := base.GetBroker()
	// Node
	node := broker.GetNodeInfo()
	collector.AsyncReport(p.Config, collector.TopicNameNode,
		&collector.Node{
			NodeName:  node.NodeName,
			NodeIp:    node.NodeIp,
			CreatedAt: node.CreatedAt,
		})
}
