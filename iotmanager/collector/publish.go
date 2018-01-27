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

	"github.com/cloustone/sentel/pkg/message"
)

// Publish
type Publish struct {
	TopicName       string
	ClientId        string `json:"clientId"`
	SubscribedTopic string `json:"topic"`
	ProductId       string `json:"product"`
}

func (p *Publish) Topic() string        { return TopicNamePublish }
func (p *Publish) SetTopic(name string) {}
func (p *Publish) Serialize(opt message.SerializeOption) ([]byte, error) {
	return message.Serialize(p, opt)
}
func (p *Publish) Deserialize(buf []byte, opt message.SerializeOption) error { return nil }

func (p *Publish) handleTopic(service *collectorService, ctx context.Context) error {
	return nil
}
