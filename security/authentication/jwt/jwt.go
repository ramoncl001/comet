package jwt

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/ramoncl001/comet/ioc"
	"github.com/ramoncl001/comet/rest"
	"github.com/ramoncl001/comet/security"
	"github.com/ramoncl001/comet/security/authentication"
)

var errInvalidToken = errors.New("invalid token")

const (
	ClaimAudience  = "aud"
	ClaimExpiresAt = "exp"
	ClaimID        = "jti"
	ClaimIssuedAt  = "iat"
	ClaimIssuer    = "iss"
	ClaimNotBefore = "nbf"
	ClaimSubject   = "sub"
	ClaimSessionID = "sid"
)

type JwtProvider interface {
	GenerateToken(claims authentication.Claims, secret string) string
	ValidateToken(token string, secret string) (authentication.Claims, error)
}

type JwtConfigurations struct {
	Issuer     string
	Audience   string
	Expiration int64
	SecretKey  string
}

type DefaultJwtSessionManager struct {
	authentication.SessionManager
	config      JwtConfigurations
	provider    JwtProvider
	userManager security.UserManager
}

func NewDefaultJwtSessionManager(config JwtConfigurations, provider JwtProvider, userManager security.UserManager) authentication.SessionManager {
	return &DefaultJwtSessionManager{
		config:      config,
		provider:    provider,
		userManager: userManager,
	}
}

func (sm *DefaultJwtSessionManager) Validate(req *rest.Request) (authentication.Claims, error) {
	authHeader := req.Headers["Authorization"][0]
	if authHeader == "" {
		return nil, errInvalidToken
	}
	return sm.provider.ValidateToken(strings.ReplaceAll(authHeader, "Bearer ", ""), sm.config.SecretKey)
}

func (sm *DefaultJwtSessionManager) GetUser(ctx context.Context) (security.ApplicationUser, error) {
	id := ctx.Value("user_id")
	if id == nil {
		return nil, errors.New("session not started")
	}

	user := sm.userManager.GetByID(fmt.Sprintf("%v", id))
	if user == nil {
		return nil, errors.New("user does not exists")
	}

	return user, nil
}

func (sm *DefaultJwtSessionManager) GetToken(claims authentication.Claims) string {
	return sm.provider.GenerateToken(claims, sm.config.SecretKey)
}

var DefaultJwtAuthenticationMiddleware = func(next rest.RequestHandler) rest.RequestHandler {
	return func(req *rest.Request) rest.Response {
		manager, err := ioc.ResolveTransient[authentication.SessionManager](req.Context())
		if err != nil {
			return rest.Error("error resolving dependency")
		}

		claims, err := manager.Validate(req)
		if err != nil {
			return rest.Unauthorized()
		}

		req = req.WithContext(context.WithValue(req.Context(), "user_id", claims["sub"]))

		return next(req)
	}
}
