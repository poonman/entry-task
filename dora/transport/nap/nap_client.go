package nap

import (
	"context"
	"errors"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/transport"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

var (
	streamID uint32
)

type clientTransport struct {
	mu sync.Mutex
	conn net.Conn

	ctx context.Context
	cancel context.CancelFunc

	lastReadAt int64

	pendingMu sync.Mutex
	pending map[uint32]*ClientStream

	keepaliveTime time.Duration
	keepaliveTimeout time.Duration

	errChan chan error

}

func (t *clientTransport) NewStream(ctx context.Context, cancel context.CancelFunc) transport.ClientStream {
	cs := &ClientStream{
		ctx: ctx,
		cancel: cancel,
		frame: &Frame{
			Header: FrameHeader{
				Type:   0,
				Flags:  0,
				Length: 0,
				ID:     atomic.AddUint32(&streamID, 1),
			},
		},
		In:       nil,
		Out:      nil,
		t:     t,
		dataChan: make(chan *Frame),
		//stopChan: make(chan struct{}),
		err:      nil,
	}

	t.pendingMu.Lock()
	defer t.pendingMu.Unlock()
	t.pending[cs.frame.Header.ID] = cs

	log.Debugf("[transport] New stream success. id:[%d]", cs.frame.Header.ID)
	//
	//go func() {
	//	<-ctx.Done()
	//	t.CloseStream(cs)
	//}()

	return cs
}

func (t *clientTransport) close(err error) error{
	t.cancel()

	for _, cs := range t.pending {
		cs.close()
	}

	if err != nil {
		go func() {
			t.errChan <- err
		}()
	}

	log.Debugf("[transport] Close client success.")
	return t.conn.Close()
}

func (t *clientTransport) Close() (err error) {
	return t.close(nil)
}

func (t *clientTransport) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *clientTransport) Local() string {
	return t.conn.LocalAddr().String()
}

func (t *clientTransport) Error() <-chan error {
	return t.errChan
}

func (t *clientTransport) keepalive() {
	timer := time.NewTimer(t.keepaliveTime)

	prevNano := time.Now().UnixNano()

	timeoutLeft := time.Duration(0)

	ping := false

	defer func() {
		timer.Stop()
	}()

	for {
		select {
		case <- timer.C:
			lastReadAt := atomic.LoadInt64(&t.lastReadAt)

			if lastReadAt > prevNano {
				ping = false
				timer.Reset(time.Duration(lastReadAt) + t.keepaliveTime - time.Duration(time.Now().UnixNano()))
				prevNano = lastReadAt
				continue
			}

			if ping && timeoutLeft <= 0 {
				_ = t.close(errors.New("keepalive ping failed to receive ACK within timeout"))
				return
			}

			if !ping {
				t.ping(false)
				timeoutLeft = t.keepaliveTimeout
				ping = true
			}

			sleepDuration := minTime(t.keepaliveTime, timeoutLeft)
			timeoutLeft = sleepDuration
			timer.Reset(sleepDuration)
		case <- t.ctx.Done():
			log.Debug("[keepalive] ctx.Done canceled")
			return
		}
	}
}

func (t *clientTransport) ping(ack bool) {
	if ack {
		pingFrame.Header.Flags = FlagPingAck
	} else {
		pingFrame.Header.Flags = 0
	}

	_ = t.send(pingFrame)
}

func (t *clientTransport) handlePing(frame *Frame) {
	//log.Debugf("[transport] Handle ping. ack:[%v]", frame.IsPingAck())

	if frame.IsPingAck() {

		return
	}

	t.ping( true)
}

func (t *clientTransport) handleData(frame *Frame) {
	log.Debugf("[transport] Handle data. id:[%d]", frame.Header.ID)

	cs := t.getStream(frame.Header.ID)
	if cs == nil {
		return
	}

	go cs.put(frame)
}

func (t *clientTransport) send(frame *Frame) error{
	t.mu.Lock()
	defer t.mu.Unlock()

	return send(t.conn, frame)
}

func (t *clientTransport) serve() {
	go t.keepalive()

	var err error
	for {
		select {
		case <-t.ctx.Done():
			log.Debug("[transport] ctx.Done canceled")

			return
		default:
			// do nothing
		}

		var in *Frame
		//log.Debugf("begin recv request...")
		in, err = recv(t.conn)
		if err != nil {
			if err == io.EOF {
				log.Warnf("[transport] Recv EOF")
			} else {
				log.Errorf("[transport] Failed to recv request. err:[%v]", err)
			}

			go func() {
				t.errChan <- err
			}()

			return
		}

		atomic.StoreInt64(&t.lastReadAt, time.Now().UnixNano())

		switch in.Header.Type {
		case FramePing:
			t.handlePing(in)

		case FrameDate:
			t.handleData(in)
		default:
			// wrong type
			_ = t.close(errors.New("invalid frame type"))
		}
	}
}

func (t *clientTransport) CloseStream(stream transport.ClientStream) {
	t.pendingMu.Lock()
	defer t.pendingMu.Unlock()

	cs := stream.(*ClientStream)
	delete(t.pending, cs.frame.Header.ID)

	cs.close()

	log.Debugf("[transport] Close stream success. id:[%d]", cs.frame.Header.ID)
}

func (t *clientTransport) getStream(id uint32) *ClientStream {
	t.pendingMu.Lock()
	defer t.pendingMu.Unlock()

	cs, ok := t.pending[id]
	if !ok {
		return nil
	}

	return cs
}

func NewClientTransport(ctx context.Context, conn net.Conn) transport.ClientTransport {
	ctx, cancel := context.WithCancel(ctx)

	t := &clientTransport{
		mu:               sync.Mutex{},
		conn:             conn,
		ctx:              ctx,
		cancel: cancel,
		lastReadAt:       0,
		pendingMu:        sync.Mutex{},
		pending:          make(map[uint32]*ClientStream),
		keepaliveTime:    5*time.Second,
		keepaliveTimeout: 20*time.Second,
		errChan:          make(chan error),
	}

	log.Debugf("[transport] New client transport success. local:[%s] <==> remote:[%s]",
		conn.LocalAddr().String(), conn.RemoteAddr().String())

	go t.serve()

	return t
}

