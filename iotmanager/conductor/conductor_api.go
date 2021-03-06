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

package conductor

import "github.com/cloustone/sentel/pkg/service"

// getConductor return conductor service instance
func getConductor() *conductorService {
	mgr := service.GetServiceManager()
	return mgr.GetService(SERVICE_NAME).(*conductorService)
}

// createTenant add tenant
func CreateTenant(tid string) error {
	c := getConductor()
	return c.createTenant(tid)
}

// RemoveTenant remove tenant from iotcluster
func RemoveTenant(tid string) error {
	c := getConductor()
	return c.removeTenant(tid)
}

// CreateProduct add product to iot cluster
func CreateProduct(tid, pid string, replicas int32) (string, error) {
	c := getConductor()
	return c.createProduct(tid, pid, replicas)
}

// RemoveProduct delete product from iot cluster
func RemoveProduct(tid string, pid string) error {
	c := getConductor()
	return c.removeProduct(tid, pid)
}
