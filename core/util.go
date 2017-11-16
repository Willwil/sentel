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

package core

import "fmt"

// GetServiceEndpoint return service address and port
// It will lookup cluster manager(k8s) to find endpoint, if failed,
// get endpoint from local configuration
func GetServiceEndpoint(c Config, serviceName string, endpointName string) (string, error) {
	// Now, just return from local configurations
	result, err := c.String(serviceName, endpointName)
	if err != nil || result == "" {
		return "", fmt.Errorf("No service endpoint:'%s' found", endpointName)
	}
	return result, nil
}
