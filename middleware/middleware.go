package middleware

import (
	"crypto/rsa"
	"github.com/dequinox/alifcore-auth-middleware/config"
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

	Keys   *rsa.PublicKey
	Config config.Config
}

type Middleware interface {
	Middleware(nextHandler gin.HandlerFunc, roles ...string) gin.HandlerFunc
}

type middleware struct {
	keys   *rsa.PublicKey
	config config.Config
}

func NewMiddleware(p Params) (Middleware, error) {

	middlewareItem := &middleware{
		keys:   p.Keys,
		config: p.Config,
	}

	return middlewareItem, nil
}

type TokenPayload struct {
	Issuer    string              `json:"iss,omitempty"`
	Subject   int64               `json:"sub,omitempty"`
	Expiry    int64               `json:"exp,omitempty"`
	NotBefore int64               `json:"nbf,omitempty"`
	IssuedAt  int64               `json:"iat,omitempty"`
	ID        string              `json:"jti,omitempty"`
	Username  string              `json:"username,omitempty"`
	Roles     map[string][]string `json:"roles"`
}

func (m *middleware) Middleware(nextHandler gin.HandlerFunc, roles ...string) gin.HandlerFunc {
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

		service := m.config.GetString("SERVICE_NAME")

		if !HasRoles(service, roles, out.Roles) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "access denied, user does not have required roles"})
			return
		}

		c.Set("id", out.Subject)
		c.Set("username", out.Username)
		c.Set("roles", out.Roles)

		nextHandler(c)
	}
}

func HasRoles(serviceName string, expectedRoles []string, actualRoles map[string][]string) bool {
	if service, ok := actualRoles[serviceName]; ok {

		var shouldHave = make(map[string]bool)
		for _, role := range service {
			shouldHave[role] = true
		}

		for _, role := range expectedRoles {
			if _, ok := shouldHave[role]; !ok {
				return false
			}
		}

		return true
	}

	return false
}
