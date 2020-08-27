// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package protos

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// PodClient is the client API for Pod service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PodClient interface {
	GetEpisodes(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Episodes, error)
	GetUserEpisode(ctx context.Context, in *Request, opts ...grpc.CallOption) (*UserEpisode, error)
	UpdateUserEpisode(ctx context.Context, in *UserEpisodeReq, opts ...grpc.CallOption) (*Response, error)
	GetSubscriptions(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Subscriptions, error)
	GetUserLastPlayed(ctx context.Context, in *Request, opts ...grpc.CallOption) (*LastPlayedRes, error)
}

type podClient struct {
	cc grpc.ClientConnInterface
}

func NewPodClient(cc grpc.ClientConnInterface) PodClient {
	return &podClient{cc}
}

var podGetEpisodesStreamDesc = &grpc.StreamDesc{
	StreamName: "GetEpisodes",
}

func (c *podClient) GetEpisodes(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Episodes, error) {
	out := new(Episodes)
	err := c.cc.Invoke(ctx, "/protos.Pod/GetEpisodes", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var podGetUserEpisodeStreamDesc = &grpc.StreamDesc{
	StreamName: "GetUserEpisode",
}

func (c *podClient) GetUserEpisode(ctx context.Context, in *Request, opts ...grpc.CallOption) (*UserEpisode, error) {
	out := new(UserEpisode)
	err := c.cc.Invoke(ctx, "/protos.Pod/GetUserEpisode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var podUpdateUserEpisodeStreamDesc = &grpc.StreamDesc{
	StreamName: "UpdateUserEpisode",
}

func (c *podClient) UpdateUserEpisode(ctx context.Context, in *UserEpisodeReq, opts ...grpc.CallOption) (*Response, error) {
	out := new(Response)
	err := c.cc.Invoke(ctx, "/protos.Pod/UpdateUserEpisode", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var podGetSubscriptionsStreamDesc = &grpc.StreamDesc{
	StreamName: "GetSubscriptions",
}

func (c *podClient) GetSubscriptions(ctx context.Context, in *Request, opts ...grpc.CallOption) (*Subscriptions, error) {
	out := new(Subscriptions)
	err := c.cc.Invoke(ctx, "/protos.Pod/GetSubscriptions", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

var podGetUserLastPlayedStreamDesc = &grpc.StreamDesc{
	StreamName: "GetUserLastPlayed",
}

func (c *podClient) GetUserLastPlayed(ctx context.Context, in *Request, opts ...grpc.CallOption) (*LastPlayedRes, error) {
	out := new(LastPlayedRes)
	err := c.cc.Invoke(ctx, "/protos.Pod/GetUserLastPlayed", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PodService is the service API for Pod service.
// Fields should be assigned to their respective handler implementations only before
// RegisterPodService is called.  Any unassigned fields will result in the
// handler for that method returning an Unimplemented error.
type PodService struct {
	GetEpisodes       func(context.Context, *Request) (*Episodes, error)
	GetUserEpisode    func(context.Context, *Request) (*UserEpisode, error)
	UpdateUserEpisode func(context.Context, *UserEpisodeReq) (*Response, error)
	GetSubscriptions  func(context.Context, *Request) (*Subscriptions, error)
	GetUserLastPlayed func(context.Context, *Request) (*LastPlayedRes, error)
}

func (s *PodService) getEpisodes(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetEpisodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/protos.Pod/GetEpisodes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetEpisodes(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *PodService) getUserEpisode(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetUserEpisode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/protos.Pod/GetUserEpisode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetUserEpisode(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *PodService) updateUserEpisode(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserEpisodeReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.UpdateUserEpisode(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/protos.Pod/UpdateUserEpisode",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.UpdateUserEpisode(ctx, req.(*UserEpisodeReq))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *PodService) getSubscriptions(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetSubscriptions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/protos.Pod/GetSubscriptions",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetSubscriptions(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}
func (s *PodService) getUserLastPlayed(_ interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Request)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return s.GetUserLastPlayed(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     s,
		FullMethod: "/protos.Pod/GetUserLastPlayed",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return s.GetUserLastPlayed(ctx, req.(*Request))
	}
	return interceptor(ctx, in, info, handler)
}

// RegisterPodService registers a service implementation with a gRPC server.
func RegisterPodService(s grpc.ServiceRegistrar, srv *PodService) {
	srvCopy := *srv
	if srvCopy.GetEpisodes == nil {
		srvCopy.GetEpisodes = func(context.Context, *Request) (*Episodes, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetEpisodes not implemented")
		}
	}
	if srvCopy.GetUserEpisode == nil {
		srvCopy.GetUserEpisode = func(context.Context, *Request) (*UserEpisode, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetUserEpisode not implemented")
		}
	}
	if srvCopy.UpdateUserEpisode == nil {
		srvCopy.UpdateUserEpisode = func(context.Context, *UserEpisodeReq) (*Response, error) {
			return nil, status.Errorf(codes.Unimplemented, "method UpdateUserEpisode not implemented")
		}
	}
	if srvCopy.GetSubscriptions == nil {
		srvCopy.GetSubscriptions = func(context.Context, *Request) (*Subscriptions, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetSubscriptions not implemented")
		}
	}
	if srvCopy.GetUserLastPlayed == nil {
		srvCopy.GetUserLastPlayed = func(context.Context, *Request) (*LastPlayedRes, error) {
			return nil, status.Errorf(codes.Unimplemented, "method GetUserLastPlayed not implemented")
		}
	}
	sd := grpc.ServiceDesc{
		ServiceName: "protos.Pod",
		Methods: []grpc.MethodDesc{
			{
				MethodName: "GetEpisodes",
				Handler:    srvCopy.getEpisodes,
			},
			{
				MethodName: "GetUserEpisode",
				Handler:    srvCopy.getUserEpisode,
			},
			{
				MethodName: "UpdateUserEpisode",
				Handler:    srvCopy.updateUserEpisode,
			},
			{
				MethodName: "GetSubscriptions",
				Handler:    srvCopy.getSubscriptions,
			},
			{
				MethodName: "GetUserLastPlayed",
				Handler:    srvCopy.getUserLastPlayed,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "podcast.proto",
	}

	s.RegisterService(&sd, nil)
}

// NewPodService creates a new PodService containing the
// implemented methods of the Pod service in s.  Any unimplemented
// methods will result in the gRPC server returning an UNIMPLEMENTED status to the client.
// This includes situations where the method handler is misspelled or has the wrong
// signature.  For this reason, this function should be used with great care and
// is not recommended to be used by most users.
func NewPodService(s interface{}) *PodService {
	ns := &PodService{}
	if h, ok := s.(interface {
		GetEpisodes(context.Context, *Request) (*Episodes, error)
	}); ok {
		ns.GetEpisodes = h.GetEpisodes
	}
	if h, ok := s.(interface {
		GetUserEpisode(context.Context, *Request) (*UserEpisode, error)
	}); ok {
		ns.GetUserEpisode = h.GetUserEpisode
	}
	if h, ok := s.(interface {
		UpdateUserEpisode(context.Context, *UserEpisodeReq) (*Response, error)
	}); ok {
		ns.UpdateUserEpisode = h.UpdateUserEpisode
	}
	if h, ok := s.(interface {
		GetSubscriptions(context.Context, *Request) (*Subscriptions, error)
	}); ok {
		ns.GetSubscriptions = h.GetSubscriptions
	}
	if h, ok := s.(interface {
		GetUserLastPlayed(context.Context, *Request) (*LastPlayedRes, error)
	}); ok {
		ns.GetUserLastPlayed = h.GetUserLastPlayed
	}
	return ns
}

// UnstablePodService is the service API for Pod service.
// New methods may be added to this interface if they are added to the service
// definition, which is not a backward-compatible change.  For this reason,
// use of this type is not recommended.
type UnstablePodService interface {
	GetEpisodes(context.Context, *Request) (*Episodes, error)
	GetUserEpisode(context.Context, *Request) (*UserEpisode, error)
	UpdateUserEpisode(context.Context, *UserEpisodeReq) (*Response, error)
	GetSubscriptions(context.Context, *Request) (*Subscriptions, error)
	GetUserLastPlayed(context.Context, *Request) (*LastPlayedRes, error)
}