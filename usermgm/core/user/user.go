package user

import (
	"context"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
	"github.com/sirupsen/logrus"
)

type Svc struct {
	store  UserStore
	logger *logrus.Entry
}

func New(rs UserStore, logger *logrus.Entry) *Svc {
	return &Svc{
		store:  rs,
		logger: logger,
	}
}

type UserStore interface {
	CreateUser(ctx context.Context, user storage.User) (string, error)
	GetUserByEmail(ctx context.Context, email string) (*storage.User, error)
}
