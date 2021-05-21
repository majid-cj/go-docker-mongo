package router

import (
	"github.com/majid-cj/go-docker-mongo/infrastructure/auth"
	"github.com/majid-cj/go-docker-mongo/infrastructure/persistence"
	"github.com/majid-cj/go-docker-mongo/router/routers"
	"github.com/majid-cj/go-docker-mongo/util/middleware"

	"github.com/kataras/iris/v12"
)

// API ...
func API(
	app *iris.Application,
	repositoryservice *persistence.Repository,
	authservice *auth.DBAuth,
) {
	token := auth.NewToken()

	authentication := routers.NewAuthenticationRouter(
		repositoryservice.Member,
		repositoryservice.VerifyCode,
		authservice.Auth,
		token,
	)
	verifycode := routers.NewVerifyCodeRouter(
		repositoryservice.VerifyCode,
	)

	v := app.Party("/api/v1")
	{
		user := v.Party("/user")
		{
			user.Post("/sign-up", authentication.SignUp)
			user.Post("/sign-in", authentication.SignIn)

			user.Post("/reset/code", verifycode.VerificationCodeFromEmail)
			user.Post("/reset/password", verifycode.ResetPasswordVerifyCode)

			user.Post("/refresh", authentication.Refresh)
			user.Use(middleware.AuthenticationJWTMiddleware, middleware.AuthMemberMiddleware)
			user.Put("/password/{id:string}", authentication.UpdatePassword)
			user.Post("/logout", authentication.Logout)

			user.Post("/verify/code", verifycode.NewVerifyCode)
			user.Post("/verify/check", verifycode.CheckVerifyCode)
			user.Post("/verify/renew", verifycode.RenewVerifyCode)
			user.Post("/verify/reset/password", verifycode.ResetPasswordVerifyCode)
		}
	}
}
