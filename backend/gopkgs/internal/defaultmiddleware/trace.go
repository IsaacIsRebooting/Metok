package defaultmiddleware

import (
	"context"
	"errors"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/metadata"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/google/uuid"
)

const (
	TraceIdHeadersKey = "x-trace-id"
	SpanIdHeadersKey  = "x-span-id"
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

// make sure every request has a traceId
func TraceIdInjector() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			meta, ok := metadata.FromServerContext(ctx)
			if !ok {
				return nil, errors.New("get metadata failed")
			}
			traceId := meta.Get(TraceIdHeadersKey)
			if traceId == "" {
				traceId, err := generateId()
				if err != nil {
					return nil, err
				}
				metadata.AppendToClientContext(ctx, TraceIdHeadersKey, traceId)
			}
			ctx = context.WithValue(ctx, traceIdKey{}, traceId)
			return handler(ctx, req)
		}
	}
}

// get traceId from context
func GetTraceId() log.Valuer {
	return func(ctx context.Context) interface{} {
		traceId := ctx.Value(traceIdKey{})
		if traceId == nil {
			return ""
		}
		return traceId
	}
}

// make sure every request has a spanId
func SpanIdInjector() middleware.Middleware {
	return func(handler middleware.Handler) middleware.Handler {
		return func(ctx context.Context, req interface{}) (interface{}, error) {
			spanId, err := generateId()
			if err != nil {
				return nil, err
			}
			ctx = context.WithValue(ctx, spanIdKey{}, spanId)
			return handler(ctx, req)
		}
	}
}

// get spanId from context
func GetSpanId() log.Valuer {
	return func(ctx context.Context) interface{} {
		spanId := ctx.Value(spanIdKey{})
		if spanId == nil {
			return ""
		}
		return spanId
	}
}

// func extractError(err error) (log.Level, string) {
// 	if err != nil {
// 		return log.LevelError, fmt.Sprintf("%+v", err)
// 	}
// 	return log.LevelInfo, ""
// }
