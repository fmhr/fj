package main

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd"
)

func exexute(config *fj.Config, seed int) (rtn *orderedmap.OrderedMap[string, any], err error) {
	// GEN
	err = fj.Gen(config, seed)
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

func run(config *fj.Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	if config.Reactive {
		return fj.ReactiveRun(config, seed)
	} else {
		return fj.RunVis(config, seed)
	}
}
