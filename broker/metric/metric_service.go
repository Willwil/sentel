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
	"container/list"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/ceilometer/collector"
	"github.com/cloustone/sentel/core"
)

const ServiceName = "metric"

type MetricService struct {
	base.ServiceBase
	keepalive    *time.Ticker
	stat         *time.Ticker
	name         string
	createdAt    string
	ip           string
	metrics      map[string]*list.List
	metricsMutex sync.Mutex
	stats        map[string]*list.List
	statsMutex   sync.Mutex
}

// MetricServiceFactory
type MetricServiceFactory struct{}

// New create apiService service factory
func (p *MetricServiceFactory) New(c core.Config, quit chan os.Signal) (base.Service, error) {
	// Get node ip, name and created time
	return &MetricService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
		metricsMutex: sync.Mutex{},
		metrics:      make(map[string]*list.List),
		statsMutex:   sync.Mutex{},
		stats:        make(map[string]*list.List),
	}, nil
}

// Name
func (p *MetricService) Name() string {
	return ServiceName
}

func (p *MetricService) Initialize() error { return nil }

// Start
func (p *MetricService) Start() error {
	// Launch timer scheduler
	duration, err := p.Config.Int("mqttbroker", "report_duration")
	if err != nil {
		duration = 2
	}
	p.keepalive = time.NewTicker(1 * time.Second)
	p.stat = time.NewTicker(time.Duration(duration) * time.Second)
	go func(p *MetricService) {
		for {
			p.WaitGroup.Add(1)
			select {
			case <-p.keepalive.C:
				p.reportKeepalive()
			case <-p.stat.C:
				p.reportHubStats()
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *MetricService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.keepalive.Stop()
	p.stat.Stop()
	p.WaitGroup.Wait()
	close(p.Quit)
}

// reportHubStats report current iothub stats
func (p *MetricService) reportHubStats() {
	// Stats
	stats := GetStats("mqtt")
	collector.AsyncReport(p.Config, collector.TopicNameStats,
		&collector.Stats{
			NodeName: p.name,
			Service:  "mqtt",
			Values:   stats,
		})

	// Metrics
	metrics := GetMetric("mqtt")
	collector.AsyncReport(p.Config, collector.TopicNameMetric,
		&collector.Metric{
			NodeName: p.name,
			Service:  "mqtt",
			Values:   metrics,
		})
}

// reportKeepalive report node information to cluster manager
func (p *MetricService) reportKeepalive() {
	/*
		broker := base.GetBroker()
			node := broker.GetNodeInfo()
			collector.AsyncReport(p.Config, collector.TopicNameNode,
				&collector.Node{
					NodeName:  node.NodeName,
					NodeIp:    node.NodeIp,
					CreatedAt: node.CreatedAt,
				})
	*/
}

// newMetrics allocate a metrics object from metric service
func (p *MetricService) newMetric(serviceName string) Metric {
	p.metricsMutex.Lock()
	defer p.metricsMutex.Unlock()
	if _, ok := p.metrics[serviceName]; !ok {
		p.metrics[serviceName] = list.New()
	}
	metric := newMetricWithLock()
	metric.element = p.metrics[serviceName].PushBack(metric)
	return metric
}

// getMetric return server metrics
func (p *MetricService) getMetric(serviceName string) map[string]uint64 {
	result := map[string]uint64{}
	p.metricsMutex.Lock()
	defer p.metricsMutex.Unlock()

	metric := newMetricWithLock()
	if _, ok := p.metrics[serviceName]; ok {
		metrics := p.metrics[serviceName]
		for e := metrics.Front(); e != nil; e = e.Next() {
			metric.AddMetric(e.Value.(Metric))
		}
		result = metric.Get()
	}
	return result
}

// getMetric return server metrics
func (p *MetricService) freeMetric(serviceName string, m Metric) {
	p.metricsMutex.Lock()
	defer p.metricsMutex.Unlock()

	if _, ok := p.metrics[serviceName]; ok {
		mm := m.(*metricWithLock)
		p.metrics[serviceName].Remove(mm.element)
	}
}

// newStats allocate a stats object from metric service
func (p *MetricService) newStats(serviceName string) Metric {
	p.statsMutex.Lock()
	if _, ok := p.stats[serviceName]; !ok {
		p.stats[serviceName] = list.New()
	}
	stat := newMetricWithLock()
	defer p.statsMutex.Unlock()
	p.stats[serviceName].PushBack(stat)
	return stat
}

// getStats return service's stats
func (p *MetricService) getStats(serviceName string) map[string]uint64 {
	result := map[string]uint64{}
	p.statsMutex.Lock()
	defer p.statsMutex.Unlock()

	metric := newMetricWithLock()
	if _, ok := p.stats[serviceName]; ok {
		stats := p.stats[serviceName]
		for e := stats.Front(); e != nil; e = e.Next() {
			metric.AddMetric(e.Value.(Metric))
		}
		result = metric.Get()
	}
	return result

}
