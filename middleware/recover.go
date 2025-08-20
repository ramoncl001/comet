package middleware

import (
	"runtime/debug"
	"time"

	"github.com/ramoncl001/comet/log"
	"github.com/ramoncl001/comet/rest"
)

type panicLog struct {
	Timestamp     time.Time   `json:"timestamp"`
	URL           string      `json:"url"`
	Method        string      `json:"method"`
	Error         interface{} `json:"error"`
	StackTrace    string      `json:"stack_trace"`
	UserAgent     string      `json:"user_agent"`
	RemoteAddress string      `json:"remote_addr"`
}

var Recover Middleware = func(next rest.RequestHandler) rest.RequestHandler {
	return func(req *rest.Request) rest.Response {
		defer func(req *rest.Request) {
			if err := recover(); err != nil {
				panicLog := panicLog{
					Timestamp:     time.Now(),
					URL:           req.Url.String(),
					Method:        req.Method,
					Error:         err,
					StackTrace:    string(debug.Stack()),
					UserAgent:     req.UserAgent,
					RemoteAddress: req.RemoteAddress,
				}

				logger := log.FromContext(req.Context())

				logger.Error("panic error received from request", "info", panicLog)
			}
		}(req)
		return next(req)
	}
}
