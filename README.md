# alifcore-auth-middleware

alifcore-auth-middleware is a package used to authenticate users from alif-service-user.

## How to use
Include the module in your fx.Options when running the app:
```go
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
// Params is the input parameter struct for the module that contains its dependencies
type Params struct {
    fx.In
    Mw     alifcore_auth_middleware.Middleware
    Srv    *gin.Engine
}

// NewPingHandler constructs a new ping.Handler.
func NewPingHandler(p Params) error {

    mw := p.Mw.Middleware
    p.Srv.GET("/ping", mw(Ping))

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
```