package download

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/net/html"
	"golang.org/x/net/publicsuffix"
)

func Logout() {
	dir, err := getConfigDir()
	if err != nil {
		log.Println("failed to get config dir:", err)
		return
	}
	cookieFile := filepath.Join(dir, "cookies.json")
	log.Println("remove cookies:", cookieFile)
	if err := os.Remove(cookieFile); err != nil {
		log.Println("failed to remove cookies:", err)
		return
	}
	fmt.Println("success")
}

func Login(url string, username string, password string) error {
	// すでにログインしているかどうか
	client := NewAtCoderClient(username, password)
	// ログインしていない場合はログイン
	if err := client.Login(); err != nil {
		return err
	}
	fmt.Println("success")
	return nil
}

const (
	baseURL   = "https://atcoder.jp"
	loginURL  = "https://atcoder.jp/login"
	submitURL = "https://atcoder.jp/contests/agc001/submit"
)

type AtCoderClient struct {
	Client   *http.Client
	Username string
	Password string
}

func NewAtCoderClient(username string, password string) *AtCoderClient {
	jar, err := loadCookieds()
	if err != nil {
		log.Println("failed to load cookies:", err)
		jar, _ = cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	}
	return &AtCoderClient{
		Client: &http.Client{
			Jar: jar,
		},
		Username: username,
		Password: password,
	}
}

// submitのページにアクセスして、ログインしているかどうかを確認する
// 200: ログイン済み
// 302: ログインしていない
func (c *AtCoderClient) IsLoggedIn() bool {
	// リダイレクトを一時的に無効
	// 関数終了時に元に戻す
	tempCheckRedirect := c.Client.CheckRedirect
	defer func() { c.Client.CheckRedirect = tempCheckRedirect }()
	// リダイレクトを検知するための、カスタムトランスポートを設定
	c.Client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	resp, err := c.Client.Get(submitURL)
	if err != nil {
		fmt.Printf("failed to access %s:%v", submitURL, err)
		return false
	}
	defer resp.Body.Close()

	// 200: ログイン済み
	if resp.StatusCode == http.StatusOK {
		return true
	}
	// リダイレクトされた場合は、ログインしていないということなのでfalseを返す
	if resp.StatusCode == http.StatusFound {
		return false
	}
	fmt.Println("unexpected status code:", resp.StatusCode)
	return false
}

func (c *AtCoderClient) Login() error {
	if c.IsLoggedIn() {
		fmt.Println("already logged in.")
		return nil
	}
	// ログインページからCSRFトークンを取得
	token, err := c.GetCSRFToken()
	if err != nil {
		return fmt.Errorf("failed to get csrf_token:%v", err)
	}
	// ログイン処理
	formData := url.Values{
		"username":   {c.Username},
		"password":   {c.Password},
		"csrf_token": {token},
	}
	req, err := http.NewRequest("POST", loginURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create request:%v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "fmhr-judge-tools/0.0.1")

	resp, err := c.Client.Do(req)
	if err != nil {
		return fmt.Errorf("unexpected error:%v", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to login:%d", resp.StatusCode)
	}
	if !c.IsLoggedIn() {
		return fmt.Errorf("failed to login. please check your username and password")
	}
	return saveCookies(c.Client.Jar.(*cookiejar.Jar))
}

func (c *AtCoderClient) GetCSRFToken() (string, error) {
	resp, err := c.Client.Get(loginURL)
	if err != nil {
		return "", fmt.Errorf("failed to access %s:%v", loginURL, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read body:%v", err)
	}
	return getCSRFToken(string(body))
}

func getCSRFToken(body string) (string, error) {
	doc, err := html.Parse(strings.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("failed to parse html:%v", err)
	}
	var token string
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "input" {
			for _, attr := range n.Attr {
				if attr.Key == "name" && attr.Val == "csrf_token" {
					for _, attr := range n.Attr {
						if attr.Key == "value" {
							token = attr.Val
							return
						}
					}
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(doc)
	if token == "" {
		return "", fmt.Errorf("csrf_token not found")
	}
	return token, nil
}

// get config dir はcookieを保存するディレクトリを決める
func getConfigDir() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user config dir:%v", err)
	}
	dir = filepath.Join(dir, "fmhr-judge-tools")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create dir:%v", err)
	}
	return dir, nil
}

func loadCookieds() (*cookiejar.Jar, error) {
	jar, _ := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})

	dir, err := getConfigDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get config dir:%v", err)
	}
	cookieFile := filepath.Join(dir, "cookies.json")
	// 存在しない場合は、jarをそのまま返す
	if _, err := os.Stat(cookieFile); err != nil {
		return jar, nil
	}
	file, err := os.Open(cookieFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var cookies []*http.Cookie
	if err := json.NewDecoder(file).Decode(&cookies); err != nil {
		return nil, err
	}
	baseURL, _ := url.Parse("https://atcoder.jp")
	jar.SetCookies(baseURL, cookies)
	return jar, nil
}

// saveCookies はログイン成功後にcookieを保存する
func saveCookies(jar *cookiejar.Jar) error {
	dir, err := getConfigDir()
	if err != nil {
		return fmt.Errorf("failed to get config dir:%v", err)
	}
	cookieFile := filepath.Join(dir, "cookies.json")
	file, err := os.Create(cookieFile)
	if err != nil {
		return fmt.Errorf("failed to create file:%v", err)
	}
	defer file.Close()
	log.Println("save cookies:", cookieFile)
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	baseURL, _ := url.Parse("https://atcoder.jp")
	return encoder.Encode(jar.Cookies(baseURL))
}
