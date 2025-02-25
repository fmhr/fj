package cmd

import (
	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
)

// RunSelector Cloud mode でなければ、reactiveRun か RunVis を呼び出す
func RunSelector(config *setup.Config, seed int) (*orderedmap.OrderedMap[string, string], error) {
	// Cloud mode なら、sendBinaryToWorker を呼び出す
	if config.CloudMode || (cloud != nil && *cloud) {
		rtn, err := requestToWorker(config, seed)
		if err != nil {
			return nil, err
		}
		return rtn, nil
	}
	// select run mode
	if config.Reactive {
		return ReactiveRun(config, seed)
	}
	return RunVis(config, seed)
}
