package middleware

import (
	"github.com/majid-cj/go-docker-mongo/infrastructure/auth"
	"github.com/majid-cj/go-docker-mongo/util"

	"github.com/kataras/iris/v12"
)

// AuthMemberMiddleware ...
func AuthMemberMiddleware(c iris.Context) {
	membertype, err := auth.ExtractMemberType(c.Request())
	if err != nil {
		util.ResponseT(err, iris.StatusUnauthorized, c)
		return
	}
	if membertype != "2" {
		util.ResponseT(util.GetError("unauthorized_access"), iris.StatusUnauthorized, c)
		return
	}
	c.Next()
}

// AuthAdminMiddleware ...
func AuthAdminMiddleware(c iris.Context) {
	membertype, err := auth.ExtractMemberType(c.Request())
	if err != nil {
		util.ResponseT(err, iris.StatusUnauthorized, c)
		return
	}
	if membertype != "1" {
		util.ResponseT(util.GetError("unauthorized_access"), iris.StatusUnauthorized, c)
		return
	}
	c.Next()
}
