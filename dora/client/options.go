package client

import (
	"crypto/tls"
	"time"
)

type Options struct {
	tlsConfig *tls.Config

	connectTimeout time.Duration
}

type Option func(options *Options)