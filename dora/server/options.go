package server

import (
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
