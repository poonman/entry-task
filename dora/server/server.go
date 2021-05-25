package server

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"github.com/poonman/entry-task/dora/codec"
	"github.com/poonman/entry-task/dora/codec/proto"
	"github.com/poonman/entry-task/dora/log"
	"github.com/poonman/entry-task/dora/protocol"
	"github.com/poonman/entry-task/dora/status"
	"io"
	"net"
	"runtime"
	"strconv"
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

	ErrServerClosed = errors.New("server closed")
)


type Server struct {
	mu sync.RWMutex
	ln net.Listener

	activeConnMap map[net.Conn]struct{}
	//seq           uint64

	serverImpl interface{}
	methods map[string]*MethodDesc

	opts *Options

	stopCh chan struct{}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		mu:            sync.RWMutex{},
		ln:            nil,
		activeConnMap: make(map[net.Conn]struct{}),
		serverImpl:    nil,
		methods:       nil,
		opts:          &Options{},
	}

	for _, o := range opts {
		o(s.opts)
	}

	return s
}

func (s *Server) Serve(address string) (err error) {
	var (
		ln net.Listener
		tempDelay time.Duration
	)

	if s.opts.tlsConfig != nil {
		ln, err = tls.Listen("tcp", address, s.opts.tlsConfig)
	} else {
		ln, err = net.Listen("tcp", address)
	}
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
			// check if server is stopped
			select {
			case <-s.stopCh:
				return ErrServerClosed
			default:
				// do nothing
			}

			if ne, ok := err1.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				log.Errorf("Failed to accept. err:[%v], tempDelay:[%v]", err1, tempDelay)
				time.Sleep(tempDelay)
				continue
			}

			return err1
		}

		tempDelay = 0

		s.mu.Lock()
		s.activeConnMap[conn] = struct{}{}
		s.mu.Unlock()

		go s.serveConn(conn)
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
		delete(s.activeConnMap, conn)
		s.mu.Unlock()
		conn.Close()

	}()

	// todo: check if server shutdown?

	if tlsConn, ok := conn.(*tls.Conn); ok {

		if err := tlsConn.Handshake(); err != nil {
			log.Errorf("Failed to handshake tls conn. remoteAddr:[%s], err:[%v]", conn.RemoteAddr(), err)
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

func (s *Server) handleRequest(ctx context.Context, reqMsg *protocol.Message) (rspMsg *protocol.Message, err error) {
	method := reqMsg.PkgHead.Method

	rspMsg = reqMsg.Clone()

	rspMsg.SetMessageType(protocol.Head_Response)

	s.mu.RLock()
	md, ok := s.methods[method]
	s.mu.RUnlock()

	if !ok {
		err = status.Newf(status.Unimplemented, "method '%s' is unimplemented", method)
		handleError(rspMsg, err)
		return
	}

	c := s.getCodec(reqMsg.PkgHead.Head.SerializeType)

	df := func(v interface{}) error {
		err := c.Unmarshal(reqMsg.Payload, v)
		if err != nil {
			return ErrInternal
		}

		return nil
	}

	var (
		rsp interface{}
		payload []byte
	)

	rsp, err = md.Handler(s.serverImpl, ctx, df, s.opts.Interceptor)
	if err != nil {
		handleError(rspMsg, err)
		return
	}

	payload, err = c.Marshal(rsp)
	if err != nil {
		err = ErrInternal
		handleError(rspMsg, err)
		return
	}

	rspMsg.Payload = payload

	return
}

func handleError(m *protocol.Message, err error) {
	st, ok := err.(*status.Status)
	if !ok {
		m.PkgHead.Meta["dora-status"] = strconv.Itoa(int(status.Unknown))
		m.PkgHead.Meta["dora-message"] = status.Unknown.String()
		return
	}

	m.PkgHead.Meta["dora-status"] = strconv.Itoa(int(st.Code))
	m.PkgHead.Meta["dora-message"] = st.Message
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


type ServiceRegistrar interface {
	RegisterService(sd *ServiceDesc, impl interface{})
}
//==============

func (s *Server) RegisterService(sd *ServiceDesc, impl interface{}) {
	for k := range sd.Methods {
		m := &sd.Methods[k]
		s.methods[m.Name] = m
	}

	s.serverImpl = impl
}