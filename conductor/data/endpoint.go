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

package data

import (
	"fmt"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
)

type DataEndpoint interface {
	Name() string
	Write(data map[string]interface{}) error
}

func NewEndpoint(c config.Config, r *registry.Rule) (DataEndpoint, error) {
	switch r.DataTarget.Type {
	case registry.DataTargetTypeTopic:
		return newTopicEndpoint(c, r)
	case registry.DataTargetTypeOuterDatabase:
		return newOuterdbEndpoint(c, r)
	case registry.DataTargetTypeInnerDatabase:
		return newInnerdbEndpoint(c, r)
	case registry.DataTargetTypeMessageService:
		return newMsgEndpoint(c, r)
	}
	return nil, fmt.Errorf("data target '%s' is not implemented", r.DataTarget.Type)
}
