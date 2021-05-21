package util

import (
	"errors"
)

var (
	ErrPasswordFormat = errors.New("invalid_password")
)

// ValidatePassword ...
func ValidatePassword(password string) error {
	if len(password) < 6 {
		return ErrPasswordFormat
	}
	return nil
}
