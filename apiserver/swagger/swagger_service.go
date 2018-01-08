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
	"strconv"
	"strings"
	"sync"

	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"path"
	"time"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/go-openapi/loads"
	"github.com/go-openapi/loads/fmts"
	"github.com/go-swagger/go-swagger/cmd/swagger/commands"

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gorilla/handlers"
	"github.com/toqueteos/webbrowser"
	"github.com/tylerb/graceful"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

var (
	// Debug is true when the SWAGGER_DEBUG env var is not empty
	Debug = os.Getenv("SWAGGER_DEBUG") != ""
)

type swaggerService struct {
	service.ServiceBase
	swaggerCmd  commands.ServeCmd
	host        string
	port        int
	swaggerFile string
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config, quit chan os.Signal) (service.Service, error) {
	hosts := c.MustString("apiserver", "swagger")
	names := strings.Split(hosts, ":")
	if len(names) != 2 {
		return nil, errors.New("swagger configuration error")
	}
	port, _ := strconv.Atoi(names[1])
	swaggerFile := c.MustString("swagger", "path")
	// swagger file
	return &swaggerService{
		ServiceBase: service.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		swaggerFile: swaggerFile,
		host:        names[0],
		port:        port,
	}, nil
}

func (p *swaggerService) Name() string { return "swagger" }

// Start
func (p *swaggerService) Start() error {
	go func(s *swaggerService) {
		p.execute()
	}(p)
	return nil
}

// Stop
func (p *swaggerService) Stop() {}

func (p *swaggerService) execute() error {
	specDoc, err := loads.Spec(p.swaggerFile)
	if err != nil {
		return err
	}
	b, err := json.MarshalIndent(specDoc.Spec(), "", "  ")
	if err != nil {
		return err
	}

	basePath := "/"
	listener, err := net.Listen("tcp4", net.JoinHostPort(p.host, strconv.Itoa(p.port)))
	if err != nil {
		return err
	}
	sh, sp, err := swag.SplitHostPort(listener.Addr().String())
	if err != nil {
		return err
	}
	if sh == "0.0.0.0" {
		sh = "localhost"
	}

	handler := http.NotFoundHandler()
	handler = middleware.Redoc(middleware.RedocOpts{
		BasePath: basePath,
		SpecURL:  path.Join(basePath, "swagger.json"),
		Path:     "docs",
	}, handler)
	visit := fmt.Sprintf("http://%s:%d%s", sh, sp, path.Join(basePath, "docs"))

	handler = handlers.CORS()(middleware.Spec(basePath, b, handler))
	go func() {
		docServer := &graceful.Server{Server: new(http.Server)}
		docServer.SetKeepAlivesEnabled(true)
		docServer.TCPKeepAlive = 3 * time.Minute
		docServer.Handler = handler

		docServer.Serve(listener)
	}()

	err = webbrowser.Open(visit)
	log.Println("serving docs at", visit)
	return nil
}
