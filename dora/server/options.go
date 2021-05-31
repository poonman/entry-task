package server

import (
	"context"
	"crypto/tls"
)

type Options struct {
	tlsConfig *tls.Config

	Interceptor Interceptor
}

type Option func(options *Options)

func WithTlsConfig(config *tls.Config) Option {
	return func(options *Options) {
		options.tlsConfig = config
	}
}

func WithInterceptor(i Interceptor) Option {
	return func(options *Options) {
		options.Interceptor = i
	}
}

type MethodHandler func(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor Interceptor) (rsp interface{}, err error)

type MethodDesc struct {
	Name    string
	Handler MethodHandler
}

func (m *MethodDesc) String() string {
	return m.Name
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
