package client

import (
	"bufio"
	"github.com/poonman/entry-task/rpc/log"
	"github.com/poonman/entry-task/rpc/protocol"
	"net"
)

type Connection struct {
	conn net.Conn
	
	opts *Options

	seq uint64
}

func NewConnection(address string, opts ...Option) (c *Connection, err error) {
	c = &Connection{
		conn: nil,
		opts: &Options{},
	}

	for _, o := range opts {
		o(c.opts)
	}

	if c.opts.tlsConfig != nil {
		//dialer := &net.Dialer{
		//	Timeout:       c.opts.connectTimeout,
		//}
	} else {
		var conn net.Conn
		conn, err = net.DialTimeout("tcp", address, c.opts.connectTimeout)
		if err != nil {
			log.Errorf("Failed to dial. err:[%v]", err)
			return
		}

		c.conn  = conn
	}

	return
}

func (c *Connection) sendRequest(msg *protocol.Message) (err error) {
	w := bufio.NewWriter(c.conn)

	return protocol.WriteMessage(w, msg)
}

func (c *Connection) recvResponse() (msg *protocol.Message, err error){
	r := bufio.NewReader(c.conn)

	return protocol.ReadMessage(r)
}