// Code generated by protoc-gen-dora. DO NOT EDIT.
// source: kv.proto

package kv

import (
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	math "math"
)

import (
	context "context"
	client "github.com/poonman/entry-task/dora/client"
	server "github.com/poonman/entry-task/dora/server"
	status "github.com/poonman/entry-task/dora/status"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type KVClient interface {
	SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error)
}

type kVClient struct {
	cc client.Invoker
}

func (c *kVClient) SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error) {
	out = &HelloRsp{}

	err = c.cc.Invoke(ctx, "SayHello", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type KVServer interface {
	SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error)
	mustEmbedUnimplementedKVServer()
}

type UnimplementedKVServer struct {
}

func (UnimplementedKVServer) SayHello(context.Context, *HelloReq) (*HelloRsp, error) {
	return nil, status.New(status.Unimplemented, "server is unimplemented")
}

func (UnimplementedKVServer) mustEmbedUnimplementedKVServer() {}

func RegisterKVServer(r server.ServiceRegistrar, impl KVServer) {
	r.RegisterService(KV_ServiceDesc, impl)
}
func _KVServer_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor server.Interceptor) (interface{}, error) {
	in := new(HelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(KVServer).SayHello(ctx, in)
	}

	info := &server.InterceptorServerInfo{
		Server: srv,
		Method: "SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KVServer).SayHello(ctx, req.(*HelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

var KV_ServiceDesc = &server.ServiceDesc{
	Methods: []server.MethodDesc{
		{
			Name:    "SayHello",
			Handler: _KVServer_SayHello_Handler,
		},
	},
}