package client

import (
	"context"
	"github.com/poonman/entry-task/dora/codec"
	"github.com/poonman/entry-task/dora/codec/proto"
	"github.com/poonman/entry-task/dora/metadata"
	"github.com/poonman/entry-task/dora/protocol"
	"github.com/poonman/entry-task/dora/status"
	"strconv"
)

type Stream struct {
	method   string
	cc       *Connection
	callInfo *callInfo

	ctx context.Context
}

type callInfo struct {
	serializeType string
}

func defaultCallInfo() *callInfo {
	return &callInfo{serializeType: proto.Name}
}

func NewStream(ctx context.Context, method string, cc *Connection) *Stream {
	s := &Stream{
		ctx:      ctx,
		method:   method,
		cc:       cc,
		callInfo: defaultCallInfo(),
	}

	return s
}

func (s *Stream) SendMsg(req interface{}) (err error) {
	//	log.Info("[dora] SendMsg begin...")

	msg, err := s.prepareMsg(req)
	if err != nil {
		return
	}

	return s.cc.sendRequest(msg)
}

func (s *Stream) prepareMsg(req interface{}) (msg *protocol.Message, err error) {
	c := codec.GetCodec(s.callInfo.serializeType)

	payload, err := c.Marshal(req)
	if err != nil {
		err = status.New(status.Internal, "marshal request error")
		return
	}

	msg = &protocol.Message{
		PkgHead: &protocol.PkgHead{
			Head: &protocol.Head{
				Version:       0,
				MessageType:   protocol.Head_Request,
				SerializeType: s.callInfo.serializeType,
				Seq:           s.cc.seq,
			},
			Method: s.method,
			Meta:   nil,
		},
		Payload: payload,
	}

	// extract context
	md, ok := metadata.FromOutgoingContext(s.ctx)
	if ok {
		msg.PkgHead.Meta = md
	}

	//	log.Debugf("md:[%v]", md)

	return
}

func (s *Stream) RecvMsg(rsp interface{}) (err error) {
	msg, err := s.cc.recvResponse()
	if err != nil {
		return
	}

	// invalid rsp seq
	if msg.PkgHead.Head.Seq != s.cc.seq {
		return
	}

	// parse dora status
	meta := msg.PkgHead.Meta
	err = parseError(meta)
	if err != nil {
		return
	}

	c := codec.GetCodec(s.callInfo.serializeType)

	err = c.Unmarshal(msg.Payload, rsp)
	if err != nil {
		//err = ErrInternal
		err = status.New(status.Internal, "unmarshal payload error")
		return
	}

	return
}

func parseError(meta map[string]string) (err error) {
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
