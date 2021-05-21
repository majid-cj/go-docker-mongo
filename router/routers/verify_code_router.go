package routers

import (
	"github.com/majid-cj/go-docker-mongo/apps"
	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/util"

	"github.com/kataras/iris/v12"
)

// VerifyCodeRouter ...
type VerifyCodeRouter struct {
	vca apps.VerifyCodeAppInterface
}

// NewVerifyCodeRouter ...
func NewVerifyCodeRouter(vca apps.VerifyCodeAppInterface) *VerifyCodeRouter {
	return &VerifyCodeRouter{
		vca: vca,
	}
}

// NewVerifyCode ...
func (vcr *VerifyCodeRouter) NewVerifyCode(c iris.Context) {
	var verifycode entity.VerificationCode
	err := c.ReadJSON(&verifycode)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	verifycode.PrepareVerificationCode(verifycode.Member, verifycode.CodeType)
	code, err := vcr.vca.CreateVerificationCode(&verifycode)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(code, iris.StatusCreated, c)
}

// VerificationCodeFromEmail ...
func (vcr *VerifyCodeRouter) VerificationCodeFromEmail(c iris.Context) {
	var verifycode entity.VerificationCode
	err := c.ReadJSON(&verifycode)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	_, err = vcr.vca.CreateVerificationCodeFromEmail(&verifycode)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusCreated, c)
}

// ResetPasswordVerifyCode ...
func (vcr *VerifyCodeRouter) ResetPasswordVerifyCode(c iris.Context) {
	var verifycode entity.VerificationCode
	err := c.ReadJSON(&verifycode)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}

	err = vcr.vca.ResetPassword(&verifycode)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusOK, c)
}

// CheckVerifyCode ...
func (vcr *VerifyCodeRouter) CheckVerifyCode(c iris.Context) {
	var verifyCode entity.VerificationCode
	err := c.ReadJSON(&verifyCode)
	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusUnprocessableEntity, c)
		return
	}

	err = vcr.vca.CheckVerificationCode(&verifyCode)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(nil, iris.StatusOK, c)
}

// RenewVerifyCode ...
func (vcr *VerifyCodeRouter) RenewVerifyCode(c iris.Context) {
	var verifycode entity.VerificationCode
	err := c.ReadJSON(&verifycode)

	if err != nil {
		util.ResponseT(util.GetError("error_parsing_data"), iris.StatusBadRequest, c)
		return
	}
	code, err := vcr.vca.RenewVerificationCode(&verifycode)
	if err != nil {
		util.ResponseT(err, iris.StatusBadRequest, c)
		return
	}
	util.Response(code, iris.StatusCreated, c)
}
