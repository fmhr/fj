package fj

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func RunParallel(cnf *Config, seeds []int) {
	// 並列実行数の設定
	concurrentNum := runtime.NumCPU() - 1
	if cnf.Jobs > 0 {
		concurrentNum = cnf.Jobs
	}
	if cnf.Cloud {
		concurrentNum = maxInt(1, maxInt(concurrentNum, cnf.ConcurrentRequests))
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrentNum)
	datas := make([]map[string]float64, 0, len(seeds))
	errorChan := make(chan string, len(seeds))
	errorSeedChan := make(chan int, len(seeds))

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
			data, err := RunSelector(cnf, seed)
			if err != nil {
				errorChan <- fmt.Sprintf("Run error: seed=%d %v\n", seed, err)
				errorSeedChan <- seed
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
	fmt.Fprintf(os.Stderr, "\n") // Newline after progress bar
	close(errorChan)
	close(errorSeedChan)
	for err := range errorChan {
		log.Println(err)
	}
	errSeeds := make([]int, 0, len(errorSeedChan))
	for seed := range errorSeedChan {
		errSeeds = append(errSeeds, seed)
	}
	sumScore := 0.0
	for i := 0; i < len(datas); i++ {
		//fmt.Println(datas[i])
		sumScore += datas[i]["Score"]
	}
	//	log.Println(datas)
	DisplayTable(datas)
	fmt.Fprintln(os.Stderr, "Error seeds:", errSeeds)
	// timeがあれば、平均と最大を表示
	_, exsit := datas[0]["time"]
	if exsit {
		sumTime := 0.0
		maxTime := 0.0
		for i := 0; i < len(datas); i++ {
			sumTime += datas[i]["time"]
			maxTime = math.Max(maxTime, datas[i]["time"])
		}
		sumTime /= float64(len(datas))
		fmt.Fprintf(os.Stderr, "avarageTime=%.2f  maxTime=%.2f\n", sumTime, maxTime)
	}
	avarageScore := sumScore / float64(len(datas))
	p := message.NewPrinter(language.English)
	p.Fprintf(os.Stderr, "(Score)sum=%.2f avarage=%.2f \n", sumScore, avarageScore)
	fmt.Printf("%.2f\n", sumScore)
	if jsonOutput != nil && *jsonOutput {
		fileContent, err := json.MarshalIndent(datas, "", " ")
		if err != nil {
			log.Fatal("json marshal error:", err)
		}
		err = createDirIfNotExist("fj/data/")
		if err != nil {
			log.Fatal("create dir error:", err)
		}
		now := time.Now()
		filename := fmt.Sprintf("fj/data/result_%s.json", fmt.Sprintf("%04d%02d%02d_%02d%02d%02d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second()))
		err = os.WriteFile(filename, fileContent, 0644)
		if err != nil {
			log.Fatal("json write error:", err)
		}
		log.Println("save json file:", filename)
	}

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

func createDirIfNotExist(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}
