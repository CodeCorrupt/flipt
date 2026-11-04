package common

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

// ContextKey is a type for context keys used across Flipt packages.
type ContextKey string

const (
	HeaderFliptEnvironment            = "X-Flipt-Environment"
	HeaderFliptNamespace              = "X-Flipt-Namespace"
	HeaderFliptOFREPStream            = "x-flipt-ofrep-stream"
	ContextKeyOFREPStream  ContextKey = "ofrep-stream"
)

// WithOFREPStream sets the OFREP stream marker on the context.
func WithOFREPStream(ctx context.Context) context.Context {
	return context.WithValue(ctx, ContextKeyOFREPStream, true)
}

// IsOFREPStream reports whether the context has the OFREP stream marker set.
func IsOFREPStream(ctx context.Context) bool {
	v := ctx.Value(ContextKeyOFREPStream)
	ok, _ := v.(bool)
	return ok
}

// IsOFREPStreamRequest reports whether an incoming gRPC request is the OFREP
// event stream alias. The HTTP gateway marks those requests with the OFREP
// stream header and strips that header from client requests, so the marker is
// only trusted on in-process requests from the gateway.
func IsOFREPStreamRequest(ctx context.Context) bool {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return false
	}

	if vals := md.Get(HeaderFliptOFREPStream); len(vals) == 0 || vals[0] != "true" {
		return false
	}

	p, ok := peer.FromContext(ctx)
	return ok && p.Addr.Network() == "inproc"
}
