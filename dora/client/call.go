package client

import (
	"context"
	"github.com/poonman/entry-task/dora/codec"
	protoCodec "github.com/poonman/entry-task/dora/codec/proto"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/protocol"
	"github.com/poonman/entry-task/dora/status"
	"github.com/poonman/entry-task/dora/transport"
	"github.com/golang/protobuf/proto"
	"strconv"
	"sync/atomic"
)

var (
	atomicSeq uint64
)

// Call is a unary call and include context and other information
type Call struct {
	// method is request method name
	method   string
	// t is served for this call
	t       transport.ClientTransport
	callInfo *callInfo
	// cs is a synchronous stream which makes the call flow easy
	cs transport.ClientStream
	// seq is the current request sequence
	seq uint64
	ctx context.Context
}

type callInfo struct {
	serializeType string
}

func defaultCallInfo() *callInfo {
	return &callInfo{serializeType: protoCodec.Name}
}

func NewCall(ctx context.Context, cancel context.CancelFunc, method string, t transport.ClientTransport) *Call {
	s := &Call{
		ctx:      ctx,
		method:   method,
		t:       t,
		callInfo: defaultCallInfo(),
		cs: t.NewStream(ctx, cancel),
	}

	return s
}

func (s *Call) SendMsg(req interface{}) (err error) {

	msg, err := s.prepareMsg(req)
	if err != nil {
		return
	}

	payload, err := proto.Marshal(msg)
	if err != nil {
		err = status.New(status.Internal, err.Error())
		return
	}

	return s.cs.Send(payload)
}

func (s *Call) prepareMsg(req interface{}) (msg *protocol.Pkg, err error) {
	c := codec.GetCodec(s.callInfo.serializeType)

	payload, err := c.Marshal(req)
	if err != nil {
		err = status.New(status.Internal, "marshal request error")
		return
	}

	s.seq = atomic.AddUint64(&atomicSeq, 1)

	msg = &protocol.Pkg{
		Head: &protocol.Head{
			Version:       0,
			SerializeType: s.callInfo.serializeType,
			Seq: s.seq,
			Method: s.method,
			Meta:   nil,
		},
		Payload: payload,
	}

	// extract context
	md, ok := metadata.FromOutgoingContext(s.ctx)
	if ok {
		msg.Head.Meta = md
	}

	return
}

func (s *Call) RecvMsg(rsp interface{}) (err error) {

	var payload []byte

	defer func() {
		s.t.CloseStream(s.cs)
	}()

	// block here until recv an incoming response
	payload, err = s.cs.Recv()
	if err != nil {
		return
	}

	msg := &protocol.Pkg{}
	err = proto.Unmarshal(payload, msg)
	if err != nil {
		err = status.New(status.Internal, "unmarshal incoming pkg error")
		return
	}

	c := codec.GetCodec(s.callInfo.serializeType)
	err = c.Unmarshal(msg.Payload, rsp)
	if err != nil {
		err = status.New(status.Internal, "unmarshal rsp error")
		return
	}

	// parse dora status
	meta := msg.Head.Meta
	err = parseError(meta)
	if err != nil {
		return
	}

	return
}

func parseError(meta map[string]string) (err error) {
	if len(meta) == 0 {
		err = nil
		return
	}

	st, ok := meta["dora-status"]
	if !ok {
		err = nil
		return
	}

	code, err := strconv.ParseUint(st, 10, 32)
	if err != nil {
		err = status.New(status.Internal, "parse status error")
		return
	}

	if code == 0 {
		err = nil
		return
	}

	msg := meta["dora-message"]

	err = status.New(status.Code(code), msg)
	return
}
