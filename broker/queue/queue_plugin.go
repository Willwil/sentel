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

import (
	"time"

	"github.com/cloustone/sentel/core"
)

type queuePlugin interface {
	// getData get a slot of data from backend
	getData() (*queueData, error)
	// pushData push a slot of data to backend
	pushData(*queueData)
}

type queueData struct {
	Id   uint32    `bson:"id"`
	Data []byte    `bson:"data"`
	Time time.Time `bson:"time"`
}

// newPlugin return backend queue plugin accoriding to configuration and id
func newPlugin(id string, c core.Config) (queuePlugin, error) {
	return newMongoQueuePlugin(id, c)

}
