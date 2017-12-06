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
	"container/list"
	"sync"
)

// Metrics declarations
type metricWithLock struct {
	element *list.Element
	metrics map[string]uint64
	mutex   *sync.Mutex
}

func newMetricWithLock() *metricWithLock {
	return &metricWithLock{
		metrics: make(map[string]uint64),
		mutex:   &sync.Mutex{},
		element: nil,
	}
}

func (p *metricWithLock) Get() map[string]uint64 { return p.metrics }
func (p *metricWithLock) Add(name string, value uint64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.metrics[name] += value
}
func (p *metricWithLock) Sub(name string, value uint64) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	p.metrics[name] -= value
}

func (p *metricWithLock) AddMetric(metric Metric) {
	p.mutex.Lock()
	defer p.mutex.Unlock()
	for k, v := range metric.Get() {
		if _, ok := p.metrics[k]; !ok {
			p.metrics[k] = v
		} else {
			p.metrics[k] += v
		}
	}
}
