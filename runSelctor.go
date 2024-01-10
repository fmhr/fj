package fj

import "github.com/elliotchance/orderedmap/v2"

// RunSelector Cloud mode でなければ、reactiveRun か RunVis を呼び出す
func RunSelector(config *Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	// Cloud mode なら、sendBinaryToWorker を呼び出す
	if config.Cloud || (cloud != nil && *cloud) {
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
