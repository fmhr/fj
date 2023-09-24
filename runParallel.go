package fj

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
)

func RunParallel(cnf *Config, seeds []int) {
	// 並列実行数の設定
	numCPUs := runtime.NumCPU() - 1
	if cnf.Jobs > 0 {
		numCPUs = cnf.Jobs
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, numCPUs)
	datas := make([]map[string]float64, 0, len(seeds))
	errorChan := make(chan string, len(seeds))

	var taskCompleted int32 = 0
	totalTask := len(seeds)

	// Ctrl+Cで中断したときに、現在実行中のseedを表示する
	currentlyRunnningSeeds := make([]int, 0, len(seeds))
	var seedMutex sync.Mutex
	var datasMutex sync.Mutex
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go handleSignals(sigCh, &wg, &currentlyRunnningSeeds)

	for _, seed := range seeds {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			seedMutex.Lock()
			currentlyRunnningSeeds = append(currentlyRunnningSeeds, seed)
			seedMutex.Unlock()

			defer func() {
				seedMutex.Lock()
				for i, s := range currentlyRunnningSeeds {
					if s == seed {
						currentlyRunnningSeeds = append(currentlyRunnningSeeds[:i], currentlyRunnningSeeds[i+1:]...)
						break
					}
				}
				seedMutex.Unlock()
				<-sem
				wg.Done()
				atomic.AddInt32(&taskCompleted, 1)
				printProgress(int(taskCompleted), totalTask)
			}()
			var data map[string]float64
			var err error
			if cnf.Reactive {
				data, err = reactiveRun(cnf, seed)
			} else {
				data, err = runVis(cnf, seed)
			}

			if err != nil {
				errorChan <- fmt.Sprintf("Run error: seed=%d %v\n", seed, err)
				return
			}
			// 処理結果を格納
			datasMutex.Lock()
			datas = append(datas, data)
			datasMutex.Unlock()
			//fmt.Fprintf(os.Stderr, "%v\n", data)
		}(seed)
	}
	wg.Wait()
	close(errorChan)

	fmt.Fprintf(os.Stderr, "\n") // Newline after progress bar

	for err := range errorChan {
		log.Println(err)
	}

	sumScore := 0.0
	for i := 0; i < len(datas); i++ {
		//fmt.Println(datas[i])
		sumScore += datas[i]["Score"]
	}
	DisplayTable(datas)
	fmt.Fprintf(os.Stderr, "sumScore=%.2f\n", sumScore)
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

func handleSignals(sigCh <-chan os.Signal, wg *sync.WaitGroup, curent *[]int) {
	for {
		sig := <-sigCh
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			fmt.Println("\nReceived signal:", sig)
			fmt.Println("Currently running seeds:", *curent)
			os.Exit(1)
		}
	}

}
