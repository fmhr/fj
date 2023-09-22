package fj

import (
	"fmt"
	"os"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

const configFileName = "fj_config.toml"

var ErrConfigNotFound = fmt.Errorf("%s not found, please run `fj -mode init`", configFileName)

type Config struct {
	Cmd         string `toml:"Cmd"`
	Reactive    bool   `toml:"Reactive"`
	TesterPath  string `toml:"TesterPath"`
	VisPath     string `toml:"VisPath"`
	GenPath     string `toml:"GenPath"`
	InfilePath  string `toml:"InfilePath"`
	OutfilePath string `toml:"OutfilePath"`
	Jobs        int    `toml:"Jobs"`
	CloudURL    string `toml:"CloudURL"`
}

func GenerateConfig() {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Printf("%s already exists\n", configFileName)
		return
	}
	numCPUs := maxInt(1, runtime.NumCPU()-1)
	conf := &Config{
		Cmd:         "",
		Reactive:    false,
		TesterPath:  "tools/target/release/tester",
		VisPath:     "tools/target/release/vis",
		GenPath:     "tools/target/release/gen",
		InfilePath:  "tools/in/",
		OutfilePath: "tmp/",
		Jobs:        numCPUs,
		CloudURL:    "http://localhost:8888",
	}
	if err := generateConfig(conf); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("%s is generated\n", configFileName)
}

func generateConfig(conf *Config) error {
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	return encoder.Encode(conf)
}

func LoadConfigFile() (*Config, error) {
	if !configExists() {
		return &Config{}, ErrConfigNotFound
	}

	conf := &Config{}
	file, err := os.Open(configFileName)
	if err != nil {
		return conf, err
	}
	defer file.Close()
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
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
	_, err := os.Stat(configFileName)
	return err == nil
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
