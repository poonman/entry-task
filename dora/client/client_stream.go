package client

import (
	"context"
	"errors"
	"github.com/poonman/entry-task/dora/codec"
	"github.com/poonman/entry-task/dora/codec/proto"
	"github.com/poonman/entry-task/dora/protocol"
)

var (
	ErrInternal = errors.New("internal error")
)
type Stream struct {
	method string
	cc *Connection
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
		ctx :ctx,
		method:   method,
		cc:       cc,
		callInfo: defaultCallInfo(),
	}

	return s
}

func (s *Stream) SendMsg(req interface{}) (err error ){
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
		err = ErrInternal
		return
	}

	msg = &protocol.Message{
		PkgHead: &protocol.PkgHead{
			Head:                 &protocol.Head{
				Version:              0,
				MessageType:          protocol.Head_Request,
				SerializeType:        s.callInfo.serializeType,
				Seq:                  s.cc.seq,
			},
			Method:               s.method,
			Meta:                 nil, // todo: auth
		},
		Payload: payload,
	}

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

	c := codec.GetCodec(s.callInfo.serializeType)

	err = c.Unmarshal(msg.Payload, rsp)
	if err != nil {
		err = ErrInternal
		return
	}

	return
}