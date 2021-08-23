package alifcore_auth_middleware

import (
	"github.com/dequinox/alifcore-auth-middleware/config"
	"github.com/dequinox/alifcore-auth-middleware/keys"
	"github.com/dequinox/alifcore-auth-middleware/middleware"
	"go.uber.org/fx"
)

var Module = fx.Options(
	config.Module,
	keys.Module,
	middleware.Module,
)
