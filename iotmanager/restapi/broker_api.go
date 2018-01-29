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

package restapi

import (
	pb "github.com/cloustone/sentel/broker/rpc"

	"github.com/golang/glog"
	context "golang.org/x/net/context"
	"google.golang.org/grpc"
)

type brokerApi struct {
	rpcapi pb.ApiClient
	conn   *grpc.ClientConn
}

func newBrokerApi(hosts string) (*brokerApi, error) {
	conn, err := grpc.Dial(hosts, grpc.WithInsecure())
	if err != nil {
		glog.Errorf("Failed to connect with iothub:%s", err)
		return nil, err
	}
	c := pb.NewApiClient(conn)
	return &brokerApi{conn: conn, rpcapi: c}, nil
}

func (s *brokerApi) version(in *pb.VersionRequest) (*pb.VersionReply, error) {
	return s.rpcapi.Version(context.Background(), in)
}
func (s *brokerApi) status(in *pb.StatusRequest) (*pb.StatusReply, error) {
	return s.rpcapi.Status(context.Background(), in)
}

func (s *brokerApi) services(in *pb.ServicesRequest) (*pb.ServicesReply, error) {
	return s.rpcapi.Services(context.Background(), in)
}

func (s *brokerApi) subscriptions(in *pb.SubscriptionsRequest) (*pb.SubscriptionsReply, error) {
	return s.rpcapi.Subscriptions(context.Background(), in)
}

func (s *brokerApi) clients(in *pb.ClientsRequest) (*pb.ClientsReply, error) {
	return s.rpcapi.Clients(context.Background(), in)
}

func (s *brokerApi) sessions(in *pb.SessionsRequest) (*pb.SessionsReply, error) {
	return s.rpcapi.Sessions(context.Background(), in)
}

func (s *brokerApi) topics(in *pb.TopicsRequest) (*pb.TopicsReply, error) {
	return s.rpcapi.Topics(context.Background(), in)
}
