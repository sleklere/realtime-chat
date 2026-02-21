package auth

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type claimsKey struct{}

// ClaimsFromCtx extracts the JWT claims stored in the request context by validateJWT.
func ClaimsFromCtx(ctx context.Context) (*Claims, bool) {
	c, ok := ctx.Value(claimsKey{}).(*Claims)
	return c, ok
}

// NewClaimsContext returns a context with the given claims stored in it.
func NewClaimsContext(ctx context.Context, c *Claims) context.Context {
	return context.WithValue(ctx, claimsKey{}, c)
}

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
