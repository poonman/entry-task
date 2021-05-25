package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"github.com/poonman/entry-task/rpc/codec"
	"github.com/poonman/entry-task/rpc/codec/proto"
	"github.com/poonman/entry-task/rpc/log"
	"github.com/poonman/entry-task/rpc/protocol"
	"io"
	"net"
	"runtime"
	"sync"
	"time"
)

const (
	// ReaderBufferSize is used for bufio reader.
	ReaderBufferSize = 1024
	// WriterBufferSize is used for bufio writer.
	WriterBufferSize = 1024
)

var (
	ErrUnsupportedMethod = errors.New("unsupported method")
	ErrInternal = errors.New("internal error")
)


type Server struct {
	mu sync.RWMutex
	ln net.Listener

	activeConnMgr map[net.Conn]struct{}
	seq uint64

	tlsConfig *tls.Config

	serverImpl interface{}
	methods map[string]*MethodDesc

	opts *Options
}

func NewServer(Opts ...Option) *Server {
	s := &Server{
		// TODO
	}




	return s
}

func (s *Server) Serve(address string) (err error) {
	var (
		ln net.Listener
		//tempDelay time.Duration
	)

	ln, err = tls.Listen("tcp", address, s.opts.tlsConfig)
	if err != nil {
		log.Errorf("Failed to listen. err:[%v]", err)
		return
	}

	s.mu.Lock()
	s.ln = ln
	s.mu.Unlock()

	for {
		conn, err1 := ln.Accept()
		if err1 != nil {
			// todo:

			return err1
		}

		//tempDelay = 0

		s.mu.Lock()
		s.activeConnMgr[conn] = struct{}{}
		s.mu.Unlock()


	}
}

func (s *Server) serveConn(conn net.Conn) {
	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Errorf("serving %s panic error: %s, stack:\n %s", conn.RemoteAddr(), err, buf)
		}

		s.mu.Lock()
		delete(s.activeConnMgr, conn)
		s.mu.Unlock()
		conn.Close()

	}()

	// todo: shutdown?

	if tlsConn, ok := conn.(*tls.Conn); ok {
		if d := s.opts.readTimeout; d != 0 {
			conn.SetReadDeadline(time.Now().Add(d))
		}
		if d := s.opts.writeTimeout; d != 0 {
			conn.SetWriteDeadline(time.Now().Add(d))
		}
		if err := tlsConn.Handshake(); err != nil {
			log.Errorf("rpcx: TLS handshake error from %s: %v", conn.RemoteAddr(), err)
			return
		}
	}

	r := bufio.NewReaderSize(conn, ReaderBufferSize)
	w := bufio.NewWriter(conn)

	for {
		// todo: readTimeout?

		ctx := context.Background()

		req, err := s.recvRequest(r)
		if err != nil {
			return
		}

		// todo: auth?

		rsp, err := s.handleRequest(ctx, req)
		if err != nil {
			// just log. rsp include err
			log.Errorf("Failed to handle request. err:[%v]", err)
		}

		err = s.sendResponse(rsp, w)
		if err != nil {
			log.Errorf("Failed to send response. err:[%v]", err)
		}

	}
}

func (s *Server) recvRequest(r io.Reader) (msg *protocol.Message, err error) {
	msg, err = protocol.ReadMessage(r)
	// unmask io.EOF
	if err == io.EOF {
		return msg, nil
	}

	return
}

func (s *Server) sendResponse(rsp *protocol.Message, w io.Writer) (err error) {
	err = protocol.WriteMessage(w, rsp)
	if err != nil {
		return err
	}

	return
}

func (s *Server) handleRequest(ctx context.Context, req *protocol.Message) (rsp *protocol.Message, err error) {
	method := req.PkgHead.Method

	rsp = req.Clone()

	rsp.SetMessageType(protocol.Head_Response)

	s.mu.RLock()
	md, ok := s.methods[method]
	s.mu.RUnlock()

	if !ok {
		err = ErrUnsupportedMethod
		return
	}

	c := s.getCodec(req.PkgHead.Head.SerializeType)

	df := func(v interface{}) error {
		err := c.Unmarshal(req.Payload, v)
		if err != nil {
			return ErrInternal
		}

		return nil
	}

	resp, err := md.Handler(s.serverImpl, ctx, df)
	if err != nil {
		handleError(rsp, err)
		// todo:
		return
	}

	payload, err := c.Marshal(resp)
	if err != nil {
		err = ErrInternal
		handleError(rsp, err)
		return
	}

	rsp.Payload = payload

	return
}

func handleError(m *protocol.Message, err error) {

}

func (s *Server) getCodec(typ string) codec.Codec {
	if typ == "" {
		return codec.GetCodec(proto.Name)
	}

	c := codec.GetCodec(typ)
	if c == nil {
		return codec.GetCodec(proto.Name)
	}

	return c
}