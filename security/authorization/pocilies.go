package authorization

import (
	"github.com/ramoncl001/comet/ioc"
	"github.com/ramoncl001/comet/rest"
	"github.com/ramoncl001/comet/security/authentication"
)

var RequireRole = func(next rest.RequestHandler, value interface{}) rest.RequestHandler {
	return func(req *rest.Request) rest.Response {
		sessionManager, err := ioc.ResolveSingleton[authentication.SessionManager](req.Context())
		if err != nil {
			return rest.Error("could not resolve session manager")
		}

		claims, err := sessionManager.Validate(req)
		if err != nil {
			return rest.Unauthorized()
		}

		if claims["role"] != value {
			return rest.Unauthorized()
		}

		return next(req)
	}
}
