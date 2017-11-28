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

package rpc

import (
	"fmt"
	"net"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/cloustone/sentel/broker/base"
	"github.com/cloustone/sentel/broker/metadata"
	"github.com/cloustone/sentel/broker/metric"
	sub "github.com/cloustone/sentel/broker/subtree"
	"github.com/cloustone/sentel/core"

	"github.com/golang/glog"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const ServiceName = "rpc"

type ApiService struct {
	base.ServiceBase
	listener net.Listener
	srv      *grpc.Server
}

// ApiServiceFactory
type ApiServiceFactory struct{}

// New create apiService service factory
func (m *ApiServiceFactory) New(c core.Config, quit chan os.Signal) (base.Service, error) {
	server := &ApiService{
		ServiceBase: base.ServiceBase{
			Config:    c,
			Quit:      quit,
			WaitGroup: sync.WaitGroup{},
		},
	}

	listen := c.MustString("rpc", "listen")
	lis, err := net.Listen("tcp", listen)
	if err != nil {
		glog.Fatal("Failed to listen: %v", err)
		return nil, err
	}
	server.listener = lis
	server.srv = grpc.NewServer()
	RegisterApiServer(server.srv, server)
	reflection.Register(server.srv)
	return server, nil

}

// Name
func (p *ApiService) Name() string {
	return ServiceName
}

func (p *ApiService) Initialize() error { return nil }

// Start
func (p *ApiService) Start() error {
	go func(p *ApiService) {
		p.srv.Serve(p.listener)
		p.WaitGroup.Add(1)
	}(p)
	return nil
}

// Stop
func (p *ApiService) Stop() {
	p.listener.Close()
	signal.Notify(p.Quit, syscall.SIGINT, syscall.SIGQUIT)
	p.WaitGroup.Wait()
	close(p.Quit)
}

//
// Wait
func (p *ApiService) Wait() {
	p.WaitGroup.Wait()
}

func (p *ApiService) Version(ctx context.Context, req *VersionRequest) (*VersionReply, error) {
	version := "0.1"
	return &VersionReply{Version: version}, nil
}

func (p *ApiService) Admins(ctx context.Context, req *AdminsRequest) (*AdminsReply, error) {
	return nil, nil
}

func (p *ApiService) Cluster(ctx context.Context, req *ClusterRequest) (*ClusterReply, error) {
	return nil, nil
}

func (p *ApiService) Status(ctx context.Context, req *StatusRequest) (*StatusReply, error) {
	return nil, nil
}

// Broker delegate broker command implementation in sentel
func (p *ApiService) Broker(ctx context.Context, req *BrokerRequest) (*BrokerReply, error) {
	switch req.Category {
	case "stats":
		stats := metric.GetStats(req.Service)
		return &BrokerReply{Stats: stats}, nil
	case "metrics":
		metrics := metric.GetMetric(req.Service)
		return &BrokerReply{Metrics: metrics}, nil
	default:
	}
	return nil, fmt.Errorf("Invalid broker request with categoru:%s", req.Category)
}

func (p *ApiService) Plugins(ctx context.Context, req *PluginsRequest) (*PluginsReply, error) {
	return nil, nil
}

// Services delegate  services command
func (p *ApiService) Services(ctx context.Context, req *ServicesRequest) (*ServicesReply, error) {
	/*
		reply := &ServicesReply{
			Header:   &ReplyMessageHeader{Success: true},
			Services: []*ServiceInfo{},
		}
		switch req.Category {
		case "list":
				services := metadata.GetAllServiceInfo()
				for _, service := range services {
					reply.Services = append(reply.Services,
						&ServiceInfo{
							ServiceName:    service.ServiceName,
							Listen:         service.Listen,
							Acceptors:      service.Acceptors,
							MaxClients:     service.MaxClients,
							CurrentClients: service.CurrentClients,
							ShutdownCount:  service.ShutdownCount,
						})
				}

		case "start":
		case "stop":
		}
	*/
	return nil, nil
}

//Subscriptions delete subscriptions command
func (p *ApiService) Subscriptions(ctx context.Context, req *SubscriptionsRequest) (*SubscriptionsReply, error) {
	reply := &SubscriptionsReply{
		Header:        &ReplyMessageHeader{Success: true},
		Subscriptions: []*SubscriptionInfo{},
	}
	switch req.Category {
	case "list":
		subs := sub.GetSubscriptions(req.Service)
		for _, sub := range subs {
			reply.Subscriptions = append(reply.Subscriptions,
				&SubscriptionInfo{
					ClientId:  sub.ClientId,
					Topic:     sub.Topic,
					Attribute: sub.Attribute,
				})
		}
	case "show":
		sub := sub.GetSubscription(req.Service, req.Subscription)
		if sub != nil {
			reply.Subscriptions = append(reply.Subscriptions,
				&SubscriptionInfo{
					ClientId:  sub.ClientId,
					Topic:     sub.Topic,
					Attribute: sub.Attribute,
				})
		}
	}
	return reply, nil
}

// Clients delegate clients command implementation in sentel
func (p *ApiService) Clients(ctx context.Context, req *ClientsRequest) (*ClientsReply, error) {
	reply := &ClientsReply{
		Clients: []*ClientInfo{},
		Header:  &ReplyMessageHeader{Success: true},
	}
	switch req.Category {
	case "list":
		// Get all client information for specified service
		clients := metadata.GetClients(req.Service)
		for _, client := range clients {
			reply.Clients = append(reply.Clients,
				&ClientInfo{
					UserName:     client.UserName,
					CleanSession: client.CleanSession,
					PeerName:     client.PeerName,
					ConnectTime:  client.ConnectTime,
				})
		}
	case "show":
		// Get client information for specified client id
		if client := metadata.GetClient(req.Service, req.ClientId); client != nil {
			reply.Clients = append(reply.Clients,
				&ClientInfo{
					UserName:     client.UserName,
					CleanSession: client.CleanSession,
					PeerName:     client.PeerName,
					ConnectTime:  client.ConnectTime,
				})
		}
	case "kick":
		/*
			broker := base.GetBroker()
			if err := broker.KickoffClient(req.Service, req.ClientId); err != nil {
				reply.Header.Success = false
				reply.Header.Reason = fmt.Sprintf("%v", err)
			}
		*/
	default:
		return nil, fmt.Errorf("Invalid category:'%s' for Clients api", req.Category)
	}
	return reply, nil
}

// Sessions delegate client sessions command
func (p *ApiService) Sessions(ctx context.Context, req *SessionsRequest) (*SessionsReply, error) {
	reply := &SessionsReply{
		Header:   &ReplyMessageHeader{Success: true},
		Sessions: []*SessionInfo{},
	}
	switch req.Category {
	case "list":
		/*
			sessions := metadata.GetSessions(req.Service, req.Conditions)
			for _, session := range sessions {
				reply.Sessions = append(reply.Sessions,
					&metadata.SessionInfo{
						ClientId:           session.ClientId,
						CreatedAt:          session.CreatedAt,
						CleanSession:       session.CleanSession,
						MessageMaxInflight: session.MessageMaxInflight,
						MessageInflight:    session.MessageInflight,
						MessageInQueue:     session.MessageInQueue,
						MessageDropped:     session.MessageDropped,
						AwaitingRel:        session.AwaitingRel,
						AwaitingComp:       session.AwaitingComp,
						AwaitingAck:        session.AwaitingAck,
					})
			}
		*/
	case "show":
		/*
			session := metadata.GetSession(req.Service, req.ClientId)
			if session != nil {
				reply.Sessions = append(reply.Sessions,
					&SessionInfo{
						ClientId:           session.ClientId,
						CreatedAt:          session.CreatedAt,
						CleanSession:       session.CleanSession,
						MessageMaxInflight: session.MessageMaxInflight,
						MessageInflight:    session.MessageInflight,
						MessageInQueue:     session.MessageInQueue,
						MessageDropped:     session.MessageDropped,
						AwaitingRel:        session.AwaitingRel,
						AwaitingComp:       session.AwaitingComp,
						AwaitingAck:        session.AwaitingAck,
					})
			}
		*/
	}
	return reply, nil
}

func (p *ApiService) Topics(ctx context.Context, req *TopicsRequest) (*TopicsReply, error) {
	reply := &TopicsReply{
		Header: &ReplyMessageHeader{Success: true},
		Topics: []*TopicInfo{},
	}
	switch req.Category {
	case "list":
		topics := sub.GetTopics(req.Service)
		for _, topic := range topics {
			reply.Topics = append(reply.Topics,
				&TopicInfo{
					Topic:     topic.Topic,
					Attribute: topic.Attribute,
				})
		}
	case "show":
		topic := sub.GetTopic(req.Service, req.Topic)
		if topic != nil {
			reply.Topics = append(reply.Topics,
				&TopicInfo{
					Topic:     topic.Topic,
					Attribute: topic.Attribute,
				})
		}
	}
	return reply, nil
}

func (p *ApiService) Routes(ctx context.Context, req *RoutesRequest) (*RoutesReply, error) {
	return nil, nil
}
