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
package ram

import (
	"errors"
	"fmt"

	l2 "github.com/cloustone/sentel/keystone/l2os"
)

type ResourceCreateOption struct {
	Name       string   `json:"name"`
	ObjectId   string   `json:"objectId"`
	Creator    string   `json:"creator"`
	Category   string   `json:"category"`
	Attributes []string `json:"attributes"`
}

type ResourceDestroyOption struct {
	ObjectId string `json:"fullname"`
	AccessId string `json:"accessId"`
}

type Resource struct {
	l2.Object
}

func (p *Resource) Access(accessorId string, action string) error {
	act := convertAction(action)
	for _, g := range p.GranteeList {
		if g.AccessorId == accessorId {
			switch act {
			case l2.ActionRead:
				if g.Right&l2.RightRead != 0 {
					return nil
				}
			case l2.ActionWrite:
				if g.Right&l2.RightWrite != 0 {
					return nil
				}
			case l2.ActionFull:
				if g.Right&l2.RightFull != 0 {
					return nil
				}
			default:
				return errors.New("invalid action")
			}
		}
	}
	return errors.New("unauthorized")
}

func (p *Resource) AccessAttribute(attr string, accessorId string, action string) error {
	if _, found := p.Attributes[attr]; !found {
		return fmt.Errorf("invalid resource '%s' attribute '%s'", p.Name, attr)
	}
	act := convertAction(action)
	granteeList := p.Attributes[attr]
	for _, g := range granteeList {
		if g.AccessorId == accessorId {
			switch act {
			case l2.ActionRead:
				if g.Right&l2.RightRead != 0 {
					return nil
				}
			case l2.ActionWrite:
				if g.Right&l2.RightWrite != 0 {
					return nil
				}
			case l2.ActionFull:
				if g.Right&l2.RightFull != 0 {
					return nil
				}
			default:
				return errors.New("invalid action")
			}
		}
	}
	return errors.New("unauthorized")
}

// AddAccessor add resource accessor
func (p *Resource) AddGrantee(accessorId string, right string) {
	p.Object.AddGrantee(l2.Grantee{AccessorId: accessorId, Right: convertRight(right)})
}

func (p *Resource) AddAttributeGrantee(attr string, accessorId string, right string) {
	p.Object.AddAttributeGrantee(attr, l2.Grantee{AccessorId: accessorId, Right: convertRight(right)})
}
