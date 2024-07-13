package helper

import (
	"context"
	"google.golang.org/grpc/metadata"
)

// DumpIncomingContext is method to
func DumpIncomingContext(ctx context.Context) string {
	md, _ := metadata.FromIncomingContext(ctx)
	return Dump(&md)
}
