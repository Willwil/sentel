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

package rpc

import (
	"github.com/cloustone/sentel/core"
	"github.com/golang/glog"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type BrokerApi struct {
	rpcapi ApiClient
	conn   *grpc.ClientConn
}

func NewBrokerApi(c core.Config) (*BrokerApi, error) {
	address, err := c.String("broker", "server")
	if err != nil {
		address = "localhost:50052"
	}

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		glog.Fatalf("Failed to connect with iothub:%s", err)
		return nil, err
	}
	rpcapi := NewApiClient(conn)
	return &BrokerApi{conn: conn, rpcapi: rpcapi}, nil
}

func (p *BrokerApi) Version(in *VersionRequest) (*VersionReply, error) {
	return p.rpcapi.Version(context.Background(), in)
}
func (p *BrokerApi) Admins(in *AdminsRequest) (*AdminsReply, error) {
	return p.rpcapi.Admins(context.Background(), in)
}

func (p *BrokerApi) Cluster(in *ClusterRequest) (*ClusterReply, error) {
	return p.rpcapi.Cluster(context.Background(), in)
}

func (p *BrokerApi) Routes(in *RoutesRequest) (*RoutesReply, error) {
	return p.rpcapi.Routes(context.Background(), in)
}

func (p *BrokerApi) Status(in *StatusRequest) (*StatusReply, error) {
	return p.rpcapi.Status(context.Background(), in)
}

func (p *BrokerApi) Broker(in *BrokerRequest) (*BrokerReply, error) {
	return p.rpcapi.Broker(context.Background(), in)
}

func (p *BrokerApi) Plugins(in *PluginsRequest) (*PluginsReply, error) {
	return p.rpcapi.Plugins(context.Background(), in)
}

func (p *BrokerApi) Services(in *ServicesRequest) (*ServicesReply, error) {
	return p.rpcapi.Services(context.Background(), in)
}

func (p *BrokerApi) Subscriptions(in *SubscriptionsRequest) (*SubscriptionsReply, error) {
	return p.rpcapi.Subscriptions(context.Background(), in)
}

func (p *BrokerApi) Clients(in *ClientsRequest) (*ClientsReply, error) {
	return p.rpcapi.Clients(context.Background(), in)
}

func (p *BrokerApi) Sessions(in *SessionsRequest) (*SessionsReply, error) {
	return p.rpcapi.Sessions(context.Background(), in)
}

func (p *BrokerApi) Topics(in *TopicsRequest) (*TopicsReply, error) {
	return p.rpcapi.Topics(context.Background(), in)
}
