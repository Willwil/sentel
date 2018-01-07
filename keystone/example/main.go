package main

import (
	"flag"

	"github.com/cloustone/sentel/keystone/client"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	c := config.New()
	c.AddConfig(map[string]map[string]string{
		"keystone": {
			"hosts": "localhost:4147",
		},
	})
	client.Initialize(c)
	glog.Error(client.CreateAccount("account1"))
	rid := ram.NewObjectId()
	glog.Error(client.CreateResource("account1", ram.ResourceCreateOption{
		ObjectId:   rid,
		Name:       "product1",
		Attributes: []string{"devices", "rules"},
	}))
	glog.Error(client.AddResourceGrantee(rid, "client1", ram.RightRead))
	glog.Error(client.AddResourceGrantee(rid+"/devices", "client1", ram.RightRead))
	glog.Error(client.AddResourceGrantee(rid+"/devices", "client1", ram.RightWrite))
	glog.Error(client.Authorize("client1", rid, ram.ActionRead))
	glog.Error(client.Authorize("client1", rid+"/devices", ram.ActionRead))
	glog.Error(client.DestroyResource(rid, "account1"))
	glog.Error(client.DestroyAccount("account1"))
}
