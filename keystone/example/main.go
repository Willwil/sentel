package main

import (
	"flag"

	"github.com/cloustone/sentel/keystone/client"
	l2 "github.com/cloustone/sentel/keystone/l2os"
	"github.com/cloustone/sentel/keystone/ram"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/golang/glog"
)

func main() {
	flag.Parse()
	c := config.New()
	c.AddConfig(map[string]map[string]string{
		"keystone": {
			"hosts": "localhost:4146",
		},
	})

	client, _ := client.New(c)
	glog.Error(client.CreateAccount("account1"))
	glog.Error(client.CreateResource("account1", ram.ResourceCreateOption{
		Name: "product1",
	}))
	glog.Error(client.AddResourceGrantee("product1", "client1", l2.RightRead))
	glog.Error(client.AccessResource("product1", "client1", l2.ActionWrite))
	glog.Error(client.DestroyResource("product1", "account1"))
	glog.Error(client.DestroyAccount("account1"))
}
