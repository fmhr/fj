package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj"
)

func exexute(config *fj.Config, seed int) (rtn *orderedmap.OrderedMap[string, any], err error) {
	// GEN
	err = fj.Gen(config, seed)
	if err != nil {
		return nil, err
	}
	// javaの場合はコンパイル
	if config.Language == "java" {
		compileCmd := fj.LanguageSets[config.Language].CompileCmd
		cmds := strings.Fields(compileCmd)
		cmd := exec.Command(cmds[0], cmds[1:]...)
		msg, err := cmd.CombinedOutput()
		if err != nil {
			log.Println(msg)
			err := fmt.Errorf("failed to compile: [%s]%v msg: %s", cmd.String(), err, string(msg))
			return nil, err
		}
	}
	// RUN
	rtn, err = run(config, seed)
	if err != nil {
		return nil, err
	}
	return rtn, nil
}

func run(config *fj.Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	if config.Reactive {
		return fj.ReactiveRun(config, seed)
	} else {
		return fj.RunVis(config, seed)
	}
}
