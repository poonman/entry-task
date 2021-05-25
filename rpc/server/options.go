package server

import (
	"crypto/tls"
	"time"
)

type Options struct {


	tlsConfig *tls.Config
	readTimeout time.Duration
	writeTimeout time.Duration
}

type Option func(options *Options)
