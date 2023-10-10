package main

import (
	"fmt"
	"os"

	"github.com/fmhr/fj"
)

func exexute(binaryPath string, language string, seed int) (rtn map[string]float64, err error) {
	var config fj.Config
	config.Cmd = "./" + binaryPath
	config.GenPath = "tools/tarrget/release/gen"
	config.VisPath = "tools/target/release/vis"
	config.TesterPath = "tools/target/release/tester"
	config.InfilePath = "in2/"
	_, err = os.Stat(config.TesterPath)
	if err != nil {
		config.Reactive = false
	} else {
		config.Reactive = true
	}
	// generate inputfile
	_, err = os.Stat(config.GenPath)
	if err != nil {
		return rtn, fmt.Errorf("gen file not found")
	}
	err = fj.Gen(&config, seed)
	if err != nil {
		return rtn, fmt.Errorf("failed to generate inputfile: %v", err)
	}
	// run
	if config.Reactive {
		return fj.ReactiveRun(&config, seed)
	} else {
		return fj.RunVis(&config, seed)
	}
}

func checkFileStat(config fj.Config) error {
	if _, err := os.Stat(config.GenPath); err != nil {
		return fmt.Errorf("gen file not found")
	}
	if _, err := os.Stat(config.TesterPath); err != nil {
		if _, err := os.Stat(config.VisPath); err != nil {
			return fmt.Errorf("tester file and vis file not found")
		}
	}
	return nil
}
