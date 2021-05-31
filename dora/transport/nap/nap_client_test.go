package nap

import (
	"context"
	"github.com/poonman/entry-task/dora/transport"
	"sync"
	"testing"
	"time"
)


func  TestClientSendAndRecv(t *testing.T) {
	server, err := NewListener(&ListenerOptions{
		Address:   "0.0.0.0:9000",
		TlsConfig: nil,
	})
	if err != nil {
		t.Fatalf("failed to listen. err:[%v]", err)
	}
	ctx, cancel := context.WithCancel(context.TODO())

	go func() {
		timer := time.NewTimer(3*time.Second)
		<-timer.C

		cancel()
	}()
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		err = server.Accept(ctx, func(tr transport.ServerTransport) {
			err = tr.Serve(func(stream transport.ServerStream) error {

				t.Logf("recv req:[%s]", string(stream.GetPayload()))

				err = tr.Send(stream, []byte("pong"))
				if err != nil {
					t.Fatalf("failed to send. err:[%v]", err)
				}
				return nil
			})
			if err != nil {
				t.Fatalf("failed to serve. err:[%v]", err)
			}

		})
		if err != nil {
			t.Fatalf("failed to accept. err:[%v]", err)
		}

		wg.Done()
	}()

	dialer := NewDialer(&DialerOptions{
		DialTimeout: 5*time.Second,
		TlsConfig:   nil,
	})

	client, err := dialer.Dial(ctx, "127.0.0.1:9000")
	if err != nil {
		t.Fatalf("failed to dial. err:[%v]", err)
	}

	cs := client.NewStream(ctx)
	err = cs.Send([]byte("ping"))
	if err != nil {
		t.Fatalf("failed to send. err:[%v]", err)
	}

	rsp, err := cs.Recv()
	if err != nil {
		t.Fatalf("failed to recv. err:[%v]", err)
	}

	t.Logf("recv rsp:[%s]", string(rsp))

	wg.Wait()

}