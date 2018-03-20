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

package client

import (
	"sync/atomic"
	"time"
)

type QPSMonitor struct {
	latestIndex  int32
	delaySecond  int32
	totalQueries []int32
}

func (p *QPSMonitor) Pulse() {
	index := int32(time.Now().Second()) % p.delaySecond

	if p.latestIndex != index {
		atomic.StoreInt32(&p.latestIndex, index)
		atomic.StoreInt32(&p.totalQueries[p.latestIndex], 0)
	}

	atomic.AddInt32(&p.totalQueries[index], 1)
}

func (p *QPSMonitor) QPS() int32 {
	var totalCount int32 = 0
	for i, queryCount := range p.totalQueries {
		if int32(i) != p.latestIndex {
			totalCount += queryCount
		}
	}
	return totalCount / (p.delaySecond - 1)
}

func NewQPSMonitor(delaySecond int32) *QPSMonitor {
	if delaySecond < 5 {
		delaySecond = 5
	}
	monitor := QPSMonitor{
		delaySecond:  delaySecond,
		totalQueries: make([]int32, delaySecond),
	}
	return &monitor
}
