package main

import (
	"flag"

	"github.com/golang/glog"

	"github.com/cloustone/sentel/iotgateway/watcher"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
)

var (
	configFileName = flag.String("c", "/etc/sentel/gateway.conf", "config file")
)

func main() {
	flag.Parse()
	glog.Info("Starting gateway server...")

	config, _ := createConfig(*configFileName)
	mgr, err := service.NewServiceManager("gateway", config)
	if err != nil {
		glog.Fatalf("gateway create failed: '%s'", err.Error())
	}
	mgr.AddService(watcher.ServiceFactory{})
	glog.Fatal(mgr.RunAndWait())
}

func createConfig(fileName string) (config.Config, error) {
	config := config.New()
	config.AddConfig(defaultConfigs)
	config.AddConfigFile(fileName)
	options := map[string]map[string]string{}
	options["gateway"] = map[string]string{}
	// options["gateway"]["zookeeper"] = os.Getenv("ZOOKEEPER_HOST")
	config.AddConfig(options)
	return config, nil
}
