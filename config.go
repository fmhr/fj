package fj

import (
	"fmt"
	"log"
	"os"
	"runtime"

	"github.com/pelletier/go-toml/v2"
)

const (
	configFileName = "fj-config.toml"
)

var ErrConfigNotFound = fmt.Errorf("%s not found, please run `fj -mode init`", configFileName)

type config struct {
	Cmd         string `toml:"Cmd"`
	Reactive    bool   `toml:"Reactive"`
	TesterPath  string `toml:"TesterPath"`
	VisPath     string `toml:"VisPath"`
	InfilePath  string `toml:"InfilePath"`
	OutfilePath string `toml:"OutfilePath"`
	Jobs        int    `toml:"Jobs"`
}

func GenerateConfig() {
	if _, err := os.Stat(configFileName); err == nil {
		fmt.Printf("%s already exists\n", configFileName)
		return
	}
	numCPUs := runtime.NumCPU() - 1
	conf := &config{
		TesterPath:  "tools/target/release/tester",
		VisPath:     "tools/target/release/vis",
		InfilePath:  "tools/in/",
		OutfilePath: "tmp/",
		Jobs:        numCPUs,
	}
	if err := generateConfig(conf); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Printf("%s is generated\n", configFileName)
}

func generateConfig(conf *config) error {
	file, err := os.Create(configFileName)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := toml.NewEncoder(file)
	err = encoder.Encode(conf)
	if err != nil {
		return err
	}
	return nil
}

func LoadConfigFile() (*config, error) {
	if _, err := os.Stat(configFileName); err != nil {
		log.Printf("%s not found, please run `fj -mode init`\n", configFileName)
		return &config{}, err
	}
	conf := &config{}
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

func checkConfigFile(cnf *config) error {
	if cnf.Cmd == "" {
		return fmt.Errorf("cmd is empty. please set cmd in %s", configFileName)
	}
	return nil
}
