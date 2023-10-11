package fj

func RunSelctor(config *Config, seed int) (map[string]float64, error) {
	if config.Cloud || *cloud {
		// return CloudRun(config, seed)
		return SendBinaryToWorker(config, seed)
	} else {
		if config.Reactive {
			return ReactiveRun(config, seed)
		} else {
			return RunVis(config, seed)
		}
	}
}
