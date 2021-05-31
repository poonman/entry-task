package nap

import (
	"context"
	"errors"
)

// ClientStream is for one call. its lifecycle is one-round trip between Send and Recv or
// based on the timeout provided by user context.
type ClientStream struct {
	ctx context.Context
	cancel context.CancelFunc
	frame *Frame
	In []byte
	Out []byte

	t *clientTransport

	dataChan chan *Frame

	//stopChan chan struct{}

	err error
}

func (cs *ClientStream) close() {
	cs.cancel()
}

func (cs *ClientStream) put(frame *Frame) {
	select {
	// this ctx has the same lifecycle with the user provided context.
	case <- cs.ctx.Done():
		return

	case cs.dataChan <- frame:

	}
}

func (cs *ClientStream) Send(in []byte) (err error) {
	length := uint32(len(in))
	if length == 0 {
		return
	}

	if length > MaxFrameSize {
		err = ErrFrameTooLarge
		return
	}

	cs.frame.Header.Type = FrameDate
	cs.frame.Header.Length = length
	cs.frame.Payload = in

	err = cs.t.send(cs.frame)
	if err != nil {
		return
	}

	return
}

func (cs *ClientStream) Recv() (out []byte, err error) {
	select {
	case <- cs.ctx.Done():
		err = errors.New("transport closed")
		return
	case frame := <- cs.dataChan:
		out = frame.Payload
		err = cs.err
	}
	return
}

type ServerStream struct {
	frame *Frame
}

func (ss *ServerStream) GetPayload() []byte{
	return ss.frame.Payload
}