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

package queue

import "github.com/cloustone/sentel/broker/base"

// Allocate queue from queue service
func NewQueue(id string, persistent bool, o Observer) (Queue, error) {
	service := base.GetService(ServiceName).(*queueService)
	return service.newQueue(id, persistent, o)
}

// GetQueue return queue by queue id
func GetQueue(id string) Queue {
	service := base.GetService(ServiceName).(*queueService)
	return service.getQueue(id)
}

// FreeQueue release queue from queue service
func DestroyQueue(id string) {
	service := base.GetService(ServiceName).(*queueService)
	service.destroyQueue(id)
}

// ReleaseQueue decrease queue's reference count, and destory the queue if reference is zero
func ReleaseQueue(id string) {
	service := base.GetService(ServiceName).(*queueService)
	service.releaseQueue(id)
}
