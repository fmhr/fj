package fj

import (
	"fmt"
	"log"
)

// --clod のときはこちらの処理に移行する
// test,testsどちらのときもこちらであることに注意

func GoogleCloudMode(config *Config, seeds []int) {
	requestGoogleCloud(config, seeds)
}

func requestGoogleCloud(config *Config, seeds []int) error {
	binaryPath, err := uploadAndReceive(config.CompilerURL, config.SourcePath, config.Language)
	if err != nil {
		return fmt.Errorf("failed to upload and receive: %v", err)
	}
	outs := make([]string, len(seeds))
	for _, seed := range seeds {
		go func(i int) {
			out, err := sendBinaryToWorker(config.WorkerURL, binaryPath, config.Language, i)
			if err != nil {
				log.Printf("failed to send binary to worker by seed:%d %v", i, err)
			}
			outs[i] = out
		}(seed)
	}
	return nil
}
