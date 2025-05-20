package middleware

import (
	"context"

	"server-template/pkg"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/transport"
	"github.com/go-kratos/kratos/v2/transport/http"
)

type RequestIDKey struct{}

const (
	TraceIDHeader = "X-Request-ID"
	CfRay         = "CF-Ray"
	RequestIDHeader = "request_id"
)

func RequestIdHandler(handler middleware.Handler) middleware.Handler {
	return func(ctx context.Context, req interface{}) (reply interface{}, err error) {
		var traceID string
		if tr, ok := transport.FromServerContext(ctx); ok {
			if ht, ok := tr.(http.Transporter); ok {
				reqHeader := ht.RequestHeader()
				traceID = reqHeader.Get(TraceIDHeader)
				if traceID == "" {
					traceID = reqHeader.Get(RequestIDHeader)
				}
				if traceID == "" {
					traceID = reqHeader.Get(CfRay)
				}
			}
		}
		if traceID == "" {
			traceID = pkg.GenShortID()
		}
		ctx = context.WithValue(ctx, RequestIDKey{}, traceID)
		return handler(ctx, req)
	}
}

func RequestId() log.Valuer {
	return func(ctx context.Context) any {
		requestId := ctx.Value(RequestIDKey{})
		if requestId == nil {
			return ""
		}
		return requestId
	}
}
