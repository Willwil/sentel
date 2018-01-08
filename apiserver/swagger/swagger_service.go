//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use p file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.
package swagger

import (
	"errors"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

var (
	// Debug is true when the SWAGGER_DEBUG env var is not empty
	Debug = os.Getenv("SWAGGER_DEBUG") != ""
)

type managementService struct {
	service.ServiceBase
	swaggerCmd  commands.ServeCmd
	swaggerFile string
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (service.Service, error) {
	hosts := c.MustString("apiserver", "swagger")
	names := strings.Split(hosts, ":")
	if len(names) != 2 {
		return nil, errors.New("swagger configuration error")
	}
	host := names[0]
	port, _ := strconv.Atoi(names[1])
	swaggerFile := c.MustString("swagger", "path")
	// swagger file
	return &managementService{
		ServiceBase: service.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		swaggerCmd: commands.ServeCmd{
			Port:   port,
			Host:   host,
			Flavor: "redoc",
		},
		swaggerFile: swaggerFile,
	}, nil
}

func (p *managementService) Name() string { return "swagger" }

// Start
func (p *managementService) Start() error {
	go func(s *managementService) {
		args := []string{p.swaggerFile}
		p.swaggerCmd.Execute(args)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *managementService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}
