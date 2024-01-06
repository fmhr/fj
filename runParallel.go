package fj

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

func RunParallel(cnf *Config, seeds []int) {
	// 並列実行数の設定
	concurrentNum := 1
	if cnf.Jobs > 0 {
		concurrentNum = cnf.Jobs
	}
	if cnf.Cloud {
		concurrentNum = maxInt(1, cnf.ConcurrentRequests)
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrentNum)
	datas := make([]*orderedmap.OrderedMap[string, any], 0, len(seeds))
	errorChan := make(chan string, len(seeds))
	errorSeedChan := make(chan int, len(seeds))

	var taskCompleted int32 = 0
	totalTask := len(seeds)

	// Ctrl+Cで中断したときに、現在実行中のseedを表示する
	currentlyRunnningSeeds := map[int]bool{}
	var datasMutex sync.Mutex
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// Ctrl+Cで中断したときに、現在実行中のseedを表示する
	go handleSignals(sigCh, &wg, &currentlyRunnningSeeds)

	printProgress(int(taskCompleted), totalTask)

	for _, seed := range seeds {
		wg.Add(1)
		sem <- struct{}{}
		currentlyRunnningSeeds[seed] = true
		time.Sleep(5 * time.Millisecond)
		go func(seed int) {
			data, err := RunSelector(cnf, seed)
			if err != nil {
				errorChan <- fmt.Sprintf("Run error: seed=%d %v\n", seed, err)
				errorSeedChan <- seed
			}
			// 後処理
			datasMutex.Lock()
			datas = append(datas, data)                  // 結果を追加
			atomic.AddInt32(&taskCompleted, 1)           // progressbar
			printProgress(int(taskCompleted), totalTask) // progressbar
			delete(currentlyRunnningSeeds, seed)         // 現在実行中のseedを削除
			wg.Done()
			<-sem
			datasMutex.Unlock()
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
	logScore := 0.0
	zeroSeeds := make([]int, 0)
	for i := 0; i < len(datas); i++ {
		//fmt.Println(datas[i])
		score, ok := datas[i].Get("Score")
		if !ok {
			log.Fatal("Score not found")
		}
		sumScore += score.(float64)
		logScore += math.Log(score.(float64))
		if score.(float64) == 0.0 {
			zeroSeeds = append(zeroSeeds, i)
		}
	}
	if displayTable != nil && *displayTable {
		DisplayTable(datas)
	}
	fmt.Fprintln(os.Stderr, "Error seeds:", errSeeds, "Zero seeds:", zeroSeeds)
	// timeがあれば、平均と最大を表示
	_, exsit := datas[0].Get("time")
	if exsit {
		sumTime := 0.0
		maxTime := 0.0
		for i := 0; i < len(datas); i++ {
			if t, ok := datas[i].Get("time"); !ok {
				log.Fatal("time is not float64")
			} else {
				sumTime += t.(float64)
				maxTime = math.Max(maxTime, t.(float64))
			}
		}
		sumTime /= float64(len(datas))
		fmt.Fprintf(os.Stderr, "avarageTime=%.2f  maxTime=%.2f\n", sumTime, maxTime)
	}
	avarageScore := sumScore / float64(len(datas))
	p := message.NewPrinter(language.English)
	p.Fprintf(os.Stderr, "(Score)sum=%.2f avarage=%.2f log=%.2f\n", sumScore, avarageScore, logScore)
	// if zeroSeeds があれば、sumScoreを０にする
	// TODO スコアが低ければいいい時は有効にする
	//if len(zeroSeeds) > 0 {
	//log.Println("Score 0 seeds:", zeroSeeds)
	//fmt.Println("0")
	//} else {
	//fmt.Printf("%.2f\n", sumScore)
	//}
	if Logscore != nil && *Logscore {
		fmt.Printf("%.4f\n", logScore)
	} else {
		fmt.Printf("%.2f\n", sumScore)
	}
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
	if csvOutput != nil && *csvOutput {
		CsvOutput(datas)
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

func handleSignals(sigCh <-chan os.Signal, wg *sync.WaitGroup, curent *map[int]bool) {
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
	} else if err != nil {
		return fmt.Errorf("error checking if directory exists: %v", err)
	}
	return nil
}
