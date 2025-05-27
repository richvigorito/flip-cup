package utils

import (
	"fmt"
	"math/rand"
)

func RandID() string {
	return fmt.Sprintf("%x", rand.Intn(999999))
}
