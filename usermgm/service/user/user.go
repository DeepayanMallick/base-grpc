package user

import (
	"context"

	"github.com/DeepayanMallick/base-grpc/usermgm/storage"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"

	upb "github.com/DeepayanMallick/base-grpc/gunk/v1/usermgm/user"
)

type Handler struct {
	upb.UnimplementedUserServiceServer
	usr    CoreUser
	logger *logrus.Entry
}

type CoreUser interface {
	CreateUser(context.Context, storage.User) (string, error)
	GetUserByEmail(context.Context, string) (*storage.User, error)
}

func New(usr CoreUser, logger *logrus.Entry) *Handler {
	return &Handler{usr: usr, logger: logger}
}

// // RegisterService with grpc server.
func (h *Handler) RegisterSvc(srv *grpc.Server) error {
	upb.RegisterUserServiceServer(srv, h)
	return nil
}

type resAct struct {
	res, act string
	pub      bool
}

func (h *Handler) Permission(ctx context.Context, mthd string) (resource, action string, pub bool) {
	p := map[string]resAct{
		"CreateUser":     {res: "user", act: "create"},
		"GetUserByEmail": {res: "user", act: "view"},
	}
	return p[mthd].res, p[mthd].act, p[mthd].pub
}
