package jwt

import (
	"crypto/hmac"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ramoncl001/comet/security"
	"github.com/ramoncl001/comet/security/authentication"
)

type defaultJwtProvider struct {
	JwtProvider
}

func NewDefaultJwtProvider() JwtProvider {
	return &defaultJwtProvider{}
}

func (pv *defaultJwtProvider) GenerateToken(claims authentication.Claims, secret string) string {
	header := struct {
		Alg string `json:"alg"`
		Typ string `json:"typ"`
	}{
		Alg: "HS256",
		Typ: "JWT",
	}

	headerJSON, err := json.Marshal(header)
	if err != nil {
		panic(err)
	}
	headerEncoded := base64.RawURLEncoding.EncodeToString(headerJSON)

	claimsJSON, err := json.Marshal(claims)
	if err != nil {
		panic(err)
	}
	claimsEncoded := base64.RawURLEncoding.EncodeToString(claimsJSON)

	signatureData := fmt.Sprintf("%s.%s", headerEncoded, claimsEncoded)
	signature := security.HMAC_SHA256(signatureData, secret)

	return fmt.Sprintf("%s.%s.%s", headerEncoded, claimsEncoded, signature)
}

func (pv *defaultJwtProvider) ValidateToken(token string, secret string) (authentication.Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, errInvalidToken
	}

	signatureData := fmt.Sprintf("%s.%s", parts[0], parts[1])
	calculatedSig := security.HMAC_SHA256(signatureData, secret)

	if !hmac.Equal([]byte(calculatedSig), []byte(parts[2])) {
		return nil, errInvalidToken
	}

	claimsJSON, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, errInvalidToken
	}

	var claims authentication.Claims
	if err := json.Unmarshal(claimsJSON, &claims); err != nil {
		return nil, errInvalidToken
	}

	currentTime := time.Now().Unix()

	if exp, ok := claims[ClaimExpiresAt].(float64); ok {
		if currentTime > int64(exp) {
			return nil, errInvalidToken
		}
	}

	if nbf, ok := claims[ClaimNotBefore].(float64); ok {
		if currentTime < int64(nbf) {
			return nil, errInvalidToken
		}
	}

	return claims, nil
}
