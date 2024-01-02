package helpers

import (
	"fmt"
	"math/rand"
)

func Spr(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func randomBytes(length int) []byte {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b = make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
