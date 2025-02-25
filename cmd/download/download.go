package download

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"time"
)

// Download tester file from URL.
func Download(urlStr string) error {
	// check if url is valid
	u, err := url.Parse(urlStr)
	if err != nil || u.Scheme == "" || u.Host == "" {
		return fmt.Errorf("invalid URL: %s", urlStr)
	}
	u.RawQuery = ""
	urlStr = u.String()

	ac := NewAtCoderClient("", "")
	// login to atcoder
	if !ac.IsLoggedIn() {
		return fmt.Errorf("not logged in")
	}
	if err := DownloadLoacaTesterZip(ac.Client, urlStr); err != nil {
		return fmt.Errorf("failed to download loacatester.zip:%v", err)
	}
	fmt.Println("[SUCCESS]Download loacatester.zip")
	defer func() {
		if err := os.Remove("loacatester.zip"); err != nil {
			log.Println("Failed to remove loacatester.zip")
		}
	}()
	if err := unzip("loacatester.zip", ""); err != nil {
		return fmt.Errorf("failed to unzip loacatester.zip:%v", err)
	}
	fmt.Println("Unzipped loacatester.zip")
	return downloadLogging(urlStr) // ダウンロード履歴をログに記録
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
	re := regexp.MustCompile(`<a href="(https://img\.atcoder\.jp/[^"]+\.zip(?:\?[^"]*)?)">.*?(ローカル版|Local version).*?</a>`)

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
	return nil
}

type DownloadLogging struct {
	TimeStamp time.Time `json:"time_stamp"`
	Directory string    `json:"directory"`
	Url       string    `json:"url"`
	Reactive  bool      `json:"reactive"`
}

func downloadLogging(url string) error {
	dlog := DownloadLogging{
		TimeStamp: time.Now(),
		Directory: "",
		Url:       url,
	}
	// カレントディレクトリを取得
	dir, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get current directory: %v", err)
	}
	dlog.Directory = dir

	appName := "fmhr-judge-tools"
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		return fmt.Errorf("failed to get user cache directory: %v", err)
	}
	// キャッシュディレクトリを作成
	cacheDir = filepath.Join(cacheDir, appName)
	if err := os.MkdirAll(cacheDir, 0755); err != nil {
		return fmt.Errorf("failed to create cache directory: %v", err)
	}
	// jsonファイルのパス
	jsonFile := filepath.Join(cacheDir, "download-history.json")

	// 既存のダウンロード履歴を読み込む
	var dlogs []DownloadLogging
	if _, err := os.Stat(jsonFile); err == nil {
		// ファイルがすでにあるとき
		file, err := os.ReadFile(jsonFile)
		if err != nil {
			return fmt.Errorf("failed to read download history: %v", err)
		}
		// dlogsにfileの内容を読み込む
		if err := json.Unmarshal(file, &dlogs); err != nil {
			return fmt.Errorf("failed to unmarshal download history: %v", err)
		}
	}
	// isReactiveはtools/README.mdを読んで確認する
	reactibe := IsReactive()
	dlog.Reactive = reactibe
	if reactibe {
		fmt.Println("Reactive Problem")
	} else {
		fmt.Println("Not Reactive Problem")
	}

	// dlogsにdlogを追加
	dlogs = append(dlogs, dlog)
	// 重複を消す
	dlogs = removeOldDownloadLogs(dlogs)

	jsonData, err := json.MarshalIndent(dlogs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal download logging: %v", err)
	}
	// ファイルに書き込む
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return fmt.Errorf("failed to write download logging to json file: %v", err)
	}
	fmt.Println("Wrote download logging to json file")

	// tools内のgen, vis, tester をbuild
	return buildTools()
}

// removeOldDownloadLogsは、dlogsを並べ替えて、ディレクトリが同じものを探して、最も新しい記録以外を消す
func removeOldDownloadLogs(dlogs []DownloadLogging) []DownloadLogging {
	sort.Slice(dlogs, func(i, j int) bool {
		return dlogs[i].TimeStamp.After(dlogs[j].TimeStamp)
	})
	newDlogs := make([]DownloadLogging, 0)
	directoryMap := make(map[string]bool)
	for _, dlog := range dlogs {
		if _, ok := directoryMap[dlog.Directory]; !ok {
			newDlogs = append(newDlogs, dlog)
			directoryMap[dlog.Directory] = true
		}
	}
	return newDlogs
}
