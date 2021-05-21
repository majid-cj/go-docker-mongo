package util

import (
	"errors"
	"strings"
)

var (
	ErrEmailFormat = errors.New("invalid_display_name")
)

// ValidateDisplayName ...
func ValidateDisplayName(name string) error {
	if len(strings.TrimSpace(name)) == 0 {
		return ErrEmailFormat
	}
	return nil
}
