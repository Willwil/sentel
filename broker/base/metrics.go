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

package base

import "sync"

// Metrics declarations
type Metrics struct {
	metrics map[string]uint64
	mutex   *sync.Mutex
}

func NewMetrics(withlock bool) *Metrics {
	if withlock {
		return &Metrics{
			metrics: make(map[string]uint64),
			mutex:   &sync.Mutex{},
		}
	} else {
		return &Metrics{
			metrics: make(map[string]uint64),
			mutex:   nil,
		}
	}
}

func (p *Metrics) Get() map[string]uint64 {
	if p.mutex != nil {
		p.mutex.Lock()
		defer p.mutex.Unlock()
	}
	return p.metrics
}
func (p *Metrics) addMetric(name string, value uint64) {
	if p.mutex != nil {
		p.mutex.Lock()
		defer p.mutex.Unlock()
	}
	p.metrics[name] += value
}

func (p *Metrics) AddMetrics(metrics *Metrics) {
	if p.mutex != nil {
		p.mutex.Lock()
		defer p.mutex.Unlock()
	}
	for k, v := range metrics.Get() {
		if _, ok := p.metrics[k]; !ok {
			p.metrics[k] = v
		} else {
			p.metrics[k] += v
		}
	}
}
