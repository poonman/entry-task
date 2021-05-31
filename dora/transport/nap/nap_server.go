package nap

import (
	"context"
	"github.com/poonman/entry-task/dora/misc/log"
	"github.com/poonman/entry-task/dora/transport"
	"io"
	"net"
	"sync"
	"sync/atomic"
	"time"
)

type serverTransport struct {
	ctx context.Context
	cancel context.CancelFunc

	mu sync.Mutex
	conn net.Conn
	lastReadAt int64

	// default 5s
	keepaliveTime time.Duration

	// default 20s
	keepaliveTimeout time.Duration
}

func (t *serverTransport) handlePing(frame *Frame) {
	if frame.IsPingAck() {
		return
	}

	t.ping( true)
}

func (t *serverTransport) send(frame *Frame) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	return send(t.conn, frame)
}

func (t *serverTransport) Serve(handle func(stream transport.ServerStream) error) (err error) {

	go t.keepalive()

	for {
		select {
		case <-t.ctx.Done():
			return
		default:
			// do nothing
		}

		var in *Frame
		//log.Debugf("begin recv request...")
		in, err = recv(t.conn)
		if err != nil {
			if err == io.EOF {
				log.Warnf("[transport[ Recv EOF. remote:[%s]", t.Remote())
			} else {
				log.Errorf("[transport] Failed to recv request. remote:[%s], err:[%v]", t.Remote(), err)
			}
			return
		}

		atomic.StoreInt64(&t.lastReadAt, time.Now().UnixNano())

		switch in.Header.Type {
		case FramePing:
			t.handlePing(in)

		case FrameDate:
			_ = handle(&ServerStream{
				frame: in,
			})
		default:
			_ = t.Close()
		}
	}
}

func (t *serverTransport) Send(stream transport.ServerStream, out []byte) (err error) {
	length := uint32(len(out))
	if length == 0 {
		return
	}

	if length > MaxFrameSize {
		err = ErrFrameTooLarge
		return
	}

	frame := stream.(*ServerStream).frame
	frame.Payload = out
	frame.Header.Length = length

	return t.send(frame)
}

func (t *serverTransport) Close() (err error) {
	t.cancel()

	return t.conn.Close()
}

func (t *serverTransport) Remote() string {
	return t.conn.RemoteAddr().String()
}

func (t *serverTransport) Local() string {
	return t.conn.LocalAddr().String()
}

var pingFrame = &Frame{
	Header:  FrameHeader{
		Type:   FramePing,
		Flags:  0,
		Length: 0,
		ID:     0,
	},
}
func (t *serverTransport) ping(ack bool) {
	if ack {
		pingFrame.Header.Flags = FlagPingAck
	} else {
		pingFrame.Header.Flags = 0
	}

	_ = t.send(pingFrame)
}

func (t *serverTransport) keepalive() {

	keepaliveTimer := time.NewTimer(10*time.Second)

	defer func() {
		keepaliveTimer.Stop()
	}()

	prevNano := time.Now().UnixNano()

	var keepaliveTimeoutLeft time.Duration

	ping := false

	for {
		select {
		case <- keepaliveTimer.C:
			lastReadAt := atomic.LoadInt64(&t.lastReadAt)

			if lastReadAt > prevNano {
				ping = false
				keepaliveTimer.Reset(time.Duration(lastReadAt) + t.keepaliveTime - time.Duration(time.Now().UnixNano()))
				prevNano = lastReadAt
				continue
			}

			if ping && keepaliveTimeoutLeft <= 0 {
				log.Infof("closing server transport due to idleness")

				_ = t.Close()
				return
			}

			if !ping {

				t.ping(false)
				keepaliveTimeoutLeft = t.keepaliveTimeout
				ping = true
			}

			sleepDuration := minTime(t.keepaliveTime, keepaliveTimeoutLeft)
			keepaliveTimeoutLeft = sleepDuration
			keepaliveTimer.Reset(sleepDuration)

		case <- t.ctx.Done():
			return
		}
	}
}

func NewServerTransport(ctx context.Context, conn net.Conn) transport.ServerTransport {
	ctx, cancel := context.WithCancel(ctx)
	t := &serverTransport{
		ctx:              ctx,
		cancel: cancel,
		mu:               sync.Mutex{},
		conn:             conn,
		lastReadAt:       0,
		keepaliveTime:    5 * time.Second,
		keepaliveTimeout: 20 * time.Second,
	}

	return t
}
