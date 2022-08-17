package middleware

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/cyb0225/gdfs/pkg/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

func UnaryServerInterceptor(removeAPI []string) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		start := time.Now()
		// get address.
		pr, ok := peer.FromContext(ctx)
		if !ok {
			return "", fmt.Errorf("missing metadata context")
		}
		addr := pr.Addr.String()

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			// log.Debugf("md[\"address\"]:%#v",md["address"])
			if v, ok := md["address"]; ok {
				addr = v[0]
				// log.Debug("get address", log.String("address", addr) )
			}
		}

		// log method.
		defer func() {
			// get required time.
			fullMethod := info.FullMethod
			patterns := strings.Split(fullMethod, "/")
			method := patterns[len(patterns) - 1]
			for i := 0; i < len(removeAPI); i++ {
				if method == removeAPI[i] {
					return
				}
			}
			durationMS := int64(time.Since(start) / time.Microsecond)
			// log.Debug("duration", log.String("start", start.Format("2006-01-02 15:04:05.000")), log.String("now", time.Now().Format("2006-01-02 15:04:05.000")))
			log.Info(method, log.String("address", addr), log.String("method", method), log.Int64("time(us)", durationMS))
		}()
		
		ctx = context.WithValue(ctx, "address", addr)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor(removeAPI []string) grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()

		// get address.
		pr, ok := peer.FromContext(ss.Context())
		if !ok {
			return fmt.Errorf("missing metadata context")
		}
		addr := pr.Addr.String()

		// log method.		
		defer func() {
			// get required time.
			fullMethod := info.FullMethod
			patterns := strings.Split(fullMethod, "/")
			method := patterns[len(patterns) - 1]
			for i := 0; i < len(removeAPI); i++ {
				if method == removeAPI[i] {
					return
				}
			}
			durationMS := int64(time.Since(start) / time.Microsecond)
			log.Info(method, log.String("address", addr), log.String("method", method), log.Int64("time(us)", durationMS))
		}()	

		return handler(srv, ss)
	}
}

// set address
// set address, this is because my server is running in wsl, 
// The address sent does not match the desired address
func UnaryClientInterceptor(address string) grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		ctx = metadata.AppendToOutgoingContext(ctx, "address", address)
		// log.Debug("put address", log.String("address", address))
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

// set address
func StreamClientInterceptor(address string) grpc.StreamClientInterceptor {
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		ctx = metadata.AppendToOutgoingContext(ctx, "address", address)
		return streamer(ctx, desc, cc, method, opts...)
	}
}
