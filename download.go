package fj

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"strings"
)

// Download tester file from URL.
func Download(url string) {
	// check if url is valid
	if !strings.HasPrefix(url, "http") {
		fmt.Println("Invalid URL")
		return
	}
	ac := NewAtCoderClient("", "")
	// login to atcoder
	if err := ac.Login(); err != nil {
		fmt.Println("Failed to login")
		return
	}
	if err := DownloadLoacaTesterZip(ac.Client, url); err != nil {
		fmt.Println("Failed to download loacatester.zip")
		return
	}
	fmt.Println("Downloaded loacatester.zip")
}

func DownloadLoacaTesterZip(client *http.Client, url string) error {
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("failed to access %s:%v", url, err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read body:%v", err)
	}

	// リンクを抽出するための正規表現
	re := regexp.MustCompile(`<a href="(https://img\.atcoder\.jp/[^"]+\.zip)">(ローカル版|Local version)</a>`)
	match := re.FindSubmatch(body)
	if match == nil {
		return fmt.Errorf("failed to find download link")
	}

	zipURL := string(match[1])
	fmt.Printf("Found zip file URL: %s\n", zipURL)

	// download zip file
	resp, err = client.Get(zipURL)
	if err != nil {
		return fmt.Errorf("failed to download zip file:%v", err)
	}
	defer resp.Body.Close()

	zipBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read zip body:%v", err)
	}

	// save zip file
	err = os.WriteFile("loacatester.zip", zipBody, 0644)
	if err != nil {
		return fmt.Errorf("failed to save zip file:%v", err)
	}
	fmt.Println("Downloaded zip file")

	return nil
}
