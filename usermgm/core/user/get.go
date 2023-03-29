package user

import (
	"context"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
)

func (s *Svc) GetUserByEmail(ctx context.Context, email string) (*storage.User, error) {
	log := s.logger.WithField("method", "Core.User.GetUser")
	res, err := s.store.GetUserByEmail(ctx, email)
	if err != nil {
		log.WithError(err).Error("failed to get user")
		return nil, err
	}

	return res, nil
}
