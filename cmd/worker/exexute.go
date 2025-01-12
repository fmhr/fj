package main

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd"
	"github.com/fmhr/fj/cmd/setup"
)

func exexute(config *setup.Config, seed int) (rtn *orderedmap.OrderedMap[string, string], err error) {
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

func run(config *setup.Config, seed int) (*orderedmap.OrderedMap[string, string], error) {
	if config.Reactive {
		return cmd.ReactiveRun(config, seed)
	} else {
		return cmd.RunVis(config, seed)
	}
}
