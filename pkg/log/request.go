package log

import (
	"github.com/sirupsen/logrus"
	mcontext "github.com/yukitsune/maestro/pkg/api/context"
	"net/http"
)

const RequestIDField = "request_id"

func ForRequest(logger *logrus.Logger, req *http.Request) (*logrus.Entry, error) {

	ctx := req.Context()
	requestId, err := mcontext.RequestID(ctx)
	if err != nil {
		return nil, err
	}

	return logger.WithContext(ctx).WithField(RequestIDField, requestId), nil
}
