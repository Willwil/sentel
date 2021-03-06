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
	"mns": {
		"loglevel":        "debug",
		"kafka":           "localhost:9092",
		"mongo":           "localhost:27017",
		"connect_timeout": 5,
		"listen":          ":50061",
	},
}
