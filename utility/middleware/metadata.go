package middleware

import (
	"context"
	"fmt"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
)

type RequestMeta interface {
	Metadata(context.Context) (context.Context, error)
}

type PublicEndpoint interface {
	PublicEndpoint(method string) bool
}

type Service struct {
	logger   logrus.FieldLogger
	metaFunc []RequestMeta
}

func NewMetadata(logger logrus.FieldLogger, metaLoaders ...RequestMeta) (*Service, error) {
	s := &Service{
		logger:   logger,
		metaFunc: make([]RequestMeta, 0, len(metaLoaders)),
	}
	for _, m := range metaLoaders {
		if m != nil {
			s.metaFunc = append(s.metaFunc, m)
		}
	}
	if len(s.metaFunc) == 0 {
		return nil, fmt.Errorf("no metadata loaders")
	}
	return s, nil
}

func noPub(string) bool { return false }

func (s *Service) UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		logger := s.logger.WithField("full_method", info.FullMethod)
		svcPub := noPub
		if p, ok := info.Server.(PublicEndpoint); ok {
			svcPub = p.PublicEndpoint
		}
		ctx = fixKeys(ctx)
		var finalCtx context.Context
		var err error
		for _, f := range s.metaFunc {
			// check public endpoints of metaloader
			metaPub := noPub
			if p, ok := f.(PublicEndpoint); ok {
				metaPub = p.PublicEndpoint
			}
			// exectue metaloader
			finalCtx, err = f.Metadata(ctx)
			if err != nil {
				logger.WithError(err).Debug("metadata failed")
				if !metaPub(info.FullMethod) && !svcPub(info.FullMethod) {
					// Private endpoints error on first failure
					return nil, status.Error(codes.PermissionDenied, "permission denied")
				}
				finalCtx = ctx
				continue
			}
			// Update context to allow chaining of metaloaders
			ctx = finalCtx
		}
		return handler(finalCtx, req)
	}
}

func fixKeys(ctx context.Context) context.Context {
	md := metadata.MD{}
	for k, v := range metautils.ExtractIncoming(ctx) {
		md.Append(k, v...)
	}
	return metautils.NiceMD(md).ToIncoming(ctx)
}
