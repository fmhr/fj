package main

import (
	"fmt"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd"
	"github.com/fmhr/fj/cmd/setup"
)

func execute(config *setup.Config, seed int) (rtn *orderedmap.OrderedMap[string, string], err error) {
	// コマンド生成
	err = cmd.Gen(config, seed)
	if err != nil {
		return nil, fmt.Errorf("コマンド生成が失敗: %w", err)
	}

	// RUN
	rtn, err = run(config, seed)
	if err != nil {
		return nil, fmt.Errorf("RUNが失敗: %w", err)
	}
	return rtn, nil
}

func run(config *setup.Config, seed int) (*orderedmap.OrderedMap[string, string], error) {
	if config.Reactive {
		rtn, err := cmd.ReactiveRun(config, seed)
		return rtn, fmt.Errorf("ReactiveRunの実行に失敗: %w", err)
	} else {
		rtn, err := cmd.RunVis(config, seed)
		return rtn, fmt.Errorf("RunVisの実行に失敗: %w", err)
	}
}
