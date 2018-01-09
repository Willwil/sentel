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

package l2

import (
	"errors"
	"sync"

	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
)

type l2osService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	quitChan  chan interface{}
	cmds      chan interface{}
	storage   objectStorage
}

// l2osServiceFactory
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	s, err := newStorage(c)
	if err != nil {
		return nil, errors.New("create storage failed")
	}
	return &l2osService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		quitChan:  make(chan interface{}),
		cmds:      make(chan interface{}),
		storage:   s,
	}, nil

}

func getService() *l2osService {
	mgr := service.GetServiceManager()
	return mgr.GetService("l2os").(*l2osService)
}

func getStorage() objectStorage {
	service := getService()
	return service.storage
}

// Name
func (p *l2osService) Name() string { return "l2os" }

// Initialize
func (p *l2osService) Initialize() error { return nil }

// Start
func (p *l2osService) Start() error {
	p.waitgroup.Add(1)
	go func(s *l2osService) {
		defer p.waitgroup.Done()
		for {
			select {
			case <-p.quitChan:
				return
			case cmd := <-p.cmds:
				p.handleRequest(cmd)
			}
		}
	}(p)
	return nil
}

// Stop
func (p *l2osService) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.cmds)
}

func (p *l2osService) handleRequest(cmd interface{}) {
	switch cmd.(type) {
	case createObjectReq:
		p.createObject(cmd)
	case destroyObjectReq:
		p.destroyObject(cmd)
	case getObjectReq:
		p.getObject(cmd)
	case setObjectAttrReq:
		p.setObjectAttr(cmd)
	case getObjectAttrReq:
		p.getObjectAttr(cmd)
	case updateObjectReq:
		p.updateObject(cmd)
	case createAccountReq:
		p.createAccount(cmd)
	case destroyAccountReq:
		p.destroyAccount(cmd)
	}
}

func (p *l2osService) createObject(cmd interface{}) {
	req := cmd.(createObjectReq)
	err := p.storage.createObject(req.object)
	req.resp <- err
}

func (p *l2osService) destroyObject(cmd interface{}) {
	req := cmd.(destroyObjectReq)
	err := p.storage.deleteObject(req.objid)
	req.resp <- err
}
func (p *l2osService) getObject(cmd interface{}) {
	req := cmd.(getObjectReq)
	obj, _ := p.storage.getObject(req.objid)
	req.resp <- obj
}

func (p *l2osService) setObjectAttr(cmd interface{}) {}
func (p *l2osService) getObjectAttr(cmd interface{}) {}

func (p *l2osService) updateObject(cmd interface{}) {
	req := cmd.(updateObjectReq)
	err := p.storage.updateObject(req.object)
	req.resp <- err
}

func (p *l2osService) createAccount(cmd interface{}) {
	req := cmd.(createAccountReq)
	err := p.storage.createAccount(req.name)
	req.resp <- err
}
func (p *l2osService) destroyAccount(cmd interface{}) {
	req := cmd.(destroyAccountReq)
	err := p.storage.destroyAccount(req.name)
	req.resp <- err
}

func (p *l2osService) sendRequest(req interface{}) {
	p.cmds <- req
}
