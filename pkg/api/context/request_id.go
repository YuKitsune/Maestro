package context

import (
	"context"
	"fmt"
)

const RequestIdKey = "maestro.requestId"

func WithRequestId(ctx context.Context, reqId string) context.Context {
	return context.WithValue(ctx, RequestIdKey, reqId)
}

func RequestId(ctx context.Context) (string, error) {
	val := ctx.Value(RequestIdKey)
	reqId, ok := val.(string)
	if ok {
		return reqId, nil
	}

	return "", fmt.Errorf("could not find request id with key \"%s\" in context", RequestIdKey)
}
