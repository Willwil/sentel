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
package apiservice

import (
	"fmt"
	"sync"

	"github.com/cloustone/sentel/apiserver/base"
	"github.com/cloustone/sentel/apiserver/util"
	"github.com/cloustone/sentel/pkg/config"
	"github.com/cloustone/sentel/pkg/goshiro"
	"github.com/cloustone/sentel/pkg/goshiro/shiro"
	"github.com/cloustone/sentel/pkg/goshiro/web"
	"github.com/cloustone/sentel/pkg/registry"
	"github.com/cloustone/sentel/pkg/service"
	jwt "github.com/dgrijalva/jwt-go"

	echo "github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
)

var resourceMaps = make(map[string]string)

type consoleService struct {
	config      config.Config
	waitgroup   sync.WaitGroup
	version     string
	echo        *echo.Echo
	securityMgr shiro.SecurityManager
}
type ServiceFactory struct{}

func (p ServiceFactory) New(c config.Config) (service.Service, error) {
	// create resource maps
	for _, res := range apiPolicies {
		resourceMaps[res.Path] = res.Resource
	}
	realm, err := base.NewAuthorizeRealm(c)
	if err != nil {
		return nil, err
	}
	securityMgr, err := goshiro.NewSecurityManager(c, realm)
	if err != nil {
		return nil, err
	}
	securityMgr.Load()

	return &consoleService{
		config:      c,
		waitgroup:   sync.WaitGroup{},
		echo:        echo.New(),
		securityMgr: securityMgr,
	}, nil
}

func (p *consoleService) Name() string { return "apiservice" }
func (p *consoleService) Initialize() error {
	c := p.config
	if err := registry.Initialize(c); err != nil {
		return fmt.Errorf("registry initialize failed:%v", err)
	}
	p.echo.HideBanner = true
	p.echo.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			ctx.Set("SecurityManager", p.securityMgr)
			return h(ctx)
		}
	})

	// Initialize middleware
	//Cross-Origin
	p.echo.Use(mw.CORSWithConfig(mw.DefaultCORSConfig))

	p.echo.Use(mw.RequestID())
	p.echo.Use(mw.LoggerWithConfig(mw.LoggerConfig{
		Format: "${time_unix},method=${method}, uri=${uri}, status=${status}\n",
	}))
	// Queue APIs
	p.echo.POST("mns/v1/api/queues/:queueName", createQueue)
	p.echo.PUT("mns/v1/api/queues/:queueName", setQueueAttribute)
	p.echo.GET("mns/v1/api/queues/:queueName", getQueueAttribute)
	p.echo.GET("mns/v1/api/queues", getQueueList)
	p.echo.DELETE("mns/v1/api/queues/:queueName", deleteQueue)
	p.echo.POST("mns/v1/api/queues/:queueName/messages", sendQueueMessage)
	p.echo.POST("mns/v1/api/queues/:queueName/messages?batch=true", batchSendQueueMessage)
	p.echo.GET("mns/v1/api/queues/:queueName/messages?ws=:ws", receiveQueueMessage)
	p.echo.GET("mns/v1/api/queues/:queueName/messages?ws=:ws&nbr=:numberOfMesssages", batchReceiveQueueMessage)
	p.echo.DELETE("mns/v1/api/queues/:queueName/messages?msgId=:msgId", deleteQueueMessage)
	p.echo.DELETE("mns/v1/api/queues/:queueName/messages?batch=true", batchDeleteQueueMessages)
	p.echo.GET("mns/v1/api/queues/:queueName/messages?ws=:ws&peekonly=true", peekQueueMessages)
	p.echo.GET("mns/v1/api/queues/:queueName/messages?peekonly=true&batch=true", batchPeekQueueMessages)

	// Topics API
	p.echo.POST("mns/v1/api/topics/:topicName", createTopic)
	p.echo.PUT("mns/v1/api/topics/:topicName", setTopicAttribute)
	p.echo.GET("mns/v1/api/topics/:topicName", getTopicAttribute)
	p.echo.DELETE("mns/v1/api/topics/:topicName", deleteTopic)
	p.echo.GET("mns/v1/api/topics", listTopics)

	// Subscriptions API
	p.echo.POST("mns/v1/api/topics/:topicName/subscriptions/:subscriptionName", subscribe)
	p.echo.PUT("mns/v1/api/topics/:topicName/subscriptions/:subscriptionName", setSubscriptionAttr)
	p.echo.GET("mns/v1/api/topics/:topicName/subscriptions/:subscriptionName", getSubscriptionAttr)
	p.echo.DELETE("mns/v1/api/topics/:topicName/subscriptions/:subscriptionName", unsubscribe)
	p.echo.GET("mns/v1/api/topics/:topicName/subscriptions?pagesize=:pageSize&pageno=:pageNo&index=:startIndex", listTopicSubscriptions)

	// Messages API
	p.echo.POST("mns/v1/api/topics/:topicName/messages", publishMessage)
	p.echo.POST("mns/v1/api/notifications", publishNotification)

	return nil

}

// Start
func (p *consoleService) Start() error {
	p.waitgroup.Add(1)
	go func(s *consoleService) {
		addr := p.config.MustStringWithSection("console", "listen")
		p.echo.Start(addr)
		p.waitgroup.Done()
	}(p)
	return nil
}

// Stop
func (p *consoleService) Stop() {
	p.echo.Close()
	p.waitgroup.Wait()
}

// setAuth setup api group 's authentication method
func (p *consoleService) setAuth(c config.Config, g *echo.Group) {
	auth := util.StringConfigWithDefaultValue(c, "auth", "jwt")
	switch auth {
	case "jwt":
		// Authentication config
		config := mw.JWTConfig{
			Claims:     &base.ApiJWTClaims{},
			SigningKey: []byte("secret"),
		}
		g.Use(mw.JWTWithConfig(config))
		g.Use(accessIdWithConfig(c))
		if util.AuthNeed(c) {
			g.Use(authorizeWithConfig(c))
		}
	default:
	}
}

func accessIdWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			// After authenticated by gateway,the authentication paramters must bevalid
			if user, ok := ctx.Get("user").(*jwt.Token); ok {
				claims := user.Claims.(*base.ApiJWTClaims)
				ctx.Set("AccessId", claims.AccessId)
			}
			return next(ctx)
		}
	}
}

func authorizeWithConfig(config config.Config) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(ctx echo.Context) error {
			securityManager := base.GetSecurityManager(ctx)
			accessId := ctx.Get("AccessId").(string)
			token := web.JWTToken{Username: accessId, Authenticated: true}
			resource, action := base.GetRequestInfo(ctx, resourceMaps)
			if err := securityManager.Authorize(token, resource, action); err != nil {
				return err
			}
			return next(ctx)
		}
	}
}
