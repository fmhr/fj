package fj

import (
	"fmt"
	"log"
	"os"

	"github.com/pelletier/go-toml"
)

var ErrConfigNotFound = fmt.Errorf("config.toml not found, please run `fj init`")

type config struct {
	Cmd         string `toml:"Cmd"`
	Reactive    bool   `toml:"Reactive"`
	ToolsPAHT   string `toml:"ToolsPAHT"`
	visPath     string `toml:"visPath"`
	infilePath  string `toml:"infilePath"`
	outfilePath string `toml:"outfilePath"`
	Jobs        int    `toml:"Jobs"`
}

func GenerateConfig() {
	if _, err := os.Stat("config.toml"); err == nil {
		fmt.Println("config.toml already exists")
		return
	}
	conf := &config{
		visPath:     "./tools/target/release/vis",
		infilePath:  "./tools/in/",
		outfilePath: "./tmp/",
	}
	if err := generateConfig(conf); err != nil {
		fmt.Println("Error: ", err)
		return
	}
	fmt.Println("config.toml is generated")
}

func generateConfig(conf *config) error {
	file, err := os.Create("config.toml")
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
	if _, err := os.Stat("config.toml"); err != nil {
		log.Println("config.toml not found, please run `fj init`")
		return &config{}, err
	}
	conf := &config{}
	file, err := os.Open("config.toml")
	if err != nil {
		return conf, err
	}
	defer file.Close()
	decoder := toml.NewDecoder(file)
	err = decoder.Decode(conf)
	if err != nil {
		return conf, err
	}
	return conf, nil
}
