package jwt

import "time"

const (
	IssuerId              = "kript.api.account"
	RefreshTokenType      = "refresh"
	RefreshTokenLife      = time.Hour * 24 * 365 * 100
	AccessTokenType       = "access"
	AccessTokenLife       = time.Hour * 24
	VerificationTokenType = "verify"
	VerificationTokenLife = time.Minute * 30
)
