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
	"testing"
	"time"

	l2 "github.com/cloustone/sentel/keystone/l2os"
	"github.com/cloustone/sentel/pkg/config"
)

var (
	rid1      = NewObjectId()
	resource1 = &Resource{
		Object: l2.Object{
			Name:        "resource1",
			ObjectId:    rid1,
			CreatedAt:   time.Now(),
			Attributes:  make(map[string][]l2.Grantee),
			Creator:     "account1",
			Category:    "product",
			GranteeList: []l2.Grantee{},
		},
	}
)

func TestResource_initialize(t *testing.T) {
	c := config.New("keystone")
	c.AddConfig(config.M{
		"keystone": {
			"hosts":           "localhost:4147",
			"mongo":           "localhost:27017",
			"connect_timeout": 2,
		},
	})
	if err := Initialize(c, "direct"); err != nil {
		t.Error(err)
	}
}

func TestResource_AddGrantee(t *testing.T) {
	resource1.AddGrantee("client1", RightRead)
	resource1.AddGrantee("client2", RightWrite)
	grantees := resource1.GranteeList
	if len(grantees) != 2 {
		t.Error("grantee list wrong")
	}
	if grantees[0].AccessorId != "client1" && grantees[1].AccessorId != "client2" {
		t.Error("accessorid is wrong")
	}
	if grantees[1].Right != convertRight(RightRead) &&
		grantees[1].Right != convertRight(RightWrite) {
		t.Error("accessor right is rong")
	}
}
func TestResource_AddAttribute(t *testing.T) {
	resource1.AddAttribute("attr1")
	resource1.AddAttribute("attr2")
	if _, err := resource1.GetAttribute("attr1"); err != nil {
		t.Error(err)
	}
	if _, err := resource1.GetAttribute("attr2"); err != nil {
		t.Error(err)
	}
	resource1.RemoveAttribute("attr1")
	resource1.RemoveAttribute("attr2")
	if _, err := resource1.GetAttribute("attr1"); err == nil {
		t.Error("remove an unexisted attribute failure")
	}
}

func TestResource_AddAttributeGrantee(t *testing.T) {
	resource1.AddAttribute("attr1")
	resource1.AddAttribute("attr2")
	resource1.AddAttributeGrantee("attr1", "client1", RightRead)
	resource1.AddAttributeGrantee("attr2", "client1", RightWrite)
	if list, err := resource1.GetAttributeGranteeList("attr1"); err != nil {
		t.Error(err)
	} else if len(list) != 1 || list[0].AccessorId != "client1" || list[0].Right != convertRight(RightRead) {
		t.Error("get attribute grantee list failure")
	}
	resource1.RemoveAttribute("attr1")
	resource1.RemoveAttribute("attr2")
}

func TestResource_Access(t *testing.T) {
	resource1.AddAttribute("attr1")
	resource1.AddAttribute("attr2")
	resource1.AddGrantee("client1", RightRead)
	resource1.AddGrantee("client2", RightWrite)
	resource1.AddAttributeGrantee("attr1", "client1", RightRead)
	resource1.AddAttributeGrantee("attr2", "client2", RightWrite)
	if err := resource1.Access("client1", ActionRead); err != nil {
		t.Error(err)
	}
	if err := resource1.Access("client1", ActionWrite); err == nil {
		t.Error("resource authorize failure")
	}
	if err := resource1.Access("client2", ActionRead); err != nil {
		t.Error(err)
	}
	if err := resource1.Access("client2", ActionWrite); err != nil {
		t.Error(err)
	}
	if err := resource1.AccessAttribute("attr1", "client1", ActionRead); err != nil {
		t.Error(err)
	}
	if err := resource1.AccessAttribute("attr1", "client1", ActionWrite); err == nil {
		t.Error("resource attribute authorize failure")
	}
	if err := resource1.AccessAttribute("attr2", "client2", ActionRead); err != nil {
		t.Error(err)
	}
	if err := resource1.AccessAttribute("attr2", "client2", ActionWrite); err != nil {
		t.Error(err)
	}
	resource1.RemoveAttribute("attr1")
	resource1.RemoveAttribute("attr2")
}
