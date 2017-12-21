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

	com "github.com/cloustone/sentel/common"
	"github.com/cloustone/sentel/common/db"
)

type DataEndpoint interface {
	Name() string
	Write(data map[string]interface{}) error
}

func NewEndpoint(c com.Config, r *db.Rule) (DataEndpoint, error) {
	switch r.DataTarget.Type {
	case db.DataTargetTypeTopic:
		return newTopicEndpoint(c, r)
	case db.DataTargetTypeOuterDatabase:
		return newOuterdbEndpoint(c, r)
	case db.DataTargetTypeInnerDatabase:
		return newInnerdbEndpoint(c, r)
	case db.DataTargetTypeMessageService:
		return newMsgEndpoint(c, r)
	}
	return nil, fmt.Errorf("data target '%s' is not implemented", r.DataTarget.Type)
}
