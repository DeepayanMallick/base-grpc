package user

import (
	"context"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	upb "github.com/DeepayanMallick/base-grpc/gunk/v1/usermgm/user"
)

func (h *Handler) CreateUser(ctx context.Context, req *upb.CreateUserRequest) (*upb.CreateUserResponse, error) {
	log := h.logger.WithField("method", "Service.User.CreateUser")

	dbUser := storage.User{
		FirstName: req.GetUser().FirstName,
		LastName:  req.GetUser().LastName,
		Username:  req.GetUser().Username,
		Email:     req.GetUser().Email,
		Password:  req.GetUser().Password,
		Status:    int32(req.User.GetStatus()),
	}

	res, err := h.usr.CreateUser(ctx, dbUser)
	if err != nil {
		errMsg := "failed to create user"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(codes.NotFound, "failed to create user")
	}

	return &upb.CreateUserResponse{
		ID: res,
	}, nil
}
