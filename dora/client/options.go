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

type Invoker interface {
	Invoke(ctx context.Context, method string, in, out interface{}) error
}