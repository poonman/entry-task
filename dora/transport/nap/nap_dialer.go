package nap

import (
	"context"
	"crypto/tls"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/transport"
	"net"
	"time"
)

type dialer struct {
	opts *DialerOptions
}

type DialerOptions struct {
	DialTimeout time.Duration
	TlsConfig *tls.Config
}

func (d *dialer) Dial(ctx context.Context, address string) (t transport.ClientTransport, err error) {
	var conn net.Conn

	if d.opts.TlsConfig != nil {

		conn, err = tls.Dial("tcp", address, d.opts.TlsConfig)
		if err != nil {
			log.Errorf("[transport] Failed to dial address. err:[%v]", err)
			return
		}

	} else {
		conn, err = net.DialTimeout("tcp", address, d.opts.DialTimeout)
		if err != nil {
			log.Errorf("[transport] Failed to dial address. err:[%v]", err)
			return
		}

	}

	t = NewClientTransport(ctx, conn)

	return
}

func NewDialer(opts *DialerOptions) transport.Dialer {
	d := &dialer{
		opts: opts,
	}

	return d
}
