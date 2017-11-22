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

package quto

import "github.com/cloustone/sentel/broker/base"

// GetQuto return object's qutotation
func GetQuto(target, id string) (*Quto, error) {
	quto := base.GetService(ServiceName).(*QutoService)
	return quto.getQuto(target, id)
}

// GetQutoItem return object's qutotation
func GetQutoItem(target, id, item string) uint64 {
	quto := base.GetService(ServiceName).(*QutoService)
	return quto.getQutoItem(target, id, item)
}

// SetQuto set object's qutotation
func SetQuto(target, id string, q *Quto) {
	quto := base.GetService(ServiceName).(*QutoService)
	quto.setQuto(target, id, q)
}

// SetQutoItem set object item's qutotations
func SetQutoItem(target, id, item string, value uint64) {
	quto := base.GetService(ServiceName).(*QutoService)
	quto.setQutoItem(target, id, item, value)
}

// AddQutoItemValue add quto items's value
func AddQutoItem(target, id, item string, value uint64) {
	quto := base.GetService(ServiceName).(*QutoService)
	quto.addQutoItem(target, id, item, value)
}

// SubQutoItemValue sub quto items's value
func SubQutoItem(target, id, item string, value uint64) {
	quto := base.GetService(ServiceName).(*QutoService)
	quto.subQutoItem(target, id, item, value)
}
