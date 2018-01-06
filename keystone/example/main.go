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
	glog.Error(client.CreateResource("account1", ram.ResourceCreateOption{
		Name:       "product1",
		Attributes: []string{"devices", "rules"},
	}))
	glog.Error(client.AddResourceGrantee("product1", "client1", ram.RightRead))
	glog.Error(client.AddResourceGrantee("product1/devices", "client1", ram.RightRead))
	glog.Error(client.AddResourceGrantee("product1/devices", "client1", ram.RightWrite))
	glog.Error(client.Authorize("product1", "client1", "r"))
	glog.Error(client.Authorize("product1/devices", "client1", "r"))
	glog.Error(client.DestroyResource("product1", "account1"))
	glog.Error(client.DestroyAccount("account1"))
}
