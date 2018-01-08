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
	"fmt"
	"strings"
	"time"

	l2 "github.com/cloustone/sentel/keystone/l2os"
	"github.com/cloustone/sentel/pkg/config"
)

const (
	RightRead   = "r"
	RightWrite  = "w"
	RightFull   = "x"
	ActionRead  = "r"
	ActionWrite = "w"
	ActionFull  = "x"
)

var (
	l2api l2.Api
)

func NewObjectId() string {
	return l2.NewObjectId()
}

func Initialize(c config.Config, apiName string) error {
	api, err := l2.NewApi(apiName, c)
	l2api = api
	return err
}

func convertAction(action string) l2.Action {
	switch action {
	case ActionRead:
		return l2.ActionRead
	case ActionWrite:
		return l2.ActionWrite
	case ActionFull:
		return l2.ActionFull
	default:
		return l2.ActionUnknown
	}
}

func convertRight(right string) l2.Right {
	switch right {
	case RightRead:
		return l2.RightRead
	case RightWrite:
		return l2.RightWrite | l2.RightRead
	case RightFull:
		return l2.RightFull | l2.RightWrite
	default:
		return l2.RightUnknown
	}
}

func Authorize(resource string, accessorId string, action string) error {
	names := strings.Split(resource, "/")
	if len(names) == 0 || len(names) > 2 {
		return fmt.Errorf("invalid resource '%s'", resource)
	}
	rid := names[0]
	if resource, err := GetResource(rid); err != nil {
		return err
	} else {
		if len(names) == 2 {
			return resource.AccessAttribute(names[1], accessorId, action)
		} else {
			return resource.Access(accessorId, action)
		}
	}
}

// Account
func CreateAccount(name string) error {
	return l2api.CreateAccount(name)
}

func DestroyAccount(name string) error {
	return l2api.DestroyAccount(name)
}

// Resource
func CreateResource(account string, opt *ResourceCreateOption) (*Resource, error) {
	if opt.ObjectId == "" {
		opt.ObjectId = l2.NewObjectId()
	}
	r := &Resource{
		Object: l2.Object{
			Name:        opt.Name,
			ObjectId:    opt.ObjectId,
			CreatedAt:   time.Now(),
			Attributes:  make(map[string][]l2.Grantee),
			Creator:     account,
			Category:    opt.Category,
			GranteeList: []l2.Grantee{},
		},
	}
	// check attribute exist
	for _, attr := range opt.Attributes {
		r.Attributes[attr] = []l2.Grantee{}
	}
	err := l2api.CreateObject(&r.Object)
	return r, err
}

func GetResource(rid string) (*Resource, error) {
	if obj, err := l2api.GetObject(rid); err == nil {
		return &Resource{Object: *obj}, nil
	} else {
		return nil, err
	}
}

func AddResourceGrantee(resource string, accessorId string, right string) error {
	names := strings.Split(resource, "/")
	if len(names) == 0 || len(names) > 2 {
		return fmt.Errorf("invalid resource '%s'", resource)
	}
	rid := names[0]
	if resource, err := GetResource(rid); err != nil {
		return err
	} else {
		if len(names) == 2 {
			resource.AddAttributeGrantee(names[1], accessorId, right)
		} else {
			resource.AddGrantee(accessorId, right)
		}
		return l2api.UpdateObject(&resource.Object)
	}
}

func DestroyResource(rid string) error {
	return l2api.DestroyObject(rid)
}