package middleware

import (
	"context"

	"github.com/google/uuid"
	"github.com/ramoncl001/comet/rest"
)

var RequestID Middleware = func(next rest.RequestHandler) rest.RequestHandler {
	return func(req *rest.Request) rest.Response {
		ctx := context.WithValue(req.Context(), "X-Request-Id", uuid.New().String())
		return next(req.WithContext(ctx))
	}
}
