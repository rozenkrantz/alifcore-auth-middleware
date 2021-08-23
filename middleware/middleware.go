package middleware

import (
	"crypto/rsa"
	"github.com/gin-gonic/gin"
	"go.uber.org/fx"
	"gopkg.in/square/go-jose.v2/jwt"
	"net/http"
	"strings"
	"time"
)

var Module = fx.Options(
	fx.Provide(NewMiddleware),
)

type Params struct {
	fx.In

	Keys *rsa.PublicKey
}

type Middleware interface {
	Middleware(nextHandler gin.HandlerFunc) gin.HandlerFunc
}

type middleware struct {
	keys *rsa.PublicKey
}

func NewMiddleware(p Params) (Middleware, error) {

	middlewareItem := &middleware{
		keys: p.Keys,
	}

	return middlewareItem, nil
}

type TokenPayload struct {
	Issuer    string                 `json:"iss,omitempty"`
	Subject   int64                  `json:"sub,omitempty"`
	Expiry    int64                  `json:"exp,omitempty"`
	NotBefore int64                  `json:"nbf,omitempty"`
	IssuedAt  int64                  `json:"iat,omitempty"`
	ID        string                 `json:"jti,omitempty"`
	Username  string                 `json:"username,omitempty"`
	Roles     map[string]interface{} `json:"roles"`
}

func (m *middleware) Middleware(nextHandler gin.HandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {

		raw := strings.TrimPrefix(c.GetHeader("Authorization"), "Bearer ")

		tok, err := jwt.ParseSigned(raw)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization failed: " + err.Error()})
			return
		}

		var out TokenPayload

		if err := tok.Claims(m.keys, &out); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token " + err.Error()})
			return
		}

		if time.Unix(out.Expiry, 0).Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Expired Token"})
			return
		}

		c.Set("id", out.Subject)
		c.Set("username", out.Username)
		c.Set("roles", out.Roles)

		nextHandler(c)
	}
}
