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
	"keystone": {
		"listen":          "localhost:4147",
		"loglevel":        "debug",
		"version":         "v1",
		"auth":            "none",
		"mongo":           "localhost:27017",
		"connect_timeout": 2,
		"storage":         "mongo",
	},
	"registry": {
		"hosts":    "localhost:27017",
		"loglevel": "debug",
	},
	"security": {
		"cafile":              "",
		"capath":              "",
		"certfile":            "",
		"keyfile":             "",
		"require_certificate": false,
	},
}
