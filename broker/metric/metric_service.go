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
	"strings"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/meter/collector"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"
)

const ServiceName = "metric"

type metricService struct {
	config       config.Config
	waitgroup    sync.WaitGroup
	quitChan     chan interface{}
	aliveTimer   *time.Timer
	metricTimer  *time.Timer
	name         string
	createdAt    string
	ip           string
	metrics      map[string]*list.List
	metricsMutex sync.Mutex
	stats        map[string]*list.List
	statsMutex   sync.Mutex
}

type ServiceFactory struct{}

// New create apiService service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// Get node ip, name and created time
	return &metricService{
		config:       c,
		quitChan:     make(chan interface{}),
		waitgroup:    sync.WaitGroup{},
		metricsMutex: sync.Mutex{},
		metrics:      make(map[string]*list.List),
		statsMutex:   sync.Mutex{},
		stats:        make(map[string]*list.List),
	}, nil
}

// Name
func (p *metricService) Name() string {
	return ServiceName
}

func (p *metricService) Initialize() error { return nil }

// Start
func (p *metricService) Start() error {
	services, err := p.config.String(ServiceName, "services")
	// If no services are specified, just return simpily
	if err != nil || services == "" {
		glog.Errorf("No metric services are specified")
		return nil
	}
	glog.Infof("metrics for '%s' is started", services)

	t1 := p.config.MustInt(ServiceName, "report_duration")
	t2 := p.config.MustInt(ServiceName, "keepalive")
	p.metricTimer = time.NewTimer(time.Duration(t1) * time.Second)
	p.aliveTimer = time.NewTimer(time.Duration(t2) * time.Second)
	p.waitgroup.Add(1)
	go func(p *metricService) {
		defer p.waitgroup.Done()
		for {
			select {
			case <-p.aliveTimer.C:
				p.reportKeepalive()
			case <-p.metricTimer.C:
				p.reportMetric()
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *metricService) Stop() {
	p.quitChan <- true
	if p.aliveTimer != nil {
		p.aliveTimer.Stop()
	}
	if p.metricTimer != nil {
		p.metricTimer.Stop()
	}
	p.waitgroup.Wait()
	close(p.quitChan)
}

// reportHubStats report current iothub stats
func (p *metricService) reportMetric() {
	val := p.config.MustString(ServiceName, "services")
	services := strings.Split(val, ",")
	for _, service := range services {
		metrics := GetMetric(service)
		collector.AsyncReport(p.config,
			&collector.Metric{
				TopicName: collector.TopicNameMetric,
				NodeName:  p.name,
				Service:   service,
				Values:    metrics,
			})
	}
}

// reportKeepalive report node information to cluster manager
func (p *metricService) reportKeepalive() {
	info := base.GetBrokerStartupInfo()
	collector.AsyncReport(p.config,
		&collector.Node{
			TopicName:  collector.TopicNameNode,
			NodeId:     info.Id,
			NodeIp:     info.Ip,
			CreatedAt:  info.CreatedAt,
			UpdatedAt:  time.Now(),
			NodeStatus: "started",
		})
}

// newMetrics allocate a metrics object from metric service
func (p *metricService) newMetric(serviceName string) Metric {
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
func (p *metricService) getMetric(serviceName string) map[string]uint64 {
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
func (p *metricService) freeMetric(serviceName string, m Metric) {
	p.metricsMutex.Lock()
	defer p.metricsMutex.Unlock()

	if _, ok := p.metrics[serviceName]; ok {
		mm := m.(*metricWithLock)
		p.metrics[serviceName].Remove(mm.element)
	}
}

// newStats allocate a stats object from metric service
func (p *metricService) newStats(serviceName string) Metric {
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
func (p *metricService) getStats(serviceName string) map[string]uint64 {
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
