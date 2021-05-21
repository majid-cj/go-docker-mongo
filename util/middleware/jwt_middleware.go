package middleware

import (
	"github.com/majid-cj/go-docker-mongo/infrastructure/auth"
	"github.com/majid-cj/go-docker-mongo/util"

	"github.com/kataras/iris/v12"
)

// AuthenticationJWTMiddleware ...
func AuthenticationJWTMiddleware(c iris.Context) {
	err := auth.TokenValid(c.Request())
	if err != nil {
		util.ResponseT(util.GetError("unauthorized_access"), iris.StatusUnauthorized, c)
		return
	}
	c.Next()
}
