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

package quto

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/service"
	"github.com/golang/glog"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	cachePolicyMemory = "memory"
	cachePolicyRedis  = "redis"
)

type qutoService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	quitChan    chan interface{}
	eventChan   chan event.Event
	cachePolicy string
	cacheMutex  sync.Mutex
	cacheMap    map[string]uint64
	rclient     *redis.Client
}

const (
	ServiceName = "quto"
)

type ServiceFactory struct{}

// New create metadata service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// check mongo db configuration
	hosts := c.MustString("mongo")
	timeout := c.MustInt("connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Connect with redis if cache policy is redis
	policy := cachePolicyMemory
	if v, err := c.StringWithSection("quto", "cace_policy"); err == nil && v == cachePolicyRedis {
		policy = cachePolicyRedis
	}

	var rclient *redis.Client = nil
	if policy == cachePolicyRedis {
		addr, _ := c.String("redis")
		password := c.MustStringWithSection("quto", "redis_password")
		db := c.MustIntWithSection("quto", "redis_db")

		rclient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		if _, err := rclient.Ping().Result(); err != nil {
			return nil, err
		}
	}

	return &qutoService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		quitChan:    make(chan interface{}),
		eventChan:   make(chan event.Event),
		cachePolicy: policy,
		cacheMutex:  sync.Mutex{},
		cacheMap:    make(map[string]uint64),
		rclient:     rclient,
	}, nil

}

// Name
func (p *qutoService) Name() string {
	return ServiceName
}

func (p *qutoService) Initialize() error { return nil }

// Start
func (p *qutoService) Start() error {
	event.Subscribe(event.QutoChange, onEventCallback, p)
	go func(p *qutoService) {
		for {
			select {
			case e := <-p.eventChan:
				p.handleQutoChanged(e)
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *qutoService) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.eventChan)
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e event.Event, ctx interface{}) {
	service := ctx.(*qutoService)
	service.eventChan <- e
}

// handleQutoChanged load changed quto into cach
func (p *qutoService) handleQutoChanged(e event.Event) {
	t := e.(*event.QutoChangeEvent)
	// check mongo db configuration
	hosts := p.config.MustString("mongo")
	timeout := p.config.MustInt("connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		glog.Errorf("quto: Access backend database failed for quto event:%d", t.Type)
		return
	}
	defer session.Close()
	c := session.DB("registry").C("qutos")

	// Load quto object
	quto := Quto{}
	if err := c.Find(bson.M{"qutoId": t.QutoId}).One(&quto); err != nil {
		glog.Errorf("quto: Failed to get qutotation '%s'", t.QutoId)
		return
	}
	// Update cach according to cach policy
	p.setCacheItem(quto.Key, quto.Value)
}

// getCacheItem get item from ache
func (p *qutoService) getCacheItem(id string) (uint64, error) {
	switch p.cachePolicy {
	case cachePolicyMemory:
		if _, ok := p.cacheMap[id]; !ok {
			return 0, fmt.Errorf("Invalid cache identifier '%s'", id)
		}
		return p.cacheMap[id], nil
	case cachePolicyRedis:
		val, err := p.rclient.Get(id).Result()
		if err != nil {
			return 0, err
		}
		return strconv.ParseUint(val, 10, 64)
	}
	return 0, errors.New("Invalid cache policy")
}

// setCacheItem set cache item internal
func (p *qutoService) setCacheItem(id string, val uint64) {
	switch p.cachePolicy {
	case cachePolicyMemory:
		p.cacheMap[id] = val
	case cachePolicyRedis:
		p.rclient.Set(id, val, 0)
	}
}

// getQuto return object's qutotation
func (p *qutoService) getQuto(id string) (uint64, error) {
	// Get quto from cache at first
	if val, err := p.getCacheItem(id); err == nil {
		return val, nil
	}

	// Read from database if not found in cache
	hosts := p.config.MustString("mongo")
	timeout := p.config.MustInt("connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		glog.Error("quto: Access backend database failed")
		return 0, err
	}
	defer session.Close()
	c := session.DB("registry").C("qutos")

	quto := Quto{}
	if err := c.Find(bson.M{"Key": id}).One(&quto); err != nil {
		return 0, err
	}
	// Set the item back to cache
	p.setCacheItem(id, quto.Value)
	return quto.Value, nil
}

// checkQutoWithAddValue check wether the newly added value is over quto
func (p *qutoService) checkQuto(id string, value uint64) bool {
	/*
		v, err := p.getCacheItem(id)
		if err != nil {
			return false
		}
		return v > value
	*/
	return true
}
