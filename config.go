package fj

import (
	"fmt"
	"os"

	"github.com/pelletier/go-toml/v2"
)

const (
	configFileName = "config.toml"
	directory      = "fj/"
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
	CloudURL           string   `toml:"CloudURL"`
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

func newConfig() *Config {
	return &Config{
		Language:           "Go",
		Cmd:                "./bin/main",
		Args:               []string{},
		Reactive:           false,
		TesterPath:         "tools/target/release/tester",
		VisPath:            "tools/target/release/vis",
		GenPath:            "tools/target/release/gen",
		InfilePath:         "tools/in/",
		OutfilePath:        "out/",
		Jobs:               4,
		Cloud:              false,
		CloudURL:           "http://localhost:8888",
		CompilerURL:        "http://localhost:8080/compiler",
		Bucket:             "",
		CompileCmd:         "go build -o bin/main src/main.go",
		Source:             "src/main.go",
		Binary:             "bin/main",
		WorkerURL:          "http://localhost:8081/worker",
		ConcurrentRequests: 1000,
		TimeLimitMS:        10000,
	}
}

func GenerateConfig() {
	if _, err := os.Stat(directory + configFileName); err == nil {
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
	conf := newConfig()

	if err := generateConfig(conf); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("%s is generated\n", configFileName)
}

func generateConfig(conf *Config) error {
	// create fj directory
	if _, err := os.Stat("fj"); err != nil {
		if err := os.Mkdir("fj", 0777); err != nil {
			return fmt.Errorf("failed to create fj directory: %v", err)
		}
	}
	// create config file
	file, err := os.Create(directory + configFileName)
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
	file, err := os.Open("fj/" + configFileName)
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
	_, err := os.Stat("fj/" + configFileName)
	return err == nil
}
