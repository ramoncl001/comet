package api

import (
	"github.com/ramoncl001/comet/middleware"
	"github.com/ramoncl001/comet/rest"
)

func chain(handler rest.RequestHandler, middlewares ...middleware.Middleware) rest.RequestHandler {
	for i := len(middlewares) - 1; i >= 0; i-- {
		handler = middlewares[i](handler)
	}

	return handler
}
