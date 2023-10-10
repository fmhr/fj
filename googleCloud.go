package fj

// すでにバイナリがある状態にしておく
func requestGoogleCloud(config *Config, seed int) (map[string]float64, error) {
	return sendBinaryToWorker(config.WorkerURL, config.BinaryPath, config.Language, seed)
}
