package server

import (
	"context"
	"github.com/poonman/entry-task/dora/codec"
	protoCodec "github.com/poonman/entry-task/dora/codec/proto"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/protocol"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/dora/transport"
	"github.com/poonman/entry-task/dora/transport/nap"
	"github.com/golang/protobuf/proto"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Server struct {
	mu              sync.RWMutex
	listener        transport.Listener
	activeTransport map[transport.ServerTransport]struct{}

	// server implement. it is a user side handler
	serverImpl interface{}
	// methods is a set of method descriptor of user side handler
	methods    map[string]*MethodDesc

	opts *Options

	// ctx is a parent ctx that manage all goroutines lifecycle
	ctx context.Context
	cancel context.CancelFunc

}

func NewServer(opts ...Option) *Server {
	s := &Server{
		mu:            sync.RWMutex{},

		listener:        nil,
		activeTransport: make(map[transport.ServerTransport]struct{}),
		serverImpl:      nil,
		methods:         make(map[string]*MethodDesc),
		opts:            &Options{},
	}

	for _, o := range opts {
		o(s.opts)
	}

	s.ctx, s.cancel = context.WithCancel(context.TODO())

	log.Info("[dora] New server success.")

	return s
}

func (s *Server) Serve(address string) (err error) {
	var (
		listener        transport.Listener
	)

	listener, err = nap.NewListener(&nap.ListenerOptions{
		Address:   address,
		TlsConfig: s.opts.tlsConfig,
	})
	if err != nil {
		log.Errorf("[dora] Failed to listen. address:[%v], err:[%v]", address, err)
		return
	}

	s.mu.Lock()
	s.listener = listener
	s.mu.Unlock()

	err = s.listener.Accept(s.ctx, s.serveTransport)
	if err != nil {
		log.Errorf("[dora] Failed to accept. listen:[%s], err:[%v]", s.listener.Address(), err)
		return
	}

	return
}

func (s *Server) serveTransport(t transport.ServerTransport) {
	log.Infof("[dora] Begin serve transport... remote:[%s]", t.Remote())

	defer func() {
		if err := recover(); err != nil {
			const size = 64 << 10
			buf := make([]byte, size)
			ss := runtime.Stack(buf, false)
			if ss > size {
				ss = size
			}
			buf = buf[:ss]
			log.Errorf("[dora] Serving %s panic error: %s, stack:\n %s", t.Remote(), err, buf)
		}

		s.mu.Lock()
		delete(s.activeTransport, t)
		s.mu.Unlock()
		t.Close()

	}()

	err := t.Serve(func(stream transport.ServerStream) (err error){
		// TODO: server side should implement a control buffer to send and recv
		ctx, _ := context.WithTimeout(context.TODO(), 5*time.Second)

		msg := &protocol.Pkg{}

		err = proto.Unmarshal(stream.GetPayload(), msg)
		if err != nil {
			log.Errorf("[dora] Failed to unmarshal incoming message. err:[%v]", err)
			err = status.New(status.Internal, err.Error())
			return
		}

		if msg.Head == nil {
			log.Errorf("[dora] A bad request. head is nil.")
			err = status.ErrBadRequest
			return
		}

		var rsp *protocol.Pkg
		var out []byte

		rsp  = s.handleRequest(ctx, msg)

		out, err = proto.Marshal(rsp)
		if err != nil {
			// should not fail. thus how to send response?
			// client-end need timeout to finish this request
			log.Errorf("[dora] Failed to marshal outgoing message. err:[%v]", err)
			err = status.New(status.Internal, err.Error())
			return
		}

		err = t.Send(stream, out)
		if err != nil {
			log.Errorf("[dora] Failed to send outgoing message. err:[%v]", err)
			return
		}
		return
	})
	if err != nil {
		log.Errorf("[dora] Failed to serve transport. err:[%v]", err)
	}
}

func (s *Server) handleRequest(ctx context.Context, in *protocol.Pkg) (out *protocol.Pkg) {

	methodName := in.Head.Method

	if len(in.Head.Meta) != 0 {
		ctx = metadata.NewIncomingContext(ctx, in.Head.Meta)
	}

	out = in.Clone()

	s.mu.RLock()
	method, ok := s.methods[methodName]
	s.mu.RUnlock()

	if !ok {
		handleError(out, status.Newf(status.Unimplemented, "method '%s' is unimplemented", methodName))
		return
	}

	c := s.getCodec(in.Head.SerializeType)

	df := func(v interface{}) error {
		err := c.Unmarshal(in.Payload, v)
		if err != nil {
			err = status.New(status.Internal, "unmarshal incoming payload error")
			return err
		}

		return nil
	}

	var (
		rsp     interface{}
		payload []byte
		err error
	)

	rsp, err = method.Handler(s.serverImpl, ctx, df, s.opts.Interceptor)
	if err != nil {
		handleError(out, err)
		return
	}

	payload, err = c.Marshal(rsp)
	if err != nil {
		err = status.New(status.Internal, "marshal outgoing payload error")
		handleError(out, err)
		return
	}

	out.Payload = payload

	return
}

func handleError(m *protocol.Pkg, err error) {
	st, ok := err.(*status.Status)
	if !ok {
		m.Head.Meta["dora-status"] = strconv.Itoa(int(status.Unknown))
		m.Head.Meta["dora-message"] = status.Unknown.String()
		return
	}

	m.Head.Meta["dora-status"] = strconv.Itoa(int(st.Code))
	m.Head.Meta["dora-message"] = st.Message
}

func (s *Server) getCodec(typ string) codec.Codec {
	if typ == "" {
		return codec.GetCodec(protoCodec.Name)
	}

	c := codec.GetCodec(typ)
	if c == nil {
		return codec.GetCodec(protoCodec.Name)
	}

	return c
}

type ServiceRegistrar interface {
	RegisterService(sd *ServiceDesc, impl interface{})
}

func (s *Server) RegisterService(sd *ServiceDesc, impl interface{}) {
	log.Info("[dora] Register service.")

	for k := range sd.Methods {
		m := &sd.Methods[k]
		s.methods[m.Name] = m
		log.Infof("[dora] Register method: %-20s", m)
	}

	s.serverImpl = impl
}

func (s *Server) Stop() (err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.listener != nil {
		err = s.listener.Close()
	}
	for socket := range s.activeTransport {
		socket.Close()
		delete(s.activeTransport, socket)
	}
	s.closeDoneChanLocked()

	log.Info("[dora] Stop server success")
	return err
}

func (s *Server) closeDoneChanLocked() {
	select {
	case <-s.ctx.Done():
		// Already closed. Don't close again.
	default:
		// Safe to close here. We're the only closer, guarded
		s.cancel()
	}
}
