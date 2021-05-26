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

type StoreClient interface {
	WriteSecureMessage(ctx context.Context, in *WriteSecureMessageReq) (out *WriteSecureMessageRsp, err error)
	ReadSecureMessage(ctx context.Context, in *ReadSecureMessageReq) (out *ReadSecureMessageRsp, err error)
}

type storeClient struct {
	cc client.Invoker
}

func (c *storeClient) WriteSecureMessage(ctx context.Context, in *WriteSecureMessageReq) (out *WriteSecureMessageRsp, err error) {
	out = &WriteSecureMessageRsp{}

	err = c.cc.Invoke(ctx, "WriteSecureMessage", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}
func (c *storeClient) ReadSecureMessage(ctx context.Context, in *ReadSecureMessageReq) (out *ReadSecureMessageRsp, err error) {
	out = &ReadSecureMessageRsp{}

	err = c.cc.Invoke(ctx, "ReadSecureMessage", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type StoreServer interface {
	WriteSecureMessage(ctx context.Context, in *WriteSecureMessageReq, out *WriteSecureMessageRsp) (err error)
	ReadSecureMessage(ctx context.Context, in *ReadSecureMessageReq, out *ReadSecureMessageRsp) (err error)
	mustEmbedUnimplementedStoreServer()
}

type UnimplementedStoreServer struct {
}

func (UnimplementedStoreServer) WriteSecureMessage(context.Context, *WriteSecureMessageReq, *WriteSecureMessageRsp) error {
	return status.New(status.Unimplemented, "server is unimplemented")
}

func (UnimplementedStoreServer) ReadSecureMessage(context.Context, *ReadSecureMessageReq, *ReadSecureMessageRsp) error {
	return status.New(status.Unimplemented, "server is unimplemented")
}

func (UnimplementedStoreServer) mustEmbedUnimplementedStoreServer() {}

func RegisterStoreServer(r server.ServiceRegistrar, impl StoreServer) {
	r.RegisterService(Store_ServiceDesc, impl)
}
func _StoreServer_WriteSecureMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor server.Interceptor) (_ interface{}, err error) {
	in := new(WriteSecureMessageReq)
	if err = dec(in); err != nil {
		return nil, err
	}

	out := new(WriteSecureMessageRsp)

	if interceptor == nil {
		err = srv.(StoreServer).WriteSecureMessage(ctx, in, out)
		return out, err
	}

	info := &server.InterceptorServerInfo{
		Server: srv,
		Method: "WriteSecureMessage",
	}
	handler := func(ctx context.Context, in, out interface{}) error {
		return srv.(StoreServer).WriteSecureMessage(ctx, in.(*WriteSecureMessageReq), out.(*WriteSecureMessageRsp))
	}
	err = interceptor(ctx, in, out, info, handler)
	return out, err
}

func _StoreServer_ReadSecureMessage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor server.Interceptor) (_ interface{}, err error) {
	in := new(ReadSecureMessageReq)
	if err = dec(in); err != nil {
		return nil, err
	}

	out := new(ReadSecureMessageRsp)

	if interceptor == nil {
		err = srv.(StoreServer).ReadSecureMessage(ctx, in, out)
		return out, err
	}

	info := &server.InterceptorServerInfo{
		Server: srv,
		Method: "ReadSecureMessage",
	}
	handler := func(ctx context.Context, in, out interface{}) error {
		return srv.(StoreServer).ReadSecureMessage(ctx, in.(*ReadSecureMessageReq), out.(*ReadSecureMessageRsp))
	}
	err = interceptor(ctx, in, out, info, handler)
	return out, err
}

var Store_ServiceDesc = &server.ServiceDesc{
	Methods: []server.MethodDesc{
		{
			Name:    "WriteSecureMessage",
			Handler: _StoreServer_WriteSecureMessage_Handler,
		},
		{
			Name:    "ReadSecureMessage",
			Handler: _StoreServer_ReadSecureMessage_Handler,
		},
	},
}
