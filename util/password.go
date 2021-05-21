package util

import (
	"errors"
	"regexp"
)

var (
	ErrPasswordFormat = errors.New("invalid_password")
	PasswordRegex     = "^(?=.*[A-Za-z])(?=.*\\d)(?=.*[@$!%*#?&])[A-Za-z\\d@$!%*#?&]{6,}$"
)

// ValidatePassword ...
func ValidatePassword(password string) error {
	if _, err := regexp.Match(PasswordRegex, []byte(password)); err != nil {
		return ErrPasswordFormat
	}
	return nil
}
