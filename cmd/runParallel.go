package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
)

func RunParallel(cnf *config, seeds []int) {
	CORE := 4
	var wg sync.WaitGroup
	sem := make(chan struct{}, CORE)
	datas := make([]map[string]float64, 0, len(seeds))
	errorChan := make(chan string, len(seeds))

	var taskCompleted int32 = 0
	totalTask := len(seeds)

	// Ctrl+Cで中断したときに、現在実行中のseedを表示する
	var currentlyRunnningSeeds []int
	var seedMutex sync.Mutex
	var datasMutex sync.Mutex
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		<-sigCh
		seedMutex.Lock()
		log.Printf("\nCurrently running seeds: %v\n", currentlyRunnningSeeds)
		os.Exit(1)
	}()
	// -----

	for _, seed := range seeds {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			seedMutex.Lock()
			currentlyRunnningSeeds = append(currentlyRunnningSeeds, seed)
			seedMutex.Unlock()

			defer func() {
				//
				seedMutex.Lock()
				for i, s := range currentlyRunnningSeeds {
					if s == seed {
						currentlyRunnningSeeds = append(currentlyRunnningSeeds[:i], currentlyRunnningSeeds[i+1:]...)
						break
					}
				}
				seedMutex.Unlock()
				//
				<-sem
				wg.Done()
				atomic.AddInt32(&taskCompleted, 1)
				printProgress(int(taskCompleted), totalTask)
			}()
			out, err := Run(cnf, seed)
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
			datasMutex.Lock()
			datas = append(datas, data)
			datasMutex.Unlock()
		}(seed)
	}
	wg.Wait()
	close(errorChan)

	fmt.Fprintf(os.Stderr, "\n") // プログレスバーを改行

	for err := range errorChan {
		log.Println(err)
	}

	sumScore := 0.0
	for i := 0; i < len(datas); i++ {
		fmt.Println(datas[i])
		sumScore += datas[i]["score"]
	}
	fmt.Printf("sumScore=%.2f\n", sumScore)
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
	fmt.Fprintf(os.Stderr, "\r[%d/%d] [%s] %.2f%%", current, total, string(progressBar), percentage*100)
}
