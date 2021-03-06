// Code generated by protoc-gen-go. DO NOT EDIT.
// source: api/protobuf-spec/backend.proto

/*
Package pb is a generated protocol buffer package.

It is generated from these files:
	api/protobuf-spec/backend.proto
	api/protobuf-spec/frontend.proto
	api/protobuf-spec/mmlogic.proto
	api/protobuf-spec/messages.proto

It has these top-level messages:
	Group
	PlayerId
	MatchObject
	Roster
	Filter
	Stats
	PlayerPool
	Player
	Result
	IlInput
	ConnectionInfo
	Assignments
*/
package pb

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

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Backend service

type BackendClient interface {
	// Run MMF once.  Return a matchobject that fits this profile.
	// INPUT: MatchObject message with these fields populated:
	//  - id
	//  - properties
	//  - [optional] roster, any fields you fill are available to your MMF.
	//  - [optional] pools, any fields you fill are available to your MMF.
	// OUTPUT: MatchObject message with these fields populated:
	//  - id
	//  - properties
	//  - error. Empty if no error was encountered
	//  - rosters, if you choose to fill them in your MMF. (Recommended)
	//  - pools, if you used the MMLogicAPI in your MMF. (Recommended, and provides stats)
	CreateMatch(ctx context.Context, in *MatchObject, opts ...grpc.CallOption) (*MatchObject, error)
	// Continually run MMF and stream matchobjects that fit this profile until
	// client closes the connection.  Same inputs/outputs as CreateMatch.
	ListMatches(ctx context.Context, in *MatchObject, opts ...grpc.CallOption) (Backend_ListMatchesClient, error)
	// Delete a matchobject from state storage manually. (Matchobjects in state
	// storage will also automatically expire after a while)
	// INPUT: MatchObject message with the 'id' field populated.
	// (All other fields are ignored.)
	DeleteMatch(ctx context.Context, in *MatchObject, opts ...grpc.CallOption) (*Result, error)
	// Write the connection info for the list of players in the
	// Assignments.messages.Rosters to state storage.  The FrontendAPI is
	// responsible for sending anything sent here to the game clients.
	// Sending a player to this function kicks off a process that removes
	// the player from future matchmaking functions by adding them to the
	// 'deindexed' player list and then deleting their player ID from state storage
	// indexes.
	// INPUT: Assignments message with these fields populated:
	//  - connection_info, anything you write to this string is sent to Frontend API
	//  - rosters. You can send any number of rosters, containing any number of
	//     player messages. All players from all rosters will be sent the connection_info.
	//     The only field in the Player object that is used by CreateAssignments is
	//     the id field.  All others are silently ignored.
	CreateAssignments(ctx context.Context, in *Assignments, opts ...grpc.CallOption) (*Result, error)
	// Remove DGS connection info from state storage for players.
	// INPUT: Roster message with the 'players' field populated.
	//    The only field in the Player object that is used by
	//    DeleteAssignments is the 'id' field.  All others are silently ignored.  If
	//    you need to delete multiple rosters, make multiple calls.
	DeleteAssignments(ctx context.Context, in *Roster, opts ...grpc.CallOption) (*Result, error)
}

type backendClient struct {
	cc *grpc.ClientConn
}

func NewBackendClient(cc *grpc.ClientConn) BackendClient {
	return &backendClient{cc}
}

func (c *backendClient) CreateMatch(ctx context.Context, in *MatchObject, opts ...grpc.CallOption) (*MatchObject, error) {
	out := new(MatchObject)
	err := grpc.Invoke(ctx, "/api.Backend/CreateMatch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) ListMatches(ctx context.Context, in *MatchObject, opts ...grpc.CallOption) (Backend_ListMatchesClient, error) {
	stream, err := grpc.NewClientStream(ctx, &_Backend_serviceDesc.Streams[0], c.cc, "/api.Backend/ListMatches", opts...)
	if err != nil {
		return nil, err
	}
	x := &backendListMatchesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type Backend_ListMatchesClient interface {
	Recv() (*MatchObject, error)
	grpc.ClientStream
}

type backendListMatchesClient struct {
	grpc.ClientStream
}

func (x *backendListMatchesClient) Recv() (*MatchObject, error) {
	m := new(MatchObject)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *backendClient) DeleteMatch(ctx context.Context, in *MatchObject, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/api.Backend/DeleteMatch", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) CreateAssignments(ctx context.Context, in *Assignments, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/api.Backend/CreateAssignments", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *backendClient) DeleteAssignments(ctx context.Context, in *Roster, opts ...grpc.CallOption) (*Result, error) {
	out := new(Result)
	err := grpc.Invoke(ctx, "/api.Backend/DeleteAssignments", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Backend service

type BackendServer interface {
	// Run MMF once.  Return a matchobject that fits this profile.
	// INPUT: MatchObject message with these fields populated:
	//  - id
	//  - properties
	//  - [optional] roster, any fields you fill are available to your MMF.
	//  - [optional] pools, any fields you fill are available to your MMF.
	// OUTPUT: MatchObject message with these fields populated:
	//  - id
	//  - properties
	//  - error. Empty if no error was encountered
	//  - rosters, if you choose to fill them in your MMF. (Recommended)
	//  - pools, if you used the MMLogicAPI in your MMF. (Recommended, and provides stats)
	CreateMatch(context.Context, *MatchObject) (*MatchObject, error)
	// Continually run MMF and stream matchobjects that fit this profile until
	// client closes the connection.  Same inputs/outputs as CreateMatch.
	ListMatches(*MatchObject, Backend_ListMatchesServer) error
	// Delete a matchobject from state storage manually. (Matchobjects in state
	// storage will also automatically expire after a while)
	// INPUT: MatchObject message with the 'id' field populated.
	// (All other fields are ignored.)
	DeleteMatch(context.Context, *MatchObject) (*Result, error)
	// Write the connection info for the list of players in the
	// Assignments.messages.Rosters to state storage.  The FrontendAPI is
	// responsible for sending anything sent here to the game clients.
	// Sending a player to this function kicks off a process that removes
	// the player from future matchmaking functions by adding them to the
	// 'deindexed' player list and then deleting their player ID from state storage
	// indexes.
	// INPUT: Assignments message with these fields populated:
	//  - connection_info, anything you write to this string is sent to Frontend API
	//  - rosters. You can send any number of rosters, containing any number of
	//     player messages. All players from all rosters will be sent the connection_info.
	//     The only field in the Player object that is used by CreateAssignments is
	//     the id field.  All others are silently ignored.
	CreateAssignments(context.Context, *Assignments) (*Result, error)
	// Remove DGS connection info from state storage for players.
	// INPUT: Roster message with the 'players' field populated.
	//    The only field in the Player object that is used by
	//    DeleteAssignments is the 'id' field.  All others are silently ignored.  If
	//    you need to delete multiple rosters, make multiple calls.
	DeleteAssignments(context.Context, *Roster) (*Result, error)
}

func RegisterBackendServer(s *grpc.Server, srv BackendServer) {
	s.RegisterService(&_Backend_serviceDesc, srv)
}

func _Backend_CreateMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchObject)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).CreateMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Backend/CreateMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).CreateMatch(ctx, req.(*MatchObject))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_ListMatches_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(MatchObject)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(BackendServer).ListMatches(m, &backendListMatchesServer{stream})
}

type Backend_ListMatchesServer interface {
	Send(*MatchObject) error
	grpc.ServerStream
}

type backendListMatchesServer struct {
	grpc.ServerStream
}

func (x *backendListMatchesServer) Send(m *MatchObject) error {
	return x.ServerStream.SendMsg(m)
}

func _Backend_DeleteMatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(MatchObject)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).DeleteMatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Backend/DeleteMatch",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).DeleteMatch(ctx, req.(*MatchObject))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_CreateAssignments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Assignments)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).CreateAssignments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Backend/CreateAssignments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).CreateAssignments(ctx, req.(*Assignments))
	}
	return interceptor(ctx, in, info, handler)
}

func _Backend_DeleteAssignments_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Roster)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BackendServer).DeleteAssignments(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.Backend/DeleteAssignments",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BackendServer).DeleteAssignments(ctx, req.(*Roster))
	}
	return interceptor(ctx, in, info, handler)
}

var _Backend_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.Backend",
	HandlerType: (*BackendServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateMatch",
			Handler:    _Backend_CreateMatch_Handler,
		},
		{
			MethodName: "DeleteMatch",
			Handler:    _Backend_DeleteMatch_Handler,
		},
		{
			MethodName: "CreateAssignments",
			Handler:    _Backend_CreateAssignments_Handler,
		},
		{
			MethodName: "DeleteAssignments",
			Handler:    _Backend_DeleteAssignments_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ListMatches",
			Handler:       _Backend_ListMatches_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/protobuf-spec/backend.proto",
}

func init() { proto.RegisterFile("api/protobuf-spec/backend.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 240 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x91, 0x3f, 0x4b, 0x04, 0x31,
	0x10, 0x47, 0xfd, 0x03, 0x0a, 0xd9, 0xc6, 0x0b, 0xd8, 0x5c, 0xa3, 0xd8, 0xdf, 0x46, 0x14, 0x51,
	0x0b, 0x15, 0xef, 0x04, 0x1b, 0x45, 0xb1, 0xb4, 0x4b, 0x72, 0x73, 0x7b, 0xd1, 0x24, 0x13, 0x32,
	0xb3, 0x5f, 0xcf, 0xcf, 0x26, 0xbb, 0x01, 0x5d, 0x71, 0x41, 0x6c, 0x1f, 0xef, 0xf1, 0x9b, 0x10,
	0x71, 0xa0, 0x93, 0x53, 0x29, 0x23, 0xa3, 0x69, 0x57, 0x33, 0x4a, 0x60, 0x95, 0xd1, 0xf6, 0x1d,
	0xe2, 0xb2, 0xee, 0xa9, 0xdc, 0xd6, 0xc9, 0x4d, 0x0f, 0x7f, 0x5b, 0x01, 0x88, 0x74, 0x03, 0x54,
	0xb4, 0x93, 0x8f, 0x2d, 0xb1, 0x3b, 0x2f, 0xa1, 0xbc, 0x12, 0xd5, 0x22, 0x83, 0x66, 0x78, 0xd4,
	0x6c, 0xd7, 0x72, 0xbf, 0xfe, 0x72, 0x7b, 0xf0, 0x64, 0xde, 0xc0, 0xf2, 0x74, 0x1c, 0x1f, 0x6d,
	0xc8, 0x1b, 0x51, 0x3d, 0x38, 0xe2, 0x1e, 0x02, 0xfd, 0x37, 0x3f, 0xde, 0x94, 0x17, 0xa2, 0xba,
	0x03, 0x0f, 0x7f, 0xec, 0xef, 0x7d, 0xe3, 0x17, 0xa0, 0xd6, 0x77, 0xd3, 0xd7, 0x62, 0x52, 0x2e,
	0xbf, 0x25, 0x72, 0x4d, 0x0c, 0x10, 0xf9, 0xc7, 0x01, 0x03, 0x3c, 0xda, 0x5f, 0x8a, 0x49, 0x59,
	0x1e, 0xf6, 0x43, 0x11, 0x89, 0x21, 0x8f, 0xa5, 0xf3, 0xf3, 0xd7, 0xb3, 0xc6, 0xf1, 0xba, 0x35,
	0xb5, 0xc5, 0xa0, 0xee, 0x11, 0x1b, 0x0f, 0x0b, 0x8f, 0xed, 0xf2, 0xd9, 0x6b, 0x5e, 0x61, 0x0e,
	0x0a, 0x13, 0xc4, 0x59, 0xe8, 0x5e, 0xa0, 0x5c, 0x64, 0xc8, 0x51, 0x7b, 0x95, 0x8c, 0xd9, 0xe9,
	0x3f, 0xe0, 0xf4, 0x33, 0x00, 0x00, 0xff, 0xff, 0x33, 0x3a, 0x1e, 0x0d, 0xca, 0x01, 0x00, 0x00,
}
