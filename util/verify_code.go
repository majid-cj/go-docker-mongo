package util

import (
	"fmt"
	"math/rand"
	"time"
)

// VerifyCode ...
func VerifyCode() string {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return fmt.Sprintf("%d", rand.Intn(max-min+1)+min)
}
