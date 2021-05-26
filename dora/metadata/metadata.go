package metadata

import (
	"context"
	"fmt"
)

type mdIncomingKey struct {}
type mdOutgoingKey struct {}

type MD map[string]string

func NewOutgoingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdOutgoingKey{}, md)
}

func NewIncomingContext(ctx context.Context, md MD) context.Context {
	return context.WithValue(ctx, mdIncomingKey{}, md)
}

func AppendToOutgoingContext(ctx context.Context, kv ...string) context.Context {
	if len(kv)%2 == 1 {
		panic(fmt.Sprintf("metadata: AppendToOutgoingContext got an odd number of input pairs for metadata: %d", len(kv)))
	}
	md, _ := ctx.Value(mdOutgoingKey{}).(MD)
	var k string
	for i, v := range kv {
		if i%2 == 0 {
			k = v
			continue
		}

		md[k] = v
	}

	return context.WithValue(ctx, mdOutgoingKey{}, md)
}

func FromOutgoingContext(ctx context.Context) (md MD, ok bool) {
	md, ok = ctx.Value(mdOutgoingKey{}).(MD)
	if !ok {
		return nil, false
	}

	return md, true
}


func FromIncomingContext(ctx context.Context) (md map[string]string, ok bool) {
	md, ok = ctx.Value(mdIncomingKey{}).(map[string]string)
	return
}

