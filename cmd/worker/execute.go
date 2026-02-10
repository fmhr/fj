package main

import (
	"fmt"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd"
	"github.com/fmhr/fj/cmd/setup"
)

func execute(config *setup.Config, seed int) (rtn *orderedmap.OrderedMap[string, string], err error) {
	// genコマンドで入力ファイルの生成
	err = cmd.Gen(config, seed)
	if err != nil {
		return nil, fmt.Errorf("genコマンドで入力ファイルの生成が失敗: %w", err)
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
		if err != nil {
			return nil, fmt.Errorf("ReactiveRunの実行に失敗: %w", err)
		}
		return rtn, nil
	} else {
		rtn, err := cmd.RunVis(config, seed)
		if err != nil {
			return nil, fmt.Errorf("RunVisの実行に失敗: %w", err)
		}
		return rtn, nil
	}
}
