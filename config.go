package fj

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const configFileName = "fj_config.toml"

var ErrConfigNotFound = fmt.Errorf("%s not found, please run `fj init`", configFileName)

type Config struct {
	Cmd         string   `toml:"Cmd"`
	Args        []string `toml:"Args"`
	Reactive    bool     `toml:"Reactive"`
	TesterPath  string   `toml:"TesterPath"`
	VisPath     string   `toml:"VisPath"`
	GenPath     string   `toml:"GenPath"`
	InfilePath  string   `toml:"InfilePath"`
	OutfilePath string   `toml:"OutfilePath"`
	Jobs        int      `toml:"Jobs"`
	Cloud       bool     `toml:"Cloud"`
	CloudURL    string   `toml:"CloudURL"`
	TimeLimitMS uint64   `toml:"TimeLimitMS"`
}

func GenerateConfig() {
	if _, err := os.Stat(configFileName); err == nil {
		if *force {
			// if force flag is set, remove config file
			err := os.Remove(configFileName)
			if err != nil {
				fmt.Println("Failed to remove config file: ", err)
				return
			}
		} else {
			fmt.Printf("%s already exists\n", configFileName)
			return
		}
	}
	//numCPUs := maxInt(1, runtime.NumCPU()-1)
	conf := &Config{
		Cmd:         "",
		Args:        []string{},
		Reactive:    false,
		TesterPath:  "tools/target/release/tester",
		VisPath:     "tools/target/release/vis",
		GenPath:     "tools/target/release/gen",
		InfilePath:  "tools/in/",
		OutfilePath: "out/",
		Jobs:        1,
		Cloud:       false,
		CloudURL:    "http://localhost:8888",
		TimeLimitMS: 20000,
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
