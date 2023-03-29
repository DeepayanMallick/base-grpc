package middleware

import (
	"runtime/debug"

	"github.com/sirupsen/logrus"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type recoveryHandler struct {
	logger       *logrus.Entry
	slackHookURL string
}

func newRecoveryHandler(l *logrus.Entry) *recoveryHandler {
	return &recoveryHandler{
		logger:       l,
	}
}

func (r *recoveryHandler) recover(p interface{}) error {
	var entry *logrus.Entry

	if err, ok := p.(error); ok {
		// runtime panic
		r.logger.WithError(err)
	} else {
		// explicit panic call
		r.logger.WithFields(logrus.Fields{
			"panic": p,
		})
	}

	stack := string(debug.Stack())
	entry.WithFields(logrus.Fields{
		"stacktrace": stack,
	}).Error("recovered from panic")

	return status.Errorf(codes.Internal, "internal error occurred")
}
