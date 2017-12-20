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
func GetQuto(id string) (uint64, error) {
	quto := base.GetService(ServiceName).(*qutoService)
	return quto.getQuto(id)
}

// CheckQutoWithAddValue check wether the newly added value is over quto
func CheckQuto(id string, value uint64) bool {
	quto := base.GetService(ServiceName).(*qutoService)
	return quto.checkQuto(id, value)
}
