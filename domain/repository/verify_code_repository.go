package repository

import "github.com/majid-cj/go-docker-mongo/domain/entity"

// VerificationCodeRepository ...
type VerificationCodeRepository interface {
	CreateVerificationCode(*entity.VerificationCode) (*entity.VerificationCode, error)
	CreateVerificationCodeFromEmail(*entity.VerificationCode) (*entity.VerificationCode, error)
	ResetPassword(*entity.VerificationCode) error
	CheckVerificationCode(*entity.VerificationCode) error
	RenewVerificationCode(*entity.VerificationCode) (*entity.VerificationCode, error)
}
