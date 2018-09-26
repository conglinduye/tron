package grpcclient

import (
	"time"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

// default setting for grpc call
var (
	defaultTimeout = 120 * time.Second
)

func getTimeoutContext(timeout time.Duration) (context.Context, func()) {
	return context.WithTimeout(context.Background(), timeout)
}

func getDefaultCallOptions() grpc.EmptyCallOption {
	return grpc.EmptyCallOption{}
}
