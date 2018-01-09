//  Licensed under the Apache License, Version 2.0 (the "License"); you may
//  not use this file except in compliance with the License. You may obtain
//  a copy of the License at
//
//        http://www.apache.org/licenses/LICENSE-2.0
//
//  Unless required by applicable law or agreed to in writing, software
//  distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  License for the specific language governing permissions and limitations
//  under the License.

package auth

import (
	"fmt"
	"sync"
	"time"

	"github.com/cloustone/sentel/broker/event"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"

	"github.com/go-redis/redis"
	"gopkg.in/mgo.v2"
)

const (
	ServiceName = "auth"
)

type ServiceFactory struct{}

// New create coap service factory
func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	s := &authService{
		config:    c,
		waitgroup: sync.WaitGroup{},
		eventChan: make(chan *event.Event),
		quitChan:  make(chan interface{}),
	}
	// check mongo db configuration
	hosts, _ := c.String("broker", "mongo")
	timeout, _ := c.Int("broker", "connect_timeout")
	session, err := mgo.DialWithTimeout(hosts, time.Duration(timeout)*time.Second)
	if err != nil {
		return nil, err
	}
	defer session.Close()

	// Connect with redis if cache policy is redis
	addr, _ := c.String("broker", "redis")
	password, _ := c.String("broker", "redis_password")
	db, _ := c.Int("auth", "redis_db")
	rc := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})
	if _, err := rc.Ping().Result(); err == nil {
		s.redis = rc
	}
	return s, nil
}

// Authentication Service
type authService struct {
	config    config.Config
	waitgroup sync.WaitGroup
	redis     *redis.Client
	eventChan chan *event.Event
	quitChan  chan interface{}
}

// Name
func (p *authService) Name() string      { return ServiceName }
func (p *authService) Initialize() error { return nil }

// Start
func (p *authService) Start() error {
	event.Subscribe(event.AuthChange, onEventCallback, p)
	p.waitgroup.Add(1)
	go func(p *authService) {
		defer p.waitgroup.Done()
		for {
			select {
			case e := <-p.eventChan:
				p.handleEvent(e)
			case <-p.quitChan:
				return
			}
		}
	}(p)
	return nil
}

// Stop
func (p *authService) Stop() {
	p.quitChan <- true
	p.waitgroup.Wait()
	close(p.quitChan)
	close(p.eventChan)
}

// CheckAcl check client's access control right
func (p *authService) authorize(ctx Context, clientid string, topic string, access int) error {
	return nil
}

// authenticate check user's name and password
func (p *authService) authenticate(ctx Context) error {
	if key, err := p.getDeviceSecretKey(ctx); err == nil {
		ctx.DeviceSecret = key
		return sign(ctx)
	}
	return fmt.Errorf("auth: Failed to get device secret key for '%s'", ctx.DeviceName)
}

func (p *authService) handleEvent(e *event.Event) {

}

// getDeviceSecretKey retrieve device secret key from cache or mongo
func (p *authService) getDeviceSecretKey(ctx Context) (string, error) {
	// Read from cache at first
	key := ctx.ProductId + "/" + ctx.DeviceName
	if p.redis != nil {
		if val, err := p.redis.Get(key).Result(); err == nil {
			return val, nil
		}
	}
	// Read from database if not found in cache
	r, err := registry.New("broker", p.config)
	if err != nil {
		return "", err
	}
	defer r.Close()
	device, err := r.GetDeviceByName(ctx.ProductId, ctx.DeviceName)
	if err != nil {
		return "", err
	}
	// Write back to redis
	if p.redis != nil {
		p.redis.Set(key, device.DeviceSecret, 0)
	}
	return device.DeviceSecret, nil
}

// onEventCallback will be called when notificaiton come from event service
func onEventCallback(e *event.Event, ctx interface{}) {
	service := ctx.(*authService)
	service.eventChan <- e
}
