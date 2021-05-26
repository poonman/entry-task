package client

import (
	"context"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/status"
	"sync"
	"sync/atomic"
)

type Client struct {
	mu      sync.RWMutex
	connMap map[uint64]*Connection

	connId uint64

	address string

	opts *Options
}

func NewClient(address string, opts ...Option) (c *Client) {
	log.Info("[dora] NewClient begin...")

	c = &Client{
		mu:      sync.RWMutex{},
		connMap: make(map[uint64]*Connection),
		address: address,
		opts:    &Options{},
	}

	for _, o := range opts {
		o(c.opts)
	}

	log.Infof("[dora] NewClient success. address:[%s]", address)

	return c
}

func (c *Client) Invoke(ctx context.Context, method string, in, out interface{}) (err error) {
	var (
		conn *Connection
	)

	conn, err = c.getConn()
	if err != nil {
		return
	}

	defer c.releaseConn(conn)

	cs := NewStream(ctx, method, conn)

	err = cs.SendMsg(in)
	if err != nil {

		return
	}

	return cs.RecvMsg(out)
}

func (c *Client) getConn() (conn *Connection, err error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.connMap) == 0 {
		conn, err = c.newConn()
		if err != nil {
			return
		}

		return
	}

	for _, v := range c.connMap {
		conn = v
		break
	}

	if conn == nil {
		err = status.New(status.Unavailable, "connection is nil")
		return
	}

	delete(c.connMap, conn.id)

	return
}

func (c *Client) releaseConn(conn *Connection) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if len(c.connMap) > c.opts.connSize {
		// todo: conn.Close()
		// do not add to connMap
		return
	}

	c.connMap[conn.id] = conn
}

func (c *Client) newConn() (conn *Connection, err error) {

	conn, err = NewConnection(c.address, c.opts)
	if err != nil {
		return
	}

	conn.id = atomic.AddUint64(&c.connId, 1)

	// todo: heartbeat

	return
}
