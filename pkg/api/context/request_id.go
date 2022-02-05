package context

import (
	"context"
	"fmt"
)

const RequestIDKey = "maestro.requestId"

func WithRequestID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, RequestIDKey, reqID)
}

func RequestID(ctx context.Context) (string, error) {
	val := ctx.Value(RequestIDKey)
	reqID, ok := val.(string)
	if ok {
		return reqID, nil
	}

	return "", fmt.Errorf("could not find request id with key \"%s\" in context", RequestIDKey)
}
