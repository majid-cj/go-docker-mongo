package util

import (
	"github.com/twinj/uuid"
)

// UUID ...
func UUID() string {
	return uuid.NewV4().String()
}
