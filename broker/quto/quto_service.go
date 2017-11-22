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
	"encoding/json"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	cachePolicyLru   = "lru"
	cachePolicyRedis = "redis"
)

// Metaservice manage broker metadata
// Broker's metadata include the following data
// - Global broker cluster data
// - Shadow device
type QutoService struct {
	core.ServiceBase
	eventChan   chan *event.Event
	cachePolicy string
	cacheMutex  sync.Mutex
	lruCache    *LRUCache
	rclient     *redis.Client
}

const (
	ServiceName = "metadata"
)

// QutoServiceFactory
type QutoServiceFactory struct{}

// New create metadata service factory
func (p *QutoServiceFactory) New(c core.Config, quit chan os.Signal) (core.Service, error) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(c, "broker", "mongo")
	timeout := c.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Connect with redis if cache policy is redis
	policy := cachePolicyLru
	if v, err := c.String("quto", "cace_policy"); err == nil && v == cachePolicyRedis {
		policy = cachePolicyRedis
	}

	var rclient *redis.Client = nil
	if policy == cachePolicyRedis {
		addr, _ := core.GetServiceEndpoint(c, "quto", "redis")
		password := c.MustString("quto", "redis_password")
		db := c.MustInt("quto", "redis_db")

		rclient = redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: password,
			DB:       db,
		})

		if _, err := rclient.Ping().Result(); err != nil {
			return nil, err
		}
	}

	return &QutoService{
		ServiceBase: core.ServiceBase{
			Config:    c,
			WaitGroup: sync.WaitGroup{},
			Quit:      quit,
		},
		eventChan:   make(chan *event.Event),
		cachePolicy: policy,
		cacheMutex:  sync.Mutex{},
		lruCache:    NewLRUCache(5),
		rclient:     rclient,
	}, nil

}

// Name
func (p *QutoService) Name() string {
	return ServiceName
}

// Start
func (p *QutoService) Start() error {
	event.Subscribe(event.QutoChanged, onEventCallback, p)
	go func(p *QutoService) {
		for {
			select {
			case e := <-p.eventChan:
				p.handleQutoChanged(e)
			case <-p.Quit:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *QutoService) Stop() {
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
	close(p.eventChan)
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*QutoService)
	service.eventChan <- e
}

// handleQutoChanged load changed quto into cach
func (p *QutoService) handleQutoChanged(e *event.Event) {
	// check mongo db configuration
	hosts, _ := core.GetServiceEndpoint(p.Config, "broker", "mongo")
	timeout := p.Config.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		glog.Errorf("quto: Access backend database failed for quto event:%d", e.Type)
		return
	}
	defer session.Close()
	c := session.DB("registry").C("qutos")

	// Load quto object
	quto := Quto{}
	if err := c.Find(bson.M{"qutoId": e.QutoId}).One(&quto); err != nil {
		glog.Errorf("quto: Failed to get qutotation '%s'", e.QutoId)
		return
	}
	// Update cach according to cach policy
	p.setQuto(quto.Target, quto.TargetId, &quto)
}

// getCacheItem get item from ache
func (p *QutoService) getCacheItem(key string) *Quto {
	switch p.cachePolicy {
	case cachePolicyLru:
		quto, found, _ := p.lruCache.Get(key)
		if found && quto != nil {
			return quto.(*Quto)
		}
	case cachePolicyRedis:
		val, err := p.rclient.Get(key).Result()
		if err == nil {
			quto := Quto{}
			if err = json.Unmarshal([]byte(val), &quto); err == nil {
				return &quto
			}
		}
	}
	return nil
}

// setCacheItem set cache item internal
func (p *QutoService) setCacheItem(key string, quto *Quto) error {
	switch p.cachePolicy {
	case cachePolicyLru:
		p.lruCache.Set(key, quto)
	case cachePolicyRedis:
		val, err := json.Marshal(quto)
		if err != nil {
			return err
		}
		p.rclient.Set(key, []byte(val), 0)
	}
	return nil
}

// getQuto return object's qutotation
func (p *QutoService) getQuto(target string, id string) (*Quto, error) {
	key := target + "/" + id
	// Get quto from cache at first
	if quto := p.getCacheItem(key); quto != nil {
		return quto, nil
	}

	// Read from database if not found in cache
	hosts, _ := core.GetServiceEndpoint(p.Config, "broker", "mongo")
	timeout := p.Config.MustInt("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		glog.Error("quto: Access backend database failed")
		return nil, err
	}
	defer session.Close()
	c := session.DB("registry").C("qutos")

	quto := Quto{}
	if err := c.Find(bson.M{"Target": target, "TargetId": id}).One(&quto); err != nil {
		return nil, err
	}
	// Set the item back to cache
	if err := p.setCacheItem(key, &quto); err != nil {
		glog.Error("quto: Failed to write quto back to cache")
		// no return
	}
	return &quto, nil
}

// getQutoItem return object's qutotation
func (p *QutoService) getQutoItem(target, id, item string) uint64 {
	return 0
}

// setQuto set object's qutotation
func (p *QutoService) setQuto(target, id string, q *Quto) {
}

// setQutoItem set object item's qutotations
func (p *QutoService) setQutoItem(taret, id, item string, value uint64) {
}

// addQutoItemValue add quto items's value
func (p *QutoService) addQutoItem(target, id, item string, value uint64) {
}

// subQutoItemValue sub quto items's value
func (p *QutoService) subQutoItem(target, id, item string, value uint64) {
}
