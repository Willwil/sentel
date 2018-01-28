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

package collector

import (
	"context"
	"time"

	mgo "gopkg.in/mgo.v2"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

const (
	TopicNameNode         = "/cluster/nodes"
	TopicNameClient       = "/cluster/clients"
	TopicNameSession      = "/cluster/sessions"
	TopicNameSubscription = "/cluster/subscriptions"
	TopicNamePublish      = "/cluster/publish"
	TopicNameMetric       = "/cluster/metrics"
	TopicNameStats        = "/cluster/stats"
)

const (
	ObjectActionRegister   = "register"
	ObjectActionUnregister = "unregister"
	ObjectActionRetrieve   = "retrieve"
	ObjectActionDelete     = "delete"
	ObjectActionUpdate     = "update"
)

type topicHandler interface {
	handleTopic(c config.Config, ctx context.Context) error
}

func SyncReport(c config.Config, msg message.Message) error {
	return message.PostMessage(c, msg)
}

func AsyncReport(c config.Config, msg message.Message) error {
	return message.PostMessage(c, msg)
}

func getDatabase(c config.Config) (*mgo.Database, error) {
	hosts := c.MustString("mongo")
	timeout := c.MustInt("connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	session.SetMode(mgo.Monotonic, true)
	return session.DB("iotmanager"), nil
}
