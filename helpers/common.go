package helpers

import (
	"fmt"
	"math/rand"
	"os"
)

func Spr(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}

func If[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

func Must[T any](value T, err error) T {
	if err != nil {
		LogErr.Printf("[ERROR]: %s\n", err.Error())
		os.Exit(1)
	}
	return value
}

func Contain[T comparable](l []T, haystack T) bool {
	for _, v := range l {
		if v == haystack {
			return true
		}
	}
	return false
}

func randomBytes(length int) []byte {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var b = make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return b
}
