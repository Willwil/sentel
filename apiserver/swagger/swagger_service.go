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

	"github.com/go-openapi/runtime/middleware"
	"github.com/go-openapi/swag"
	"github.com/gorilla/handlers"
	"github.com/toqueteos/webbrowser"
	"github.com/tylerb/graceful"
)

func init() {
	loads.AddLoader(fmts.YAMLMatcher, fmts.YAMLDoc)
}

type swaggerService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	host      string
	port      int
	docServer *graceful.Server
}

type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	hosts := c.MustStringWithSection("swagger", "listen")
	names := strings.Split(hosts, ":")
	if len(names) != 2 {
		return nil, errors.New("swagger configuration error")
	}
	port, _ := strconv.Atoi(names[1])
	return &swaggerService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		host:      names[0],
		port:      port,
	}, nil
}

func (p *swaggerService) Name() string      { return "swagger" }
func (p *swaggerService) Initialize() error { return nil }

// Start
func (p *swaggerService) Start() error {
	specFile, err := p.config.StringWithSection("swagger", "path")
	if err != nil || specFile == "" {
		return errors.New("invalid swagger file")
	}
	specDoc, err := loads.Spec(specFile)
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
	p.waitgroup.Add(1)
	go func(p *swaggerService) {
		p.docServer = &graceful.Server{
			Server:           new(http.Server),
			NoSignalHandling: true,
		}
		p.docServer.SetKeepAlivesEnabled(true)
		p.docServer.TCPKeepAlive = 3 * time.Minute
		p.docServer.Handler = handler

		p.docServer.Serve(listener)
		p.waitgroup.Done()
	}(p)

	if ui, err := p.config.BoolWithSection("swagger", "open_browser"); err == nil && ui == true {
		webbrowser.Open(visit)
	}
	log.Println("serving docs at", visit)
	return nil

}

// Stop
func (p *swaggerService) Stop() {
	if p.docServer != nil {
		p.docServer.Stop(time.Duration(1 * time.Second))
	}
	p.waitgroup.Wait()
}
