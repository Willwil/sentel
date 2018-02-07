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
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/hailocab/gocassa"
)

type cassLoader struct {
	productId string
	target    registry.RuleDataTarget
	keyspace  gocassa.KeySpace
}

func newCassLoader(c config.Config) (Loader, error) {
	target := c.MustValue("datatarget").(registry.RuleDataTarget)
	productId, _ := c.String("productId")
	hosts, _ := c.String("cassandra")

	keyspace, err := gocassa.ConnectToKeySpace(productId, []string{hosts}, target.Username, target.Password)
	if err != nil {
		return nil, err
	}
	return &cassLoader{
		productId: productId,
		target:    target,
		keyspace:  keyspace,
	}, nil
}

func (p *cassLoader) Name() string { return "elastic" }
func (p *cassLoader) Close()       {}

func (p *cassLoader) Load(f *data.DataFrame) error {
	ftype := f.PrettyJson()
	table := p.keyspace.Table(p.productId, ftype, gocassa.Keys{
		PartitionKeys: []string{"Id"},
	})
	err := table.Create()
	if err == nil {
		return table.Set(ftype).Run()
	}
	return err
}
