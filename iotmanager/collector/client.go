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

	"github.com/cloustone/sentel/iotmanager/mgrdb"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/message"
)

// ClientTopic
type ClientTopic struct {
	TopicName string
	mgrdb.Client
}

func (p *ClientTopic) Topic() string        { return TopicNameClient }
func (p *ClientTopic) SetTopic(name string) {}
func (p *ClientTopic) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *ClientTopic) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *ClientTopic) handleTopic(c config.Config, ctx context.Context) error {
	db, err := mgrdb.New(c)
	if err != nil {
		return err
	}
	defer db.Close()
	return db.UpdateClient(p.Client)

}
