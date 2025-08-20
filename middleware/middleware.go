package middleware

import (
	"github.com/ramoncl001/comet/rest"
)

type Middleware func(next rest.RequestHandler) rest.RequestHandler
