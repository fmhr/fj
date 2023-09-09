package main

import (
	"fmt"
	"log"
	"sync"
	"sync/atomic"
)

func RunParallel(seeds []int) {
	CORE := 4
	var wg sync.WaitGroup
	sem := make(chan struct{}, CORE)
	datas := make([]map[string]float64, len(seeds))

	var taskCompleted int32 = 0
	totalTask := len(seeds)

	errorChan := make(chan string, len(seeds))

	for _, seed := range seeds {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			defer func() {
				<-sem
				wg.Done()
				atomic.AddInt32(&taskCompleted, 1)
				printProgress(int(taskCompleted), totalTask)
			}()
			out, err := Run(seed)
			if err != nil {
				errorChan <- fmt.Sprintf("seed=%d %v\n", seed, err)
				return
			}
			data, err := ExtractKeyValuePairs(string(out))
			if err != nil {
				errorChan <- fmt.Sprintf("seed=%d %v\n", seed, err)
				return
			}
			data["seed"] = float64(seed)
			datas[seed] = data
		}(seed)
	}
	wg.Wait()
	close(errorChan)

	fmt.Println() // プログレスバーを改行

	for err := range errorChan {
		log.Println(err)
	}

	sumScore := 0.0
	for i := 0; i < len(datas); i++ {
		log.Println(datas[i])
		sumScore += datas[i]["score"]
	}
	log.Printf("sumScore=%.2f\n", sumScore)
	fmt.Printf("%.2f\n", sumScore)
}

const progressBarWidth = 40

func printProgress(current, total int) {
	percentage := float64(current) / float64(total)
	barLength := int(percentage * float64(progressBarWidth))
	progressBar := make([]rune, progressBarWidth)

	for i := 0; i < progressBarWidth; i++ {
		if i < barLength {
			progressBar[i] = '■'
		} else {
			progressBar[i] = ' '
		}
	}
	fmt.Printf("\r[%d/%d] [%s] %.2f%%", current, total, string(progressBar), percentage*100)
}
