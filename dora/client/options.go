package client

import (
	"context"
	"crypto/tls"
	"time"
)

type Options struct {
	tlsConfig *tls.Config

	connectTimeout time.Duration

	connSize int
}

type Option func(options *Options)

func WithTlsConfig(config *tls.Config) Option {
	return func(options *Options) {
		options.tlsConfig = config
	}
}

func WithConnSize(size int) Option {
	return func(options *Options) {
		options.connSize = size
	}
}

type Invoker interface {
	Invoke(ctx context.Context, method string, in, out interface{}) error
}
