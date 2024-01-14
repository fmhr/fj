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

type Config struct {
	Language           string   `toml:"Language"`
	Cmd                string   `toml:"Cmd"`
	Args               []string `toml:"Args"`
	Reactive           bool     `toml:"Reactive"`
	TesterPath         string   `toml:"TesterPath"`
	VisPath            string   `toml:"VisPath"`
	GenPath            string   `toml:"GenPath"`
	InfilePath         string   `toml:"InfilePath"`
	OutfilePath        string   `toml:"OutfilePath"`
	Jobs               int      `toml:"Jobs"`
	Cloud              bool     `toml:"Cloud"`
	CompilerURL        string   `toml:"CompilerURL"`
	Source             string   `toml:"Source"`
	CompileCmd         string   `toml:"CompileCmd"`
	Binary             string   `toml:"Binary"`
	TmpBinary          string   `toml:"tmpBinary"`
	Bucket             string   `toml:"Bucket"`
	WorkerURL          string   `toml:"WorkerURL"`
	ConcurrentRequests int      `toml:"ConcurrentRequests"`
	TimeLimitMS        uint64   `toml:"TimeLimitMS"`
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
	return conf, checkConfigFile(conf)
}

func checkConfigFile(cnf *Config) error {
	if cnf.Cmd == "" {
		return fmt.Errorf("cmd is empty. please set cmd in %s", CONFIG_FILE)
	}
	return nil
}

func configExists() bool {
	_, err := os.Stat("fj/" + CONFIG_FILE)
	return err == nil
}
