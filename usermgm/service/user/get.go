package user

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	upb "github.com/DeepayanMallick/base-grpc/gunk/v1/usermgm/user"
)

func (h *Handler) GetUserByEmail(ctx context.Context, req *upb.GetUserByEmailRequest) (*upb.GetUserByEmailResponse, error) {
	log := h.logger.WithField("method", "service.user.GetUser")
	r, err := h.usr.GetUserByEmail(ctx, req.GetEmail())
	if err != nil {
		errMsg := "failed to get user by id"
		log.WithError(err).Error(errMsg)
		return nil, status.Error(codes.NotFound, "failed to get user by id")
	}

	resp := &upb.GetUserByEmailResponse{
		User: &upb.User{
			ID:        r.ID,
			FirstName: r.FirstName,
			LastName:  r.LastName,
			Email:     r.Email,
			Username:  r.Username,
			Password:  r.Password,
			Status:    upb.Status(r.Status),
			CreatedAt: timestamppb.New(r.CreatedAt),
			UpdatedAt: timestamppb.New(r.UpdatedAt),
		},
	}

	return resp, nil
}
