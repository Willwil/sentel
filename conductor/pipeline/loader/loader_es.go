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

package loader

import (
	"github.com/cloustone/sentel/conductor/pipeline/data"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/olivere/elastic"
)

type elasticLoader struct {
	productId string
	client    *elastic.Client
}

func newESLoader(c config.Config) (Loader, error) {
	client, err := elastic.NewClient()
	if err != nil {
		return nil, err
	}
	if _, _, err := client.Ping().Do(); err != nil {
		return nil, err
	}
	productId := c.MustString("productId")
	if exists, err := client.IndexExists(productId).Do(); err != nil {
		return nil, err
	} else if !exists {
		if _, err := client.CreateIndex(productId).Do(); err != nil {
			return nil, err
		}
	}
	return &elasticLoader{
		productId: productId,
		client:    client,
	}, nil
}

func (p *elasticLoader) Name() string { return "elastic" }
func (p *elasticLoader) Close()       {}

func (p *elasticLoader) Load(f *data.DataFrame) error {
	topic := f.Context("topic").(string)
	_, err := p.client.Index().Index(p.productId).Type(topic).BodyJson(f.PrettyJson()).Do()
	return err
}
