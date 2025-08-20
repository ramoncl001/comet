package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
)

func HMAC_SHA256(content, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(content))
	return hex.EncodeToString(h.Sum(nil))
}

func SHA256(content string) string {
	hash := sha256.Sum256([]byte(content))
	return hex.EncodeToString(hash[:])
}
