package transport

import (
	"context"
)

// Dialer provide a base function Dial to connect a server
type Dialer interface {
	Dial(ctx context.Context, address string) (t ClientTransport, err error)
}

// Listener provide a base function Accept to accept connection from client
type Listener interface {
	Address() string
	Accept(ctx context.Context, accept func(t ServerTransport)) (err error)
	Close() error
}

// ClientStream is for one unary call, it make sure the same context when Send and Recv
type ClientStream interface {
	Send(in []byte)(err error)
	Recv() (out []byte, err error)
}

// ClientTransport is a client-end connection between client and server.
type ClientTransport interface {
	NewStream(ctx context.Context, cancelFunc context.CancelFunc) (stream ClientStream)
	Close() (err error)
	CloseStream(stream ClientStream)
	Remote() string
	Local() string
	Error() <-chan error
}

// ServerStream is for one unary call and include context
type ServerStream interface {
	GetPayload() []byte
}

// ServerTransport is a server-end connection between server and client
type ServerTransport interface {
	Serve(func(stream ServerStream) error) (err error)
	Send(stream ServerStream, out []byte) (err error)
	Close() (err error)
	Remote() string
	Local() string
}
