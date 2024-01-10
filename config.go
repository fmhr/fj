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
	configFileName = "config.toml"
	directory      = "fj"
)

var ErrConfigNotFound = fmt.Errorf("%s not found, please run `fj init`", configFileName)

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

func generateConfig() error {
	// fj ディレクトリを作成
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		if err := os.Mkdir(directory, 0755); err != nil {
			return err
		}
	}

	// config.tomlの内容を所得
	configBytes, err := configContent.ReadFile("config.toml")
	if err != nil {
		return err
	}

	// config.tomlを作成
	filePath := filepath.Join(directory, configFileName)
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

func LoadConfigFile() (*Config, error) {
	if !configExists() {
		return &Config{}, ErrConfigNotFound
	}
	conf := &Config{}
	file, err := os.Open("fj/" + configFileName)
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
		return fmt.Errorf("cmd is empty. please set cmd in %s", configFileName)
	}
	return nil
}

func configExists() bool {
	_, err := os.Stat("fj/" + configFileName)
	return err == nil
}
