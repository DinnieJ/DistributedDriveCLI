package helpers

import "fmt"

func Spr(format string, a ...any) string {
	return fmt.Sprintf(format, a...)
}
