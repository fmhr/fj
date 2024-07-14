package fj

import (
	"fmt"
	"strings"
)

// Download tester file from URL.
func Download(url string) {
	// check if url is valid
	if !strings.HasPrefix(url, "http") {
		fmt.Println("Invalid URL")
		return
	}
}
