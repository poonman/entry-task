package client

import (
	"context"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/dora/transport"
	"github.com/poonman/entry-task/dora/transport/nap"
	"io"
	"sync"
	"time"
)

type Client struct {
	// ctx is a parent context that manage all transport lifecycle
	ctx context.Context
	cancel context.CancelFunc

	mu         sync.RWMutex
	dialer     transport.Dialer

	// transport length is usually relative to the concurrency.
	// Another way to implement transport pool is to use slice and load balanced by robin instead
	transports map[transport.ClientTransport]struct{}

	address string

	opts *Options
}

func NewClient(address string, opts ...Option) (c *Client) {

	ctx, cancel := context.WithCancel(context.TODO())
	c = &Client{
		ctx:        ctx,
		cancel:     cancel,
		mu:         sync.RWMutex{},
		dialer:     nil,
		transports: make(map[transport.ClientTransport]struct{}),
		address:    address,
		opts:       &Options{},
	}

	for _, o := range opts {
		o(c.opts)
	}

	c.dialer = nap.NewDialer(&nap.DialerOptions{
		DialTimeout: c.opts.dialTimeout,
		TlsConfig:   c.opts.tlsConfig,
	})

	log.Infof("[dora] New client success. address:[%s]", address)

	return c
}

func (c *Client) Invoke(ctx context.Context, method string, in, out interface{}) (err error) {
	var (
		t transport.ClientTransport
	)

	t, err = c.getTransport()
	if err != nil {
		return
	}

	defer func() {
		if err == io.EOF {
			log.Errorf("[dora] Server has closed transport. err:[%v]", err)
			err = status.ErrUnavailable
		} else {
			c.releaseTransport(t)
		}
	}()

	var cancel context.CancelFunc
	// timeout to cancel this call
	ctx, cancel = context.WithTimeout(ctx, 5*time.Second)

	call := NewCall(ctx, cancel, method, t)

	err = call.SendMsg(in)
	if err != nil {
		return
	}

	return call.RecvMsg(out)
}

func (c *Client) getTransport() (t transport.ClientTransport, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.transports) == 0 {
		t, err = c.newTransport()
		if err != nil {
			return
		}

		return
	}

	for t_ := range c.transports {
		t = t_
		break
	}

	if t == nil {
		err = status.New(status.Unavailable, "connection is nil")
		return
	}

	delete(c.transports, t)

	return
}

func (c *Client) deleteTransport(t transport.ClientTransport) {
	c.mu.Lock()
	defer c.mu.Unlock()

	delete(c.transports, t)
}

func (c *Client) releaseTransport(t transport.ClientTransport) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.transports) > c.opts.connSize {
		// todo: conn.Close()
		// do not add to transports
		return
	}

	c.transports[t] = struct{}{}

}

func (c *Client) newTransport() (t transport.ClientTransport, err error) {

	t, err = c.dialer.Dial(c.ctx, c.address)
	if err != nil {
		return
	}

	go func() {
		select {
		case <- c.ctx.Done():
			return
		case err := <- t.Error():
			c.deleteTransport(t)
			log.Infof("[dora] Client transport closed. err:[%v]", err)
		}
	}()

	return
}

func (c *Client) Stop() {
	c.cancel()
}