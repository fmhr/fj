package setup

import (
	"fmt"
	"log"
	"os"

	"github.com/fmhr/fj/cmd/download"
	"github.com/pelletier/go-toml/v2"
)

const (
	CONFIG_FILE  = "config.toml"
	FJ_DIRECTORY = "fj/"
)

var ErrConfigNotFound = fmt.Errorf("%s not found, please run `fj init`", CONFIG_FILE)

// Config はfjの設定を保持します
// ソースファイルのパス、バイナリーパスが不正確だと、コンパイルコマンド、実行コマンドが正常に動作しない
type Config struct {
	Language           string   `toml:"Language"`
	ExecuteCmd         string   `toml:"ExecuteCmd"`         // 実行コマンド
	Args               []string `toml:"Args"`               // 実行コマンドのオプション
	Contest            string   `toml:"Contest"`            // コンテスト名
	Reactive           bool     `toml:"Reactive"`           // 問題の種類
	TimeLimitMS        uint64   `toml:"TimeLimitMS"`        // 問題の制限時間　または強制終了するまでの時間
	TesterPath         string   `toml:"TesterPath"`         // リアクティブ問題のテスターのパス
	VisPath            string   `toml:"VisPath"`            // ノンリアクティブ問題の可視化プログラム(Score計算用)
	GenPath            string   `toml:"GenPath"`            // 問題生成プログラム サーバー上で使う
	InfilePath         string   `toml:"InfilePath"`         // 問題の入力ファイル
	OutfilePath        string   `toml:"OutfilePath"`        // 問題の出力ファイル
	Jobs               int      `toml:"Jobs"`               // ローカル実行時の並列実行数
	BinaryPath         string   `toml:"BinaryPath"`         // コンテナ内のバイナリの保存先
	CloudMode          bool     `toml:"Cloud"`              // デフォルトの実行モード
	CompilerURL        string   `toml:"CompilerURL"`        // クラウド上のコンパイラのURL
	SourceFilePath     string   `toml:"Source"`             // クラウドにアップロードするソースファイル
	TmpBinary          string   `toml:"tmpBinary"`          // クラウドにアップロードするバイナリファイルのランダム生成された名前
	Bucket             string   `toml:"Bucket"`             // バイナリの保存先
	WorkerURL          string   `toml:"WorkerURL"`          // クラウドジャッジコンテナのURL ここが多数立ち上がる
	ConcurrentRequests int      `toml:"ConcurrentRequests"` // クラウドジャッジコンテナの並列アクセス数
}

// GenerateConfig はfj init コマンドで読み出されて fj/config.tomlを生成する
func GenerateConfig() error {
	// fj ディレクトリを作成
	if _, err := os.Stat(FJ_DIRECTORY); os.IsNotExist(err) {
		if err := os.Mkdir(FJ_DIRECTORY, 0755); err != nil {
			return err
		}
	}

	return nil
}

// SetConfig はconfig.tomlを読み込む
func SetConfig() (*Config, error) {
	if !configExists() {
		// configファイルがない場合(*これがv2のデフォルト)
		// 最小限にする
		conf := newConfig()
		return &conf, nil
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
	overWrite(conf)
	return conf, nil
}

func overWrite(config *Config) error {
	if config.BinaryPath == "" {
		config.BinaryPath = LanguageSets[config.Language].BinaryPath
	}
	return nil
}

func configExists() bool {
	_, err := os.Stat("fj/" + CONFIG_FILE)
	return err == nil
}

func newConfig() (c Config) {
	c.Reactive = download.IsReactive()
	c.TesterPath = "tools/target/release/tester"
	c.VisPath = "tools/target/release/vis"
	c.GenPath = "tools/target/release/gen"
	c.InfilePath = "tools/in/"
	c.OutfilePath = "out/"
	c.TimeLimitMS = 5000
	return c
}
