package middleware

import (
	"github.com/ramoncl001/comet/log"
	"github.com/ramoncl001/comet/rest"
)

var RequestLogging Middleware = func(next rest.RequestHandler) rest.RequestHandler {
	return func(req *rest.Request) rest.Response {
		logger := log.FromContext(req.Context())

		logger.Debug("request received", "method", req.Method, "path", req.Url.Path)

		return next(req)
	}
}
