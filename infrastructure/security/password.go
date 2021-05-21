package security

import (
	"github.com/majid-cj/go-docker-mongo/util"
	"golang.org/x/crypto/bcrypt"
)

// HashPassword ...
func HashPassword(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// VerifyPassword ...
func VerifyPassword(hash, password string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		return util.GetError("password_mismatch")
	}
	return nil
}
