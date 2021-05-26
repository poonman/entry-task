package server

import "context"

type MethodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor Interceptor) (rsp interface{}, err error)

type MethodDesc struct {
	Name    string
	Handler MethodHandler
}

type ServiceDesc struct {
	Methods []MethodDesc
}

type InterceptorServerInfo struct {
	Server interface{}
	Method string
}

type Handler func(ctx context.Context, in, out interface{}) (err error)

type Interceptor func(ctx context.Context, in, out interface{}, serverInfo *InterceptorServerInfo, handler Handler) (err error)
