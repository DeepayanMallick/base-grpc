package middleware

import (
	"google.golang.org/grpc"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_logrus "github.com/grpc-ecosystem/go-grpc-middleware/logging/logrus"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"

	"github.com/sirupsen/logrus"
)

type Config struct {
	LogOpts []grpc_logrus.Option
}

// New creates a gRPC middleware chain.
func New(env string, logger *logrus.Entry, c Config, ints ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	rh := newRecoveryHandler(logger)
	interceptors := []grpc.UnaryServerInterceptor{
		grpc_ctxtags.UnaryServerInterceptor(grpc_ctxtags.WithFieldExtractor(grpc_ctxtags.CodeGenRequestFieldExtractor)),
		grpc_logrus.UnaryServerInterceptor(logger, c.LogOpts...),
		grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(rh.recover)),
	}
	if len(ints) > 0 {
		interceptors = append(interceptors, ints...)
	}

	return grpc_middleware.ChainUnaryServer(interceptors...)
}
