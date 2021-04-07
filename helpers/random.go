package helpers

import (
	"math/rand"
	"strings"
	"time"
)

// RandomString returns a (pseudo) string of specified length.
func RandomString(length int) string {
	chars := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := strings.Builder{}

	for i := 0; i < length; i++ {
		r := rand.Intn(len(chars))
		c := chars[r]
		str.WriteByte(c)
	}

	return str.String()
}

func init() {
	rand.Seed(time.Now().UnixNano())
}
