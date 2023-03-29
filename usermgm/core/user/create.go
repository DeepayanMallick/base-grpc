package user

import (
	"context"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
)

func (s *Svc) CreateUser(ctx context.Context, user storage.User) (string, error) {
	log := s.logger.WithField("method", "Core.User.CreateUser")
	usr, err := s.store.CreateUser(ctx, user)
	if err != nil {
		log.WithError(err).Error("Failed to create user")
		return "", err
	}

	return usr, nil
}
