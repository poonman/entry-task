package nap

import (
	"context"
	"crypto/tls"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/dora/transport"
	"net"
	"time"
)

type listener struct {
	listener net.Listener
	opts *ListenerOptions
}

type ListenerOptions struct {
	Address string
	TlsConfig *tls.Config
}


func (l *listener) Close() error {
	return l.listener.Close()
}

func (l *listener) Accept(ctx context.Context, accept func(t transport.ServerTransport)) (err error) {
	var (
		tempDelay time.Duration
		conn net.Conn
	)

	for {
		conn, err = l.listener.Accept()
		if err != nil {
			select {
			case <- ctx.Done():
				return status.ErrInternal
			default:
				// do nothing
			}
			// check if server is stopped

			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Errorf("Failed to accept. err:[%v], tempDelay:[%v]", err, tempDelay)
				time.Sleep(tempDelay)
				return
			}

			return err
		}

		tempDelay = 0

		t := NewServerTransport(ctx, conn)

		go accept(t)
	}

}

func (l *listener) Address() string {
	return l.listener.Addr().String()
}

func NewListener(opts *ListenerOptions) (_ transport.Listener, err error) {
	l := &listener{
		listener: nil,
		opts: opts,
	}

	if l.opts.TlsConfig != nil {
		l.listener, err = tls.Listen("tcp", l.opts.Address, l.opts.TlsConfig)
	} else {
		l.listener, err = net.Listen("tcp", l.opts.Address)
	}
	if err != nil {
		log.Errorf("[transport] Failed to listen. err:[%v]", err)
		return
	}

	log.Infof("[transport] Listen success... address:[%s]", l.listener.Addr().String())

	return l, nil
}