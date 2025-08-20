package authentication

import (
	"context"

	"github.com/ramoncl001/comet/rest"
	"github.com/ramoncl001/comet/security"
)

type Claims map[string]interface{}

type SessionManager interface {
	GetToken(claims Claims) string
	Validate(req *rest.Request) (Claims, error)
	GetUser(ctx context.Context) (security.ApplicationUser, error)
}
