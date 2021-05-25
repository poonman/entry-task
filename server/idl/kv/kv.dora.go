package kv

import (
	"context"
	"github.com/poonman/entry-task/dora/client"
	"github.com/poonman/entry-task/dora/server"
	"github.com/poonman/entry-task/dora/status"
)

type HelloReq struct {

}

type HelloRsp struct {

}

type KvClient interface {
	SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error)
}

type kvClient struct {
	cc *client.Client
}

func (c *kvClient) SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error) {
	out = &HelloRsp{}

	err = c.cc.Invoke(ctx, "SayHello", in, out)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type KvServer interface {
	SayHello(ctx context.Context, in *HelloReq) (out *HelloRsp, err error)
}

// UnimplementedKvServer must be embedded to have forward compatible implementations.
type UnimplementedKvServer struct {
}

func (UnimplementedKvServer) SayHello(context.Context, *HelloReq) (*HelloRsp, error) {
	return nil, status.New(status.Unimplemented, "server is unimplemented")
}
func (UnimplementedKvServer) mustEmbedUnimplementedGreeterServer() {}

func RegisterKvServer(r server.ServiceRegistrar, impl KvServer) {
	r.RegisterService(Kv_ServiceDesc, impl)
}

//=========
func _KvServer_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor server.Interceptor) (interface{}, error) {
	in := new(HelloReq)
	if err := dec(in); err != nil {
		return nil, err
	}

	if interceptor == nil {
		return srv.(KvServer).SayHello(ctx, in)
	}


	info := &server.InterceptorServerInfo{
		Server:     srv,
		Method: "SayHello",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(KvServer).SayHello(ctx, req.(*HelloReq))
	}
	return interceptor(ctx, in, info, handler)
}

var Kv_ServiceDesc = &server.ServiceDesc{
	Methods: []server.MethodDesc{
		{
			Name:    "SayHello",
			Handler: _KvServer_SayHello_Handler,
		},
	}}