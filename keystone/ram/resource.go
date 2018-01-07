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

func (p *Resource) Access(accessorId string, action l2.Action) error {
	for _, g := range p.GranteeList {
		if g.AccessorId == accessorId {
			switch action {
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

func (p *Resource) AccessAttribute(attr string, accessorId string, action l2.Action) error {
	if _, found := p.Attributes[attr]; !found {
		return fmt.Errorf("invalid resource '%s' attribute '%s'", p.Name, attr)
	}
	granteeList := p.Attributes[attr]
	for _, g := range granteeList {
		if g.AccessorId == accessorId {
			switch action {
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
func (p *Resource) AddGrantee(accessorId string, right l2.Right) error {
	p.Object.AddGrantee(l2.Grantee{AccessorId: accessorId, Right: right})
	return l2api.UpdateObject(&p.Object)
}

func (p *Resource) AddAttributeGrantee(attr string, accessorId string, right l2.Right) error {
	p.Object.AddAttributeGrantee(attr, l2.Grantee{AccessorId: accessorId, Right: right})
	return l2api.UpdateObject(&p.Object)
}

func (p *Resource) AddAttribute(attr string, granteeList []l2.Grantee) error {
	p.Object.AddAttribute(attr, granteeList)
	return l2api.UpdateObject(&p.Object)
}

func (p *Resource) RemoveAttribute(attr string) error {
	p.Object.RemoveAttribute(attr)
	return l2api.UpdateObject(&p.Object)
}

func (p *Resource) RemoveAttributeGrantee(attr string, accessorId string) error {
	p.Object.RemoveAttributeGrantee(attr, accessorId)
	return l2api.UpdateObject(&p.Object)
}
