package middlewares

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/uuid"
)

type traceIdKey struct{}
type spanIdKey struct{}

func generateId() (string, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func TraceIdInjector() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			traceId, err := generateId()
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, traceIdKey{}, traceId)
			return handler(ctx, req)
		}
	}
}
