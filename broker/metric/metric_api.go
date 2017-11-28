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

import "github.com/cloustone/sentel/broker/base"

// NewMetrics allocate a metrics object from metric service
func NewMetric(serviceName string) Metric {
	metric := base.GetService(ServiceName).(*metricService)
	return metric.newMetric(serviceName)
}

// GetMetrics return server metrics
func GetMetric(serviceName string) map[string]uint64 {
	metric := base.GetService(ServiceName).(*metricService)
	return metric.getMetric(serviceName)
}

// FreeMetric delete metric from metric service
func FreeMetric(serviceName string, m Metric) {
	metric := base.GetService(ServiceName).(*metricService)
	metric.freeMetric(serviceName, m)
}

// NewMetrics allocate a metrics object from metric service
func NewStats(serviceName string) Metric {
	metric := base.GetService(ServiceName).(*metricService)
	return metric.newStats(serviceName)
}

// GetStats return server's stats
func GetStats(serviceName string) map[string]uint64 {
	metric := base.GetService(ServiceName).(*metricService)
	return metric.getStats(serviceName)
}
