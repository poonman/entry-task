package server

import (
	"crypto/tls"
)

type Options struct {


	tlsConfig *tls.Config

	Interceptor Interceptor
	//readTimeout time.Duration
	//writeTimeout time.Duration
}

type Option func(options *Options)
