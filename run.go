package fj

func Run(config *Config, seed int) (map[string]float64, error) {
	if config.Reactive {
		return ReactiveRun(config, seed)
	} else {
		return RunVis(config, seed)
	}
}
