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

package main

import "github.com/cloustone/sentel/pkg/config"

var defaultConfigs = config.M{
	"apiserver": {
		"loglevel": "debug",
		"kafka":    "localhost:9092",
		"version":  "v1",
		"auth":     "jwt",
		"mongo":    "localhost:27017",
		"swagger":  ":53384",
	},
	"security": {
		"cafile":              "",
		"capath":              "",
		"certfile":            "",
		"keyfile":             "",
		"require_certificate": false,
	},
	"console": {
		"listen": "0.0.0.0:4145",
	},
	"management": {
		"listen": "0.0.0.0:4146",
	},
	"swagger": {
		"open_browser": true,
	},
	"security_manager": {
		"realms": "registry,sentel",
	},
}
