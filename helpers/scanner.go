package helpers

import (
	"bufio"
	"os"
)

func GetInput(placeholder string) string {
	var scanner = bufio.NewScanner(os.Stdin)
	var text string
	for {
		LogInfo.Printf(placeholder)
		scanner.Scan()
		text = scanner.Text()
		if len(text) > 0 {
			break
		}
	}
	return text
}
