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
)

type sqldbLoader struct {
	config config.Config
}

func newSqldbLoader(c config.Config) (Loader, error) {
	return &sqldbLoader{config: c}, nil
}

func (p *sqldbLoader) Name() string { return "elastic" }
func (p *sqldbLoader) Close()       {}

func (p *sqldbLoader) Load(f *data.DataFrame) error {
	return nil
}
