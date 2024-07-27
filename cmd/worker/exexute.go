package main

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd"
)

func exexute(config *cmd.Config, seed int) (rtn *orderedmap.OrderedMap[string, any], err error) {
	// GEN
	err = cmd.Gen(config, seed)
	if err != nil {
		return nil, err
	}

	// RUN
	rtn, err = run(config, seed)
	if err != nil {
		return nil, err
	}
	return rtn, nil
}

func run(config *cmd.Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	if config.Reactive {
		return cmd.ReactiveRun(config, seed)
	} else {
		return cmd.RunVis(config, seed)
	}
}
