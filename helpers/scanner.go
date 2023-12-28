package helpers

import (
	"bufio"
	"os"

	"github.com/fatih/color"
)

var printInfo = color.New(color.FgGreen, color.Bold)
var printErr = color.New(color.FgRed)

func GetInput(placeholder string) string {
	var scanner = bufio.NewScanner(os.Stdin)
	var text string
	for {
		printInfo.Printf(placeholder)
		scanner.Scan()
		text = scanner.Text()
		if len(text) > 0 {
			break
		}
	}
	return text
}
