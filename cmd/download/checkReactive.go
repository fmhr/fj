package download

import (
	"log"
	"os"
	"strings"
)

// README.mdにtesterがあるかを確認する
func IsReactive() bool {
	// read tools/README.md
	file, err := os.ReadFile("tools/README.md")
	if err != nil {
		log.Println("Failed to read README.md, ", err)
		return false
	}
	// check if "tester" is in the file
	if strings.Contains(string(file), "tester") {
		return true
	}
	return false
}
