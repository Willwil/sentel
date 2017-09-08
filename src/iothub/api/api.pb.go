// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api.proto

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	api.proto

It has these top-level messages:
	VersionRequest
	VersionReply
	AdminsRequest
	AdminsReply
	ClusterRequest
	ClusterReply
	RoutesRequest
	RoutesReply
	StatusRequest
	StatusReply
	BrokerRequest
	BrokerReply
	PluginsRequest
	PluginsReply
	ServicesRequest
	ServicesReply
	SubscriptionsRequest
	SubscriptionsReply
	ClientsRequest
	ClientsReply
	SessionsRequest
	SessionsReply
	TopicsRequest
	TopicsReply
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type VersionRequest struct {
}

func (m *VersionRequest) Reset()                    { *m = VersionRequest{} }
func (m *VersionRequest) String() string            { return proto.CompactTextString(m) }
func (*VersionRequest) ProtoMessage()               {}
func (*VersionRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type VersionReply struct {
	Version string `protobuf:"bytes,1,opt,name=version" json:"version,omitempty"`
}

func (m *VersionReply) Reset()                    { *m = VersionReply{} }
func (m *VersionReply) String() string            { return proto.CompactTextString(m) }
func (*VersionReply) ProtoMessage()               {}
func (*VersionReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

func (m *VersionReply) GetVersion() string {
	if m != nil {
		return m.Version
	}
	return ""
}

type AdminsRequest struct {
}

func (m *AdminsRequest) Reset()                    { *m = AdminsRequest{} }
func (m *AdminsRequest) String() string            { return proto.CompactTextString(m) }
func (*AdminsRequest) ProtoMessage()               {}
func (*AdminsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

type AdminsReply struct {
}

func (m *AdminsReply) Reset()                    { *m = AdminsReply{} }
func (m *AdminsReply) String() string            { return proto.CompactTextString(m) }
func (*AdminsReply) ProtoMessage()               {}
func (*AdminsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ClusterRequest struct {
}

func (m *ClusterRequest) Reset()                    { *m = ClusterRequest{} }
func (m *ClusterRequest) String() string            { return proto.CompactTextString(m) }
func (*ClusterRequest) ProtoMessage()               {}
func (*ClusterRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type ClusterReply struct {
}

func (m *ClusterReply) Reset()                    { *m = ClusterReply{} }
func (m *ClusterReply) String() string            { return proto.CompactTextString(m) }
func (*ClusterReply) ProtoMessage()               {}
func (*ClusterReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

type RoutesRequest struct {
}

func (m *RoutesRequest) Reset()                    { *m = RoutesRequest{} }
func (m *RoutesRequest) String() string            { return proto.CompactTextString(m) }
func (*RoutesRequest) ProtoMessage()               {}
func (*RoutesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

type RoutesReply struct {
}

func (m *RoutesReply) Reset()                    { *m = RoutesReply{} }
func (m *RoutesReply) String() string            { return proto.CompactTextString(m) }
func (*RoutesReply) ProtoMessage()               {}
func (*RoutesReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

type StatusRequest struct {
}

func (m *StatusRequest) Reset()                    { *m = StatusRequest{} }
func (m *StatusRequest) String() string            { return proto.CompactTextString(m) }
func (*StatusRequest) ProtoMessage()               {}
func (*StatusRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type StatusReply struct {
}

func (m *StatusReply) Reset()                    { *m = StatusReply{} }
func (m *StatusReply) String() string            { return proto.CompactTextString(m) }
func (*StatusReply) ProtoMessage()               {}
func (*StatusReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

type BrokerRequest struct {
}

func (m *BrokerRequest) Reset()                    { *m = BrokerRequest{} }
func (m *BrokerRequest) String() string            { return proto.CompactTextString(m) }
func (*BrokerRequest) ProtoMessage()               {}
func (*BrokerRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

type BrokerReply struct {
}

func (m *BrokerReply) Reset()                    { *m = BrokerReply{} }
func (m *BrokerReply) String() string            { return proto.CompactTextString(m) }
func (*BrokerReply) ProtoMessage()               {}
func (*BrokerReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

type PluginsRequest struct {
}

func (m *PluginsRequest) Reset()                    { *m = PluginsRequest{} }
func (m *PluginsRequest) String() string            { return proto.CompactTextString(m) }
func (*PluginsRequest) ProtoMessage()               {}
func (*PluginsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{12} }

type PluginsReply struct {
}

func (m *PluginsReply) Reset()                    { *m = PluginsReply{} }
func (m *PluginsReply) String() string            { return proto.CompactTextString(m) }
func (*PluginsReply) ProtoMessage()               {}
func (*PluginsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{13} }

type ServicesRequest struct {
}

func (m *ServicesRequest) Reset()                    { *m = ServicesRequest{} }
func (m *ServicesRequest) String() string            { return proto.CompactTextString(m) }
func (*ServicesRequest) ProtoMessage()               {}
func (*ServicesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{14} }

type ServicesReply struct {
}

func (m *ServicesReply) Reset()                    { *m = ServicesReply{} }
func (m *ServicesReply) String() string            { return proto.CompactTextString(m) }
func (*ServicesReply) ProtoMessage()               {}
func (*ServicesReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{15} }

type SubscriptionsRequest struct {
}

func (m *SubscriptionsRequest) Reset()                    { *m = SubscriptionsRequest{} }
func (m *SubscriptionsRequest) String() string            { return proto.CompactTextString(m) }
func (*SubscriptionsRequest) ProtoMessage()               {}
func (*SubscriptionsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{16} }

type SubscriptionsReply struct {
}

func (m *SubscriptionsReply) Reset()                    { *m = SubscriptionsReply{} }
func (m *SubscriptionsReply) String() string            { return proto.CompactTextString(m) }
func (*SubscriptionsReply) ProtoMessage()               {}
func (*SubscriptionsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{17} }

type ClientsRequest struct {
}

func (m *ClientsRequest) Reset()                    { *m = ClientsRequest{} }
func (m *ClientsRequest) String() string            { return proto.CompactTextString(m) }
func (*ClientsRequest) ProtoMessage()               {}
func (*ClientsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{18} }

type ClientsReply struct {
}

func (m *ClientsReply) Reset()                    { *m = ClientsReply{} }
func (m *ClientsReply) String() string            { return proto.CompactTextString(m) }
func (*ClientsReply) ProtoMessage()               {}
func (*ClientsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{19} }

type SessionsRequest struct {
}

func (m *SessionsRequest) Reset()                    { *m = SessionsRequest{} }
func (m *SessionsRequest) String() string            { return proto.CompactTextString(m) }
func (*SessionsRequest) ProtoMessage()               {}
func (*SessionsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{20} }

type SessionsReply struct {
}

func (m *SessionsReply) Reset()                    { *m = SessionsReply{} }
func (m *SessionsReply) String() string            { return proto.CompactTextString(m) }
func (*SessionsReply) ProtoMessage()               {}
func (*SessionsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{21} }

type TopicsRequest struct {
}

func (m *TopicsRequest) Reset()                    { *m = TopicsRequest{} }
func (m *TopicsRequest) String() string            { return proto.CompactTextString(m) }
func (*TopicsRequest) ProtoMessage()               {}
func (*TopicsRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{22} }

type TopicsReply struct {
}

func (m *TopicsReply) Reset()                    { *m = TopicsReply{} }
func (m *TopicsReply) String() string            { return proto.CompactTextString(m) }
func (*TopicsReply) ProtoMessage()               {}
func (*TopicsReply) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{23} }

func init() {
	proto.RegisterType((*VersionRequest)(nil), "api.VersionRequest")
	proto.RegisterType((*VersionReply)(nil), "api.VersionReply")
	proto.RegisterType((*AdminsRequest)(nil), "api.AdminsRequest")
	proto.RegisterType((*AdminsReply)(nil), "api.AdminsReply")
	proto.RegisterType((*ClusterRequest)(nil), "api.ClusterRequest")
	proto.RegisterType((*ClusterReply)(nil), "api.ClusterReply")
	proto.RegisterType((*RoutesRequest)(nil), "api.RoutesRequest")
	proto.RegisterType((*RoutesReply)(nil), "api.RoutesReply")
	proto.RegisterType((*StatusRequest)(nil), "api.StatusRequest")
	proto.RegisterType((*StatusReply)(nil), "api.StatusReply")
	proto.RegisterType((*BrokerRequest)(nil), "api.BrokerRequest")
	proto.RegisterType((*BrokerReply)(nil), "api.BrokerReply")
	proto.RegisterType((*PluginsRequest)(nil), "api.PluginsRequest")
	proto.RegisterType((*PluginsReply)(nil), "api.PluginsReply")
	proto.RegisterType((*ServicesRequest)(nil), "api.ServicesRequest")
	proto.RegisterType((*ServicesReply)(nil), "api.ServicesReply")
	proto.RegisterType((*SubscriptionsRequest)(nil), "api.SubscriptionsRequest")
	proto.RegisterType((*SubscriptionsReply)(nil), "api.SubscriptionsReply")
	proto.RegisterType((*ClientsRequest)(nil), "api.ClientsRequest")
	proto.RegisterType((*ClientsReply)(nil), "api.ClientsReply")
	proto.RegisterType((*SessionsRequest)(nil), "api.SessionsRequest")
	proto.RegisterType((*SessionsReply)(nil), "api.SessionsReply")
	proto.RegisterType((*TopicsRequest)(nil), "api.TopicsRequest")
	proto.RegisterType((*TopicsReply)(nil), "api.TopicsReply")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Api service

type ApiClient interface {
	Version(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionReply, error)
	Admins(ctx context.Context, in *AdminsRequest, opts ...grpc.CallOption) (*AdminsReply, error)
	Cluster(ctx context.Context, in *ClusterRequest, opts ...grpc.CallOption) (*ClusterReply, error)
	Routes(ctx context.Context, in *RoutesRequest, opts ...grpc.CallOption) (*RoutesReply, error)
	Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error)
	Broker(ctx context.Context, in *BrokerRequest, opts ...grpc.CallOption) (*BrokerReply, error)
	Plugins(ctx context.Context, in *PluginsRequest, opts ...grpc.CallOption) (*PluginsReply, error)
	Services(ctx context.Context, in *ServicesRequest, opts ...grpc.CallOption) (*ServicesReply, error)
	Subscriptions(ctx context.Context, in *SubscriptionsRequest, opts ...grpc.CallOption) (*SubscriptionsReply, error)
	Clients(ctx context.Context, in *ClientsRequest, opts ...grpc.CallOption) (*ClientsReply, error)
	Sessions(ctx context.Context, in *SessionsRequest, opts ...grpc.CallOption) (*SessionsReply, error)
	Topics(ctx context.Context, in *TopicsRequest, opts ...grpc.CallOption) (*TopicsReply, error)
}

type apiClient struct {
	cc *grpc.ClientConn
}

func NewApiClient(cc *grpc.ClientConn) ApiClient {
	return &apiClient{cc}
}

func (c *apiClient) Version(ctx context.Context, in *VersionRequest, opts ...grpc.CallOption) (*VersionReply, error) {
	out := new(VersionReply)
	err := grpc.Invoke(ctx, "/api.Api/Version", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Admins(ctx context.Context, in *AdminsRequest, opts ...grpc.CallOption) (*AdminsReply, error) {
	out := new(AdminsReply)
	err := grpc.Invoke(ctx, "/api.Api/Admins", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Cluster(ctx context.Context, in *ClusterRequest, opts ...grpc.CallOption) (*ClusterReply, error) {
	out := new(ClusterReply)
	err := grpc.Invoke(ctx, "/api.Api/Cluster", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Routes(ctx context.Context, in *RoutesRequest, opts ...grpc.CallOption) (*RoutesReply, error) {
	out := new(RoutesReply)
	err := grpc.Invoke(ctx, "/api.Api/Routes", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Status(ctx context.Context, in *StatusRequest, opts ...grpc.CallOption) (*StatusReply, error) {
	out := new(StatusReply)
	err := grpc.Invoke(ctx, "/api.Api/Status", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Broker(ctx context.Context, in *BrokerRequest, opts ...grpc.CallOption) (*BrokerReply, error) {
	out := new(BrokerReply)
	err := grpc.Invoke(ctx, "/api.Api/Broker", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Plugins(ctx context.Context, in *PluginsRequest, opts ...grpc.CallOption) (*PluginsReply, error) {
	out := new(PluginsReply)
	err := grpc.Invoke(ctx, "/api.Api/Plugins", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Services(ctx context.Context, in *ServicesRequest, opts ...grpc.CallOption) (*ServicesReply, error) {
	out := new(ServicesReply)
	err := grpc.Invoke(ctx, "/api.Api/Services", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Subscriptions(ctx context.Context, in *SubscriptionsRequest, opts ...grpc.CallOption) (*SubscriptionsReply, error) {
	out := new(SubscriptionsReply)
	err := grpc.Invoke(ctx, "/api.Api/Subscriptions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Clients(ctx context.Context, in *ClientsRequest, opts ...grpc.CallOption) (*ClientsReply, error) {
	out := new(ClientsReply)
	err := grpc.Invoke(ctx, "/api.Api/Clients", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Sessions(ctx context.Context, in *SessionsRequest, opts ...grpc.CallOption) (*SessionsReply, error) {
	out := new(SessionsReply)
	err := grpc.Invoke(ctx, "/api.Api/Sessions", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *apiClient) Topics(ctx context.Context, in *TopicsRequest, opts ...grpc.CallOption) (*TopicsReply, error) {
	out := new(TopicsReply)
	err := grpc.Invoke(ctx, "/api.Api/Topics", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Api service

type ApiServer interface {
	Version(context.Context, *VersionRequest) (*VersionReply, error)
	Admins(context.Context, *AdminsRequest) (*AdminsReply, error)
	Cluster(context.Context, *ClusterRequest) (*ClusterReply, error)
	Routes(context.Context, *RoutesRequest) (*RoutesReply, error)
	Status(context.Context, *StatusRequest) (*StatusReply, error)
	Broker(context.Context, *BrokerRequest) (*BrokerReply, error)
	Plugins(context.Context, *PluginsRequest) (*PluginsReply, error)
	Services(context.Context, *ServicesRequest) (*ServicesReply, error)
	Subscriptions(context.Context, *SubscriptionsRequest) (*SubscriptionsReply, error)
	Clients(context.Context, *ClientsRequest) (*ClientsReply, error)
	Sessions(context.Context, *SessionsRequest) (*SessionsReply, error)
	Topics(context.Context, *TopicsRequest) (*TopicsReply, error)
}

func RegisterApiServer(s *grpc.Server, srv ApiServer) {
	s.RegisterService(&_Api_serviceDesc, srv)
}

func _Api_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VersionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Version(ctx, req.(*VersionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Admins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AdminsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Admins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Admins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Admins(ctx, req.(*AdminsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Cluster_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClusterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Cluster(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Cluster",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Cluster(ctx, req.(*ClusterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Routes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RoutesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Routes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Routes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Routes(ctx, req.(*RoutesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(StatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Status(ctx, req.(*StatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Broker_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(BrokerRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Broker(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Broker",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Broker(ctx, req.(*BrokerRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Plugins_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PluginsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Plugins(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Plugins",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Plugins(ctx, req.(*PluginsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Services_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ServicesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Services(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Services",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Services(ctx, req.(*ServicesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Subscriptions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SubscriptionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Subscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Subscriptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Subscriptions(ctx, req.(*SubscriptionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Clients_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ClientsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Clients(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Clients",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Clients(ctx, req.(*ClientsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Sessions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SessionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Sessions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Sessions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Sessions(ctx, req.(*SessionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Api_Topics_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(TopicsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ApiServer).Topics(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Api/Topics",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ApiServer).Topics(ctx, req.(*TopicsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Api_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Api",
	HandlerType: (*ApiServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _Api_Version_Handler,
		},
		{
			MethodName: "Admins",
			Handler:    _Api_Admins_Handler,
		},
		{
			MethodName: "Cluster",
			Handler:    _Api_Cluster_Handler,
		},
		{
			MethodName: "Routes",
			Handler:    _Api_Routes_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _Api_Status_Handler,
		},
		{
			MethodName: "Broker",
			Handler:    _Api_Broker_Handler,
		},
		{
			MethodName: "Plugins",
			Handler:    _Api_Plugins_Handler,
		},
		{
			MethodName: "Services",
			Handler:    _Api_Services_Handler,
		},
		{
			MethodName: "Subscriptions",
			Handler:    _Api_Subscriptions_Handler,
		},
		{
			MethodName: "Clients",
			Handler:    _Api_Clients_Handler,
		},
		{
			MethodName: "Sessions",
			Handler:    _Api_Sessions_Handler,
		},
		{
			MethodName: "Topics",
			Handler:    _Api_Topics_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api.proto",
}

func init() { proto.RegisterFile("api.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 410 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x6c, 0x94, 0xcd, 0x4e, 0x32, 0x31,
	0x14, 0x86, 0x21, 0xe4, 0x83, 0x8f, 0x03, 0xc3, 0x4f, 0x25, 0x8a, 0xb3, 0x32, 0xb3, 0x62, 0x45,
	0x8c, 0x24, 0xee, 0xd1, 0xb8, 0x37, 0x60, 0xdc, 0x03, 0x36, 0xa6, 0x71, 0xa4, 0x75, 0xda, 0x21,
	0xe1, 0x16, 0xbc, 0x6a, 0xd3, 0xf6, 0xb4, 0xf4, 0x10, 0x96, 0xf3, 0xf6, 0xfc, 0xbf, 0x4f, 0x06,
	0xba, 0x1b, 0x25, 0xe6, 0xaa, 0x92, 0x46, 0xb2, 0xd6, 0x46, 0x89, 0x62, 0x04, 0x83, 0x77, 0x5e,
	0x69, 0x21, 0xf7, 0x2b, 0xfe, 0x53, 0x73, 0x6d, 0x8a, 0x19, 0xf4, 0xa3, 0xa2, 0xca, 0x23, 0x9b,
	0x42, 0xe7, 0xe0, 0xbf, 0xa7, 0xcd, 0xbb, 0xe6, 0xac, 0xbb, 0x0a, 0x9f, 0xc5, 0x10, 0xb2, 0xe5,
	0xc7, 0xb7, 0xd8, 0xeb, 0x90, 0x9a, 0x41, 0x2f, 0x08, 0xaa, 0x3c, 0xda, 0xda, 0xcf, 0x65, 0xad,
	0x0d, 0xaf, 0x42, 0xc0, 0x00, 0xfa, 0x51, 0xb1, 0x11, 0x43, 0xc8, 0x56, 0xb2, 0x36, 0x3c, 0xad,
	0x10, 0x04, 0x7c, 0x5f, 0x9b, 0x8d, 0xa9, 0xd3, 0xf7, 0x20, 0xe0, 0xfb, 0x53, 0x25, 0xbf, 0x4e,
	0x0d, 0x32, 0xe8, 0x05, 0x01, 0x27, 0x78, 0x2d, 0xeb, 0xcf, 0x64, 0xc4, 0x01, 0xf4, 0xa3, 0x62,
	0x23, 0xc6, 0x30, 0x5c, 0xf3, 0xea, 0x20, 0x76, 0xa7, 0x19, 0x6c, 0xd3, 0x28, 0xd9, 0x98, 0x6b,
	0x98, 0xac, 0xeb, 0xad, 0xde, 0x55, 0x42, 0x19, 0x21, 0x4f, 0xb5, 0x26, 0xc0, 0xce, 0xf4, 0xb8,
	0xb5, 0xe0, 0x7b, 0xa3, 0xc9, 0xd6, 0xa8, 0xc4, 0x9e, 0x5a, 0xa7, 0xa5, 0x5c, 0xcf, 0x20, 0xe1,
	0x66, 0x6f, 0x52, 0x89, 0x5d, 0xba, 0x79, 0x10, 0x54, 0x79, 0x7c, 0xf8, 0xfd, 0x07, 0xad, 0xa5,
	0x12, 0x6c, 0x01, 0x1d, 0x74, 0x8b, 0x5d, 0xcd, 0xad, 0xb7, 0xd4, 0xcd, 0x7c, 0x4c, 0x45, 0x5b,
	0xba, 0xc1, 0xee, 0xa1, 0xed, 0x7d, 0x62, 0xcc, 0x3d, 0x13, 0x17, 0xf3, 0x11, 0xd1, 0x7c, 0xc6,
	0x02, 0x3a, 0x68, 0x1c, 0xb6, 0xa1, 0xc6, 0x62, 0x1b, 0xe2, 0xad, 0x6b, 0xe3, 0xcd, 0xc4, 0x36,
	0xc4, 0x6a, 0x6c, 0x93, 0xba, 0xed, 0x32, 0xbc, 0xbd, 0x98, 0x41, 0xcc, 0xc7, 0x8c, 0xd4, 0x7f,
	0x97, 0xe1, 0x0d, 0xc7, 0x0c, 0x82, 0x03, 0x66, 0xa4, 0x44, 0xb8, 0x55, 0x90, 0x00, 0x5c, 0x85,
	0x12, 0x82, 0xab, 0x10, 0x48, 0x1a, 0xec, 0x11, 0xfe, 0x07, 0x26, 0xd8, 0xc4, 0x8f, 0x41, 0xa9,
	0xc9, 0xd9, 0x99, 0xea, 0xf3, 0x5e, 0x20, 0x23, 0x88, 0xb0, 0x5b, 0x1f, 0x76, 0x01, 0xa7, 0xfc,
	0xe6, 0xd2, 0x53, 0x72, 0x7e, 0x47, 0x50, 0x3c, 0x7f, 0x4a, 0x58, 0x3c, 0x7f, 0x02, 0x19, 0xce,
	0xec, 0x99, 0x8a, 0x33, 0x13, 0xea, 0xe2, 0xcc, 0x29, 0x78, 0xee, 0xa4, 0x9e, 0x34, 0x3c, 0x29,
	0xe1, 0x10, 0x4f, 0x9a, 0xa0, 0x58, 0x34, 0xb6, 0x6d, 0xf7, 0x43, 0x59, 0xfc, 0x05, 0x00, 0x00,
	0xff, 0xff, 0x18, 0x89, 0xdf, 0x8e, 0x5d, 0x04, 0x00, 0x00,
}
