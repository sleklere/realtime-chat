package auth

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Config holds parameters used to issue and validate JWT access tokens.
type Config struct {
	JWTSecret []byte
	Issuer    string
	AccessTTL time.Duration
}

// Claims contains application-specific JWT claims embedded in access tokens.
type Claims struct {
	UserID   int64  `json:"uid"`
	Username string `json:"un"`
	jwt.RegisteredClaims
}
