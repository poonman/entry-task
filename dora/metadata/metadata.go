package metadata

import (
	"context"
)

type mdIncomingKey struct{}
type mdOutgoingKey struct{}

type MD map[string]string

func NewOutgoingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, md)
}

func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

func FromOutgoingContext(ctx context.Context) (md MD, ok bool) {
	md, ok = ctx.Value(mdOutgoingKey{}).(MD)
	if !ok {
		return nil, false
	}

	return md, true
}

func FromIncomingContext(ctx context.Context) (md MD, ok bool) {
	md, ok = ctx.Value(mdIncomingKey{}).(MD)
	return
}
