// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.12.4
// source: lider.proto

package lider

import (
	context "context"
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

type Interaction struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	PlayerId int32  `protobuf:"varint,1,opt,name=player_id,json=playerId,proto3" json:"player_id,omitempty"` // player id
	Play     int32  `protobuf:"varint,2,opt,name=play,proto3" json:"play,omitempty"`                         // content of the play
	Response string `protobuf:"bytes,3,opt,name=response,proto3" json:"response,omitempty"`                  // response to the play
}

func (x *Interaction) Reset() {
	*x = Interaction{}
	if protoimpl.UnsafeEnabled {
		mi := &file_lider_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Interaction) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Interaction) ProtoMessage() {}

func (x *Interaction) ProtoReflect() protoreflect.Message {
	mi := &file_lider_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Interaction.ProtoReflect.Descriptor instead.
func (*Interaction) Descriptor() ([]byte, []int) {
	return file_lider_proto_rawDescGZIP(), []int{0}
}

func (x *Interaction) GetPlayerId() int32 {
	if x != nil {
		return x.PlayerId
	}
	return 0
}

func (x *Interaction) GetPlay() int32 {
	if x != nil {
		return x.Play
	}
	return 0
}

func (x *Interaction) GetResponse() string {
	if x != nil {
		return x.Response
	}
	return ""
}

var File_lider_proto protoreflect.FileDescriptor

var file_lider_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x6c,
	0x69, 0x64, 0x65, 0x72, 0x22, 0x5a, 0x0a, 0x0b, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x1b, 0x0a, 0x09, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x5f, 0x69, 0x64,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x08, 0x70, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x49, 0x64,
	0x12, 0x12, 0x0a, 0x04, 0x70, 0x6c, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x04,
	0x70, 0x6c, 0x61, 0x79, 0x12, 0x1a, 0x0a, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x08, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x32, 0xba, 0x01, 0x0a, 0x0d, 0x50, 0x6c, 0x61, 0x79, 0x65, 0x72, 0x53, 0x65, 0x72, 0x76, 0x69,
	0x63, 0x65, 0x12, 0x39, 0x0a, 0x0d, 0x47, 0x65, 0x74, 0x43, 0x6f, 0x6e, 0x6e, 0x65, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x12, 0x12, 0x2e, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x12, 0x2e, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x2e,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x12, 0x34, 0x0a,
	0x08, 0x53, 0x65, 0x6e, 0x64, 0x50, 0x6c, 0x61, 0x79, 0x12, 0x12, 0x2e, 0x6c, 0x69, 0x64, 0x65,
	0x72, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x12, 0x2e,
	0x6c, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x22, 0x00, 0x12, 0x38, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x50, 0x72, 0x69, 0x7a, 0x65, 0x50,
	0x6f, 0x6f, 0x6c, 0x12, 0x12, 0x2e, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x2e, 0x49, 0x6e, 0x74, 0x65,
	0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x1a, 0x12, 0x2e, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x2e,
	0x49, 0x6e, 0x74, 0x65, 0x72, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e, 0x22, 0x00, 0x42, 0x08, 0x5a,
	0x06, 0x2f, 0x6c, 0x69, 0x64, 0x65, 0x72, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_lider_proto_rawDescOnce sync.Once
	file_lider_proto_rawDescData = file_lider_proto_rawDesc
)

func file_lider_proto_rawDescGZIP() []byte {
	file_lider_proto_rawDescOnce.Do(func() {
		file_lider_proto_rawDescData = protoimpl.X.CompressGZIP(file_lider_proto_rawDescData)
	})
	return file_lider_proto_rawDescData
}

var file_lider_proto_msgTypes = make([]protoimpl.MessageInfo, 1)
var file_lider_proto_goTypes = []interface{}{
	(*Interaction)(nil), // 0: lider.Interaction
}
var file_lider_proto_depIdxs = []int32{
	0, // 0: lider.PlayerService.GetConnection:input_type -> lider.Interaction
	0, // 1: lider.PlayerService.SendPlay:input_type -> lider.Interaction
	0, // 2: lider.PlayerService.GetPrizePool:input_type -> lider.Interaction
	0, // 3: lider.PlayerService.GetConnection:output_type -> lider.Interaction
	0, // 4: lider.PlayerService.SendPlay:output_type -> lider.Interaction
	0, // 5: lider.PlayerService.GetPrizePool:output_type -> lider.Interaction
	3, // [3:6] is the sub-list for method output_type
	0, // [0:3] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_lider_proto_init() }
func file_lider_proto_init() {
	if File_lider_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_lider_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Interaction); i {
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
			RawDescriptor: file_lider_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   1,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_lider_proto_goTypes,
		DependencyIndexes: file_lider_proto_depIdxs,
		MessageInfos:      file_lider_proto_msgTypes,
	}.Build()
	File_lider_proto = out.File
	file_lider_proto_rawDesc = nil
	file_lider_proto_goTypes = nil
	file_lider_proto_depIdxs = nil
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// PlayerServiceClient is the client API for PlayerService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type PlayerServiceClient interface {
	GetConnection(ctx context.Context, in *Interaction, opts ...grpc.CallOption) (*Interaction, error)
	SendPlay(ctx context.Context, in *Interaction, opts ...grpc.CallOption) (*Interaction, error)
	GetPrizePool(ctx context.Context, in *Interaction, opts ...grpc.CallOption) (*Interaction, error)
}

type playerServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPlayerServiceClient(cc grpc.ClientConnInterface) PlayerServiceClient {
	return &playerServiceClient{cc}
}

func (c *playerServiceClient) GetConnection(ctx context.Context, in *Interaction, opts ...grpc.CallOption) (*Interaction, error) {
	out := new(Interaction)
	err := c.cc.Invoke(ctx, "/lider.PlayerService/GetConnection", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerServiceClient) SendPlay(ctx context.Context, in *Interaction, opts ...grpc.CallOption) (*Interaction, error) {
	out := new(Interaction)
	err := c.cc.Invoke(ctx, "/lider.PlayerService/SendPlay", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playerServiceClient) GetPrizePool(ctx context.Context, in *Interaction, opts ...grpc.CallOption) (*Interaction, error) {
	out := new(Interaction)
	err := c.cc.Invoke(ctx, "/lider.PlayerService/GetPrizePool", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlayerServiceServer is the server API for PlayerService service.
type PlayerServiceServer interface {
	GetConnection(context.Context, *Interaction) (*Interaction, error)
	SendPlay(context.Context, *Interaction) (*Interaction, error)
	GetPrizePool(context.Context, *Interaction) (*Interaction, error)
}

// UnimplementedPlayerServiceServer can be embedded to have forward compatible implementations.
type UnimplementedPlayerServiceServer struct {
}

func (*UnimplementedPlayerServiceServer) GetConnection(context.Context, *Interaction) (*Interaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetConnection not implemented")
}
func (*UnimplementedPlayerServiceServer) SendPlay(context.Context, *Interaction) (*Interaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendPlay not implemented")
}
func (*UnimplementedPlayerServiceServer) GetPrizePool(context.Context, *Interaction) (*Interaction, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPrizePool not implemented")
}

func RegisterPlayerServiceServer(s *grpc.Server, srv PlayerServiceServer) {
	s.RegisterService(&_PlayerService_serviceDesc, srv)
}

func _PlayerService_GetConnection_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Interaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServiceServer).GetConnection(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lider.PlayerService/GetConnection",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServiceServer).GetConnection(ctx, req.(*Interaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerService_SendPlay_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Interaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServiceServer).SendPlay(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lider.PlayerService/SendPlay",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServiceServer).SendPlay(ctx, req.(*Interaction))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlayerService_GetPrizePool_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Interaction)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlayerServiceServer).GetPrizePool(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/lider.PlayerService/GetPrizePool",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlayerServiceServer).GetPrizePool(ctx, req.(*Interaction))
	}
	return interceptor(ctx, in, info, handler)
}

var _PlayerService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "lider.PlayerService",
	HandlerType: (*PlayerServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetConnection",
			Handler:    _PlayerService_GetConnection_Handler,
		},
		{
			MethodName: "SendPlay",
			Handler:    _PlayerService_SendPlay_Handler,
		},
		{
			MethodName: "GetPrizePool",
			Handler:    _PlayerService_GetPrizePool_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "lider.proto",
}
