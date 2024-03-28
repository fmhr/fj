package fj

import (
	"embed"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml/v2"
)

//go:embed config.toml
var configContent embed.FS

const (
	CONFIG_FILE  = "config.toml"
	FJ_DIRECTORY = "fj/"
)

var ErrConfigNotFound = fmt.Errorf("%s not found, please run `fj init`", CONFIG_FILE)

// Config はfjの設定を保持します
// ソースファイルのパス、バイナリーパスが不正確だと、コンパイルコマンド、実行コマンドが正常に動作しない
type Config struct {
	Language           string   `toml:"Language"`
	Args               []string `toml:"Args"`               // 実行時の引数
	Contest            string   `toml:"Contest"`            // コンテスト名 例: abc123 TODO
	Reactive           bool     `toml:"Reactive"`           // 問題の種類
	TimeLimitMS        uint64   `toml:"TimeLimitMS"`        // 問題の制限時間
	TesterPath         string   `toml:"TesterPath"`         // リアクティブ問題のテスター
	VisPath            string   `toml:"VisPath"`            // ノンリアクティブ問題の可視化プログラム(Score計算用)
	GenPath            string   `toml:"GenPath"`            // 問題生成プログラム サーバー上で使う
	InfilePath         string   `toml:"InfilePath"`         // ノンリアクティブ問題の入力ファイル
	OutfilePath        string   `toml:"OutfilePath"`        // ノンリアクティブ問題の出力ファイル
	Jobs               int      `toml:"Jobs"`               // ローカル実行時の並列実行数
	BinaryPath         string   `toml:"BinaryPath"`         // バイナリの保存先
	CloudMode          bool     `toml:"Cloud"`              // デフォルトの実行モード
	CompilerURL        string   `toml:"CompilerURL"`        // クラウド上のコンパイラのURL
	SourceFilePath     string   `toml:"Source"`             // クラウドにアップロードするソースファイル
	TmpBinary          string   `toml:"tmpBinary"`          // クラウドにアップロードするバイナリファイルのランダム生成された名前
	Bucket             string   `toml:"Bucket"`             // バイナリの保存先
	WorkerURL          string   `toml:"WorkerURL"`          // クラウドワーカーのURL ここが多数立ち上がる
	ConcurrentRequests int      `toml:"ConcurrentRequests"` // クラウドワーカーの並列アクセス数
	CustomExeCmd       string   `toml:"CustomExeCmd"`       // 実行コマンドの上書き TODO
	CustomCompileCmd   string   `toml:"CustomCompileCmd"`   // コンパイルコマンドの上書き TODO
}

// generateConfig はfj init コマンドで読み出されて fj/config.tomlを生成する
func generateConfig() error {
	// fj ディレクトリを作成
	if _, err := os.Stat(FJ_DIRECTORY); os.IsNotExist(err) {
		if err := os.Mkdir(FJ_DIRECTORY, 0755); err != nil {
			return err
		}
	}

	// config.tomlの内容を所得
	configBytes, err := configContent.ReadFile("config.toml")
	if err != nil {
		return err
	}

	// config.tomlを作成
	filePath := filepath.Join(FJ_DIRECTORY, CONFIG_FILE)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// config.tomlに書き込む
	if _, err := file.Write(configBytes); err != nil {
		return err
	}

	return nil
}

// setConfig はconfig.tomlを読み込む
func setConfig() (*Config, error) {
	if !configExists() {
		return &Config{}, ErrConfigNotFound
	}
	conf := &Config{}
	file, err := os.Open(FJ_DIRECTORY + CONFIG_FILE)
	if err != nil {
		return conf, err
	}
	defer file.Close()
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		log.Println("failed to decode config file: ", err)
		return conf, err
	}
	return conf, nil
}

func configExists() bool {
	_, err := os.Stat("fj/" + CONFIG_FILE)
	return err == nil
}
