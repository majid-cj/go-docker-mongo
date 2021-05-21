package apps

import (
	"github.com/majid-cj/go-docker-mongo/domain/entity"
	"github.com/majid-cj/go-docker-mongo/domain/repository"
)

// VerifyCodeApp ...
type VerifyCodeApp struct {
	code repository.VerificationCodeRepository
}

// VerifyCodeAppInterface ...
type VerifyCodeAppInterface interface {
	CreateVerificationCode(*entity.VerificationCode) (*entity.VerificationCode, error)
	CreateVerificationCodeFromEmail(*entity.VerificationCode) (*entity.VerificationCode, error)
	ResetPassword(*entity.VerificationCode) error
	CheckVerificationCode(*entity.VerificationCode) error
	RenewVerificationCode(*entity.VerificationCode) (*entity.VerificationCode, error)
}

var _ VerifyCodeAppInterface = &VerifyCodeApp{}

// CreateVerificationCode ...
func (vca *VerifyCodeApp) CreateVerificationCode(verifycode *entity.VerificationCode) (*entity.VerificationCode, error) {
	return vca.code.CreateVerificationCode(verifycode)
}

// CreateVerificationCodeFromEmail ...
func (vca *VerifyCodeApp) CreateVerificationCodeFromEmail(verifycode *entity.VerificationCode) (*entity.VerificationCode, error) {
	return vca.code.CreateVerificationCodeFromEmail(verifycode)
}

// ResetPassword ...
func (vca *VerifyCodeApp) ResetPassword(verifycode *entity.VerificationCode) error {
	return vca.code.ResetPassword(verifycode)
}

// CheckVerificationCode ...
func (vca *VerifyCodeApp) CheckVerificationCode(verifycode *entity.VerificationCode) error {
	return vca.code.CheckVerificationCode(verifycode)
}

// RenewVerificationCode ...
func (vca *VerifyCodeApp) RenewVerificationCode(verifycode *entity.VerificationCode) (*entity.VerificationCode, error) {
	return vca.code.RenewVerificationCode(verifycode)
}
