// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.25.0
// 	protoc        v3.12.3
// source: auth.proto

package protos

import (
	context "context"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// This is a compile-time assertion that a sufficiently up-to-date version
// of the legacy proto package is being used.
const _ = proto.ProtoPackageIsVersion4

type AuthReq struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Username string `protobuf:"bytes,1,opt,name=username,proto3" json:"username,omitempty"`
	Password string `protobuf:"bytes,2,opt,name=password,proto3" json:"password,omitempty"`
	// used only for Authoriation
	SessionKey   string `protobuf:"bytes,3,opt,name=sessionKey,proto3" json:"sessionKey,omitempty"`
	UserAgent    string `protobuf:"bytes,4,opt,name=userAgent,proto3" json:"userAgent,omitempty"`
	StayLoggedIn bool   `protobuf:"varint,5,opt,name=stayLoggedIn,proto3" json:"stayLoggedIn,omitempty"`
}

func (x *AuthReq) Reset() {
	*x = AuthReq{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthReq) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthReq) ProtoMessage() {}

func (x *AuthReq) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthReq.ProtoReflect.Descriptor instead.
func (*AuthReq) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{0}
}

func (x *AuthReq) GetUsername() string {
	if x != nil {
		return x.Username
	}
	return ""
}

func (x *AuthReq) GetPassword() string {
	if x != nil {
		return x.Password
	}
	return ""
}

func (x *AuthReq) GetSessionKey() string {
	if x != nil {
		return x.SessionKey
	}
	return ""
}

func (x *AuthReq) GetUserAgent() string {
	if x != nil {
		return x.UserAgent
	}
	return ""
}

func (x *AuthReq) GetStayLoggedIn() bool {
	if x != nil {
		return x.StayLoggedIn
	}
	return false
}

// AuthRes contains the status of the request
// success == true : session key and user data will be populated
type AuthRes struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Success    bool   `protobuf:"varint,1,opt,name=success,proto3" json:"success,omitempty"`
	SessionKey string `protobuf:"bytes,2,opt,name=sessionKey,proto3" json:"sessionKey,omitempty"`
	User       *User  `protobuf:"bytes,15,opt,name=user,proto3" json:"user,omitempty"`
}

func (x *AuthRes) Reset() {
	*x = AuthRes{}
	if protoimpl.UnsafeEnabled {
		mi := &file_auth_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AuthRes) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AuthRes) ProtoMessage() {}

func (x *AuthRes) ProtoReflect() protoreflect.Message {
	mi := &file_auth_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AuthRes.ProtoReflect.Descriptor instead.
func (*AuthRes) Descriptor() ([]byte, []int) {
	return file_auth_proto_rawDescGZIP(), []int{1}
}

func (x *AuthRes) GetSuccess() bool {
	if x != nil {
		return x.Success
	}
	return false
}

func (x *AuthRes) GetSessionKey() string {
	if x != nil {
		return x.SessionKey
	}
	return ""
}

func (x *AuthRes) GetUser() *User {
	if x != nil {
		return x.User
	}
	return nil
}

var File_auth_proto protoreflect.FileDescriptor

var file_auth_proto_rawDesc = []byte{
	0x0a, 0x0a, 0x61, 0x75, 0x74, 0x68, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x06, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x73, 0x1a, 0x0a, 0x75, 0x73, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x22, 0xa3, 0x01, 0x0a, 0x07, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x12, 0x1a, 0x0a, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08,
	0x75, 0x73, 0x65, 0x72, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x70, 0x61, 0x73, 0x73,
	0x77, 0x6f, 0x72, 0x64, 0x12, 0x1e, 0x0a, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4b,
	0x65, 0x79, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f,
	0x6e, 0x4b, 0x65, 0x79, 0x12, 0x1c, 0x0a, 0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65, 0x6e,
	0x74, 0x18, 0x04, 0x20, 0x01, 0x28, 0x09, 0x52, 0x09, 0x75, 0x73, 0x65, 0x72, 0x41, 0x67, 0x65,
	0x6e, 0x74, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x74, 0x61, 0x79, 0x4c, 0x6f, 0x67, 0x67, 0x65, 0x64,
	0x49, 0x6e, 0x18, 0x05, 0x20, 0x01, 0x28, 0x08, 0x52, 0x0c, 0x73, 0x74, 0x61, 0x79, 0x4c, 0x6f,
	0x67, 0x67, 0x65, 0x64, 0x49, 0x6e, 0x22, 0x65, 0x0a, 0x07, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65,
	0x73, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x18, 0x01, 0x20, 0x01,
	0x28, 0x08, 0x52, 0x07, 0x73, 0x75, 0x63, 0x63, 0x65, 0x73, 0x73, 0x12, 0x1e, 0x0a, 0x0a, 0x73,
	0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x0a, 0x73, 0x65, 0x73, 0x73, 0x69, 0x6f, 0x6e, 0x4b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x04, 0x75,
	0x73, 0x65, 0x72, 0x18, 0x0f, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x73, 0x2e, 0x55, 0x73, 0x65, 0x72, 0x52, 0x04, 0x75, 0x73, 0x65, 0x72, 0x32, 0x99, 0x01,
	0x0a, 0x04, 0x41, 0x75, 0x74, 0x68, 0x12, 0x32, 0x0a, 0x0c, 0x41, 0x75, 0x74, 0x68, 0x65, 0x6e,
	0x74, 0x69, 0x63, 0x61, 0x74, 0x65, 0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x2f, 0x0a, 0x09, 0x41, 0x75,
	0x74, 0x68, 0x6f, 0x72, 0x69, 0x7a, 0x65, 0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73,
	0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x73, 0x2e, 0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x06, 0x4c,
	0x6f, 0x67, 0x6f, 0x75, 0x74, 0x12, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e, 0x41,
	0x75, 0x74, 0x68, 0x52, 0x65, 0x71, 0x1a, 0x0f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x73, 0x2e,
	0x41, 0x75, 0x74, 0x68, 0x52, 0x65, 0x73, 0x22, 0x00, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x3b, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x73, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_auth_proto_rawDescOnce sync.Once
	file_auth_proto_rawDescData = file_auth_proto_rawDesc
)

func file_auth_proto_rawDescGZIP() []byte {
	file_auth_proto_rawDescOnce.Do(func() {
		file_auth_proto_rawDescData = protoimpl.X.CompressGZIP(file_auth_proto_rawDescData)
	})
	return file_auth_proto_rawDescData
}

var file_auth_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_auth_proto_goTypes = []interface{}{
	(*AuthReq)(nil), // 0: protos.AuthReq
	(*AuthRes)(nil), // 1: protos.AuthRes
	(*User)(nil),    // 2: protos.User
}
var file_auth_proto_depIdxs = []int32{
	2, // 0: protos.AuthRes.user:type_name -> protos.User
	0, // 1: protos.Auth.Authenticate:input_type -> protos.AuthReq
	0, // 2: protos.Auth.Authorize:input_type -> protos.AuthReq
	0, // 3: protos.Auth.Logout:input_type -> protos.AuthReq
	1, // 4: protos.Auth.Authenticate:output_type -> protos.AuthRes
	1, // 5: protos.Auth.Authorize:output_type -> protos.AuthRes
	1, // 6: protos.Auth.Logout:output_type -> protos.AuthRes
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_auth_proto_init() }
func file_auth_proto_init() {
	if File_auth_proto != nil {
		return
	}
	file_user_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_auth_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthReq); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_auth_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AuthRes); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_auth_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_auth_proto_goTypes,
		DependencyIndexes: file_auth_proto_depIdxs,
		MessageInfos:      file_auth_proto_msgTypes,
	}.Build()
	File_auth_proto = out.File
	file_auth_proto_rawDesc = nil
	file_auth_proto_goTypes = nil
	file_auth_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// AuthClient is the client API for Auth service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type AuthClient interface {
	Authenticate(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthRes, error)
	Authorize(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthRes, error)
	Logout(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthRes, error)
}

type authClient struct {
	cc grpc.ClientConnInterface
}

func NewAuthClient(cc grpc.ClientConnInterface) AuthClient {
	return &authClient{cc}
}

func (c *authClient) Authenticate(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthRes, error) {
	out := new(AuthRes)
	err := c.cc.Invoke(ctx, "/protos.Auth/Authenticate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Authorize(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthRes, error) {
	out := new(AuthRes)
	err := c.cc.Invoke(ctx, "/protos.Auth/Authorize", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *authClient) Logout(ctx context.Context, in *AuthReq, opts ...grpc.CallOption) (*AuthRes, error) {
	out := new(AuthRes)
	err := c.cc.Invoke(ctx, "/protos.Auth/Logout", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AuthServer is the server API for Auth service.
type AuthServer interface {
	Authenticate(context.Context, *AuthReq) (*AuthRes, error)
	Authorize(context.Context, *AuthReq) (*AuthRes, error)
	Logout(context.Context, *AuthReq) (*AuthRes, error)
}

// UnimplementedAuthServer can be embedded to have forward compatible implementations.
type UnimplementedAuthServer struct {
}

func (*UnimplementedAuthServer) Authenticate(context.Context, *AuthReq) (*AuthRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authenticate not implemented")
}
func (*UnimplementedAuthServer) Authorize(context.Context, *AuthReq) (*AuthRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Authorize not implemented")
}
func (*UnimplementedAuthServer) Logout(context.Context, *AuthReq) (*AuthRes, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}

func RegisterAuthServer(s *grpc.Server, srv AuthServer) {
	s.RegisterService(&_Auth_serviceDesc, srv)
}

func _Auth_Authenticate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Authenticate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Auth/Authenticate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Authenticate(ctx, req.(*AuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Authorize_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Authorize(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Auth/Authorize",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Authorize(ctx, req.(*AuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

func _Auth_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthReq)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AuthServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/protos.Auth/Logout",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AuthServer).Logout(ctx, req.(*AuthReq))
	}
	return interceptor(ctx, in, info, handler)
}

var _Auth_serviceDesc = grpc.ServiceDesc{
	ServiceName: "protos.Auth",
	HandlerType: (*AuthServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Authenticate",
			Handler:    _Auth_Authenticate_Handler,
		},
		{
			MethodName: "Authorize",
			Handler:    _Auth_Authorize_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _Auth_Logout_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "auth.proto",
}
