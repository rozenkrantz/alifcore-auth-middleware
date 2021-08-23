# alifcore-auth-middleware

alifcore-auth-middleware is a package used to authenticate users from alif-service-user.

## How to use
Include the module in your fx.Options when running the app:
```go

import "github.com/dequinox/alifcore-auth-middleware/middleware"

...
modules := fx.Options(
	...
	alifcore_auth_middleware.Module,
	...
)

fx.New(modules).Run()
...
```

Pass the middleware on route endpoints:
```go

import "github.com/dequinox/alifcore-auth-middleware/middleware"

// Params is the input parameter struct for the module that contains its dependencies
type Params struct {
    fx.In
    Mw     middleware.Middleware
    Srv    *gin.Engine
}

// NewPingHandler constructs a new ping.Handler.
func NewPingHandler(p Params) error {

    mw := p.Mw.Middleware
    p.Srv.GET("/ping", mw(Ping))
    p.Srv.Get("/ping-access", mw(Ping, "admin", "moder")) // Endpoint requiring admin and moder access

    return nil
}

func Ping(c *gin.Context) {

    c.JSON(200, gin.H{
        "message": "pong",
    })

}
```

## Environment Variables
Place .env file in the base of your project
```dotenv
PUB_KEY_URI="http://alif-core-service-user.eu-central-1.elasticbeanstalk.com/service_user/auth/public_key"
PUB_KEY_DATA="{\"service_name\": \"alif-shop-settings\"}"
SERVICE_NAME="alif-shop-settings"
```

## Checking info
3 values are passed from middleware:
```go
id       string
username string
roles    map[string][]string
```
Get values passed from middleware in handlers:
```go
func (t *handler) foobar(c *gin.Context) {
    if id, ok := c.Get("id"); ok {
        ...
    }
    ...
```