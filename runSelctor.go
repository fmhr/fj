package fj

func RunSelctor(config *Config, seed int) (map[string]float64, error) {
	if !config.Cloud && !*cloud {
		if config.Reactive {
			return ReactiveRun(config, seed)
		} else {
			return RunVis(config, seed)
		}
	} else {
		// Cloud mode
		rtn, err := sendBinaryToWorker(config, seed)
		if err != nil {
			return nil, TraceErrorf("failed to run: %v", err)
		}
		return rtn, nil
	}
}
