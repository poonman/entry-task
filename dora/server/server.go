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
	ErrInternal          = errors.New("internal error")

	ErrServerClosed = errors.New("server closed")
)

type Server struct {
	mu sync.RWMutex
	ln net.Listener

	activeConnMap map[net.Conn]struct{}
	//seq           uint64

	serverImpl interface{}
	methods    map[string]*MethodDesc

	opts *Options

	stopCh chan struct{}
}

func NewServer(opts ...Option) *Server {
	s := &Server{
		mu:            sync.RWMutex{},
		ln:            nil,
		activeConnMap: make(map[net.Conn]struct{}),
		serverImpl:    nil,
		methods:       make(map[string]*MethodDesc),
		opts:          &Options{},
		stopCh:        make(chan struct{}),
	}

	for _, o := range opts {
		o(s.opts)
	}

	log.Info("NewServer success...")

	return s
}

func (s *Server) Serve(address string) (err error) {
	var (
		ln        net.Listener
		tempDelay time.Duration
	)

	log.Info("[dora] Serve begin...")

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

	log.Info("[dora] Listen success...")

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
	log.Infof("[dora] serverConn begin... remote:[%s]", conn.RemoteAddr().String())

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
	//w := bufio.NewWriter(conn)

	for {
		// check if server is stopped
		select {
		case <-s.stopCh:
			return
		default:
			// do nothing
		}

		ctx := context.Background()

		log.Debugf("begin recv request...")

		req, err := s.recvRequest(r)
		if err != nil {
			if err == io.EOF {
				// close connect
				//_ = conn.Close()
				// close on defer, just log here
				log.Warnf("recv request error. ")
			} else {
				log.Errorf("Failed to recv request. err:[%v]", err)
			}

			return
		}

		rsp, err := s.handleRequest(ctx, req)
		if err != nil {
			// just log. rsp include err
			log.Errorf("Failed to handle request. err:[%v]", err)
		}

		err = s.sendResponse(rsp, conn)
		if err != nil {
			log.Errorf("Failed to send response. err:[%v]", err)
		}

	}
}

func (s *Server) recvRequest(r io.Reader) (msg *protocol.Message, err error) {
	msg, err = protocol.ReadMessage(r)

	log.Infof("[dora] recvRequest. msg:[%+v]", msg)

	return
}

func (s *Server) sendResponse(rsp *protocol.Message, w net.Conn /*w io.Writer*/) (err error) {
	log.Infof("[dora] sendResponse begin. rsp:[%+v]", rsp)

	err = protocol.WriteMessage(w, rsp)
	if err != nil {
		return err
	}

	log.Infof("[dora] sendResponse success.")
	return
}

func (s *Server) handleRequest(ctx context.Context, reqMsg *protocol.Message) (rspMsg *protocol.Message, err error) {

	log.Infof("[dora] handleRequest begin. ")

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
		rsp     interface{}
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

	log.Infof("handleRequest success.")

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
	log.Debugf("RegisterService begin...")
	for k := range sd.Methods {
		m := &sd.Methods[k]
		s.methods[m.Name] = m
	}

	s.serverImpl = impl

	log.Debug("RegisterService success")
}

func (s *Server) Stop() (err error){
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ln != nil {
		err = s.ln.Close()
	}
	for conn := range s.activeConnMap {
		conn.Close()
		delete(s.activeConnMap, conn)
	}
	s.closeDoneChanLocked()
	return err
}

func (s *Server) closeDoneChanLocked() {
	select {
	case <-s.stopCh:
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		// by s.mu.RegisterName
		close(s.stopCh)
	}
}