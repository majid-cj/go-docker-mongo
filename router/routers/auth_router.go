package routers

import (
	"fmt"
	"os"

	"github.com/majid-cj/go-docker-mongo/apps"
	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/infrastructure/auth"
	"github.com/majid-cj/go-docker-mongo/infrastructure/security"
	"github.com/majid-cj/go-docker-mongo/util"
	"github.com/majid-cj/go-docker-mongo/util/fileupload"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris/v12"
)

// AuthenticationRouter ...
type AuthenticationRouter struct {
	memberapp  apps.MemberAppInterface
	auth       auth.AuthenticationInterface
	token      auth.TokenInterface
	fileupload fileupload.UploadFileInterface
	verify     apps.VerifyCodeAppInterface
}

// NewAuthenticationRouter ...
func NewAuthenticationRouter(
	memberapp apps.MemberAppInterface,
	verify apps.VerifyCodeAppInterface,
	auth auth.AuthenticationInterface,
	token auth.TokenInterface,
) *AuthenticationRouter {
	return &AuthenticationRouter{
		memberapp:  memberapp,
		verify:     verify,
		auth:       auth,
		token:      token,
	}
}

// SignUp ...
func (router *AuthenticationRouter) SignUp(c iris.Context) {
	var member entity.Member
	var code entity.VerificationCode

	err := c.ReadJSON(&member)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	err = member.ValidateSignUp()
	if err != nil {
		util.ResponseT(err, iris.StatusUnprocessableEntity, c)
		return
	}
	password, err := security.HashPassword(member.Password)
	if err != nil {
		util.ResponseT(util.GetError("general_error"), iris.StatusUnprocessableEntity, c)
		return
	}
	member.PrepareMember(password)
	newMember, err := router.memberapp.CreateMember(&member)
	if err != nil {
		util.ResponseT(err, iris.StatusInternalServerError, c)
		return
	}

	token, err := router.token.CreateJWTToken(
		newMember.ID,
		fmt.Sprintf("%d", newMember.Type),
	)
	if err != nil {
		util.ResponseT(util.GetError("general_error"), iris.StatusUnprocessableEntity, c)
		return
	}

	createError := router.auth.CreatToken(newMember.ID, token)
	if createError != nil {
		util.ResponseT(util.GetError("general_error"), iris.StatusInternalServerError, c)
		return
	}

	code.PrepareVerificationCode(member.ID, 1)
	_, err = router.verify.CreateVerificationCode(&code)
	if err != nil {
		util.ResponseT(util.GetError("general_error"), iris.StatusInternalServerError, c)
		return
	}

	memberResponse := make(map[string]interface{})

	memberResponse["token"] = token
	memberResponse["member"] = newMember.GetMemberSerializer()

	util.Response(memberResponse, iris.StatusCreated, c)
}

// SignIn ...
func (router *AuthenticationRouter) SignIn(c iris.Context) {
	var member *entity.Member
	var code entity.VerificationCode

	err := c.ReadJSON(&member)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	err = member.ValidateSignIn()
	if err != nil {
		util.ResponseT(err, iris.StatusUnprocessableEntity, c)
		return
	}
	memberLogin, err := router.memberapp.GetMemberByEmailAndPassword(member)
	if err != nil {
		util.ResponseT(err, iris.StatusNotFound, c)
		return
	}
	token, err := router.token.CreateJWTToken(
		memberLogin.ID,
		fmt.Sprintf("%d", memberLogin.Type),
	)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}

	err = router.auth.CreatToken(memberLogin.ID, token)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}

	if !memberLogin.Verified {
		code.PrepareVerificationCode(memberLogin.ID, 1)
		_, err = router.verify.CreateVerificationCode(&code)
		if err != nil {
			util.ResponseT(util.GetError("general_error"), iris.StatusInternalServerError, c)
			return
		}
	}

	memberResponse := make(map[string]interface{})
	memberResponse["token"] = token
	memberResponse["member"] = memberLogin.GetMemberSerializer()

	util.Response(memberResponse, iris.StatusOK, c)
}

// Logout ...
func (router *AuthenticationRouter) Logout(c iris.Context) {
	token, err := router.token.ExtractJWTTokenMetadata(c.Request())
	if err != nil {
		util.ResponseT(err, iris.StatusUnauthorized, c)
		return
	}
	err = router.auth.DeleteAccessToken(token)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusOK, c)
}

// Refresh ...
func (router *AuthenticationRouter) Refresh(c iris.Context) {
	var data struct {
		Refresh string `json:"refresh"`
	}

	err := c.ReadJSON(&data)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	refreshToken := data.Refresh
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, util.GetError("general_error")
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	if err != nil {
		util.ResponseT(err, iris.StatusUnauthorized, c)
		return
	}

	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		util.ResponseT(err, iris.StatusUnauthorized, c)
		return
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		refreshUUID, ok := claims["refresh_uuid"].(string)
		if !ok {
			util.ResponseT(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}

		userID, ok := claims["user_id"].(string)
		if !ok {
			util.ResponseT(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}

		userType, ok := claims["user_type"].(string)
		if !ok {
			util.ResponseT(util.GetError("general_error"), iris.StatusUnauthorized, c)
			return
		}

		err := router.auth.DeleteRefreshToken(refreshUUID)
		if err != nil {
			util.ResponseT(err, iris.StatusUnauthorized, c)
			return
		}
		newToken, err := router.token.CreateJWTToken(userID, userType)
		if err != nil {
			util.ResponseT(err, iris.StatusUnauthorized, c)
			return
		}

		err = router.auth.CreatToken(userID, newToken)
		if err != nil {
			util.ResponseT(err, iris.StatusUnauthorized, c)
			return
		}

		util.Response(newToken, iris.StatusOK, c)
	} else {
		util.ResponseT(util.GetError("general_error"), iris.StatusUnauthorized, c)
		return
	}
}

// UpdatePassword ...
func (router *AuthenticationRouter) UpdatePassword(c iris.Context) {
	var data struct {
		Password        string `json:"password"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	err := c.ReadJSON(&data)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	member, err := router.memberapp.GetMember(c.Params().GetString("id"))
	if err != nil {
		util.ResponseT(err, iris.StatusNotFound, c)
		return
	}
	err = security.VerifyPassword(member.Password, data.Password)
	if err != nil {
		util.ResponseT(err, iris.StatusUnauthorized, c)
		return
	}

	if data.NewPassword != data.ConfirmPassword {
		util.ResponseT(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	password, err := security.HashPassword(data.ConfirmPassword)
	if err != nil {
		util.ResponseT(util.GetError("general_error"), iris.StatusBadRequest, c)
		return
	}

	member.Password = util.EscapeString(string(password))
	member.UpdateAt = util.GetTimeNow()

	err = member.ValidateSignUp()
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}

	err = router.memberapp.UpdatePassword(member)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}

	util.Response(nil, iris.StatusOK, c)
}
