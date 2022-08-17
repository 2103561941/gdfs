package middleware

import (
	"github.com/cyb0225/gdfs/pkg/log"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	// add middlewares
	recoveryOpts = []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandler(func(p interface{}) error {
			err := status.Errorf(codes.Unknown, "panic triggered: %v", p)
			log.Error("recovery panic", log.Err(err))
			return err
		}),
	}
)

func UneryRecovery() grpc.UnaryServerInterceptor {
	return grpc_recovery.UnaryServerInterceptor(recoveryOpts...)
}

func StreamRecovery() grpc.StreamServerInterceptor {
	return grpc_recovery.StreamServerInterceptor(recoveryOpts...)
}
