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

import (
	"fmt"
	"net"
	"os"
)

const (
	ServiceNameMongo   = "mongo"
	ServiceNameRedis   = "redis"
	ServiceNameKafka   = "kafka"
	SentelEnvKafkaHost = "SENTEL_KAFKA_HOST"
	SentelEnvRedisHost = "SENTEL_REDIS_HOST"
	SentelEnvMongoHost = "SENTEL_MONGO_HOST"
)

// GetServiceEndpoint return service address and port
// It will lookup cluster manager(k8s) to find endpoint, if failed,
// get endpoint from local configuration

func GetServiceEndpoint(c Config, serviceName string, endpointName string) (string, error) {
	// trieve endpoint from environment variable at first
	switch endpointName {
	case ServiceNameMongo:
		if v := os.Getenv(SentelEnvMongoHost); v != "" {
			return v, nil
		}
	case ServiceNameRedis:
		if v := os.Getenv(SentelEnvRedisHost); v != "" {
			return v, nil
		}
	case ServiceNameKafka:
		if v := os.Getenv(SentelEnvKafkaHost); v != "" {
			addr, err := net.ResolveTCPAddr("tcp", v)
			if err == nil && addr != nil {
				return addr.String(), nil
			}
		}
	}
	// Now, just return from local configurations
	result, err := c.String(serviceName, endpointName)
	if err != nil || result == "" {
		return "", fmt.Errorf("No service endpoint:'%s' found", endpointName)
	}
	return result, nil
}
