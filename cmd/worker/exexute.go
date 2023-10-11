package main

import (
	"github.com/fmhr/fj"
)

func exexute(config *fj.Config, seed int) (rtn map[string]float64, err error) {
	// GEN
	err = fj.Gen(config, seed)
	if err != nil {
		return nil, fj.TraceErrorf("failed to gen: %v", err)
	}
	// RUN
	rtn, err = run(config, seed)
	if err != nil {
		return nil, fj.TraceErrorf("failed to run: %v", err)
	}
	return rtn, nil
}

func run(config *fj.Config, seed int) (map[string]float64, error) {
	if config.Reactive {
		return fj.ReactiveRun(config, seed)
	} else {
		return fj.RunVis(config, seed)
	}
}
