package cmd

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

// RunParallel 複数のシードに対して並列にテストを実行する
func RunParallel(cnf *setup.Config, seeds []int) {
	// 並列実行数の設定
	concurrentNum := 1
	if cnf.Jobs > 0 {
		concurrentNum = cnf.Jobs
	}
	if cnf.CloudMode {
		concurrentNum = max(1, cnf.ConcurrentRequests)
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrentNum)
	datas := make([]*orderedmap.OrderedMap[string, string], 0, len(seeds))
	errorChan := make(chan string, len(seeds))
	errorSeedChan := make(chan int, len(seeds))

	var taskCompleted int32 = 0
	totalTask := len(seeds)

	var datasMutex sync.Mutex
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	printProgress(int(taskCompleted), totalTask) // プログレスバーの表示

	// エラーが出たらそこで打ち止めにする
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Ctrl+Cで中断したときに、現在実行中のseedを表示して、それ以降はキャンセルする
	go func() {
		sig := <-sigCh
		fmt.Printf("\nReceived signal: %v - waiting for running tasks to complete\n", sig)
		cancel()
	}()
overloop:
	for _, seed := range seeds {
		select {
		case <-ctx.Done():
			break overloop
		default:
			wg.Add(1)
			sem <- struct{}{}
			time.Sleep(5 * time.Millisecond)
			go func(seed int) {
				defer wg.Done()
				defer func() { <-sem }()
				data, err := RunSelector(cnf, seed)
				if err != nil {
					errorChan <- fmt.Sprintf("Run error: seed=%d %v\n", seed, err)
					errorSeedChan <- seed
					log.Printf("seed=%d has Error %v\n", seed, err)
				} else {
					datasMutex.Lock()
					datas = append(datas, data) // 結果を追加
					datasMutex.Unlock()
					currentTaskCompleted := atomic.AddInt32(&taskCompleted, 1) // progressbar
					printProgress(int(currentTaskCompleted), totalTask)        // progressbar
				}
			}(seed)
		}
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
	sumScore := 0
	logScore := 0.0
	zeroSeeds := make([]int, 0)
	tleSeeds := make([]int, 0)
	for i := 0; i < len(datas); i++ {
		seedStr, ok := datas[i].Get("seed")
		if !ok {
			log.Println("seed not found")
			continue
		}
		v, ok := datas[i].Get("result")
		if ok {
			if v == "TLE" {
				seed, err := strconv.Atoi(seedStr)
				if err != nil {
					log.Println("seed not found")
					continue
				}
				tleSeeds = append(tleSeeds, seed)
			}
		}

		// error のときnilになる
		if datas[i] != nil {
			scoreStr, ok := datas[i].Get("Score")
			if !ok {
				log.Println("Score not found")
				continue
			}
			score, err := strconv.Atoi(scoreStr)
			if err != nil {
				log.Println("Score is not int")
				continue
			}
			if scoreStr == "-1" {
				continue
			}
			sumScore += score
			logScore += math.Max(0, math.Log(float64(score)))
			seed, ok := datas[i].Get("seed")
			if !ok {
				log.Println("seed not found")
				continue
			}
			if score == 0 {
				seed, err := strconv.Atoi(seed)
				if err != nil {
					continue
				}
				zeroSeeds = append(zeroSeeds, seed)
			}
		}
	}
	// テーブル表示 代替を考えるまでtrue
	if true {
		bestScores, err := GetBestScores()
		if err != nil {
			log.Println("Error GetBestScores:", err)
		}
		// delete stdErr
		for i := 0; i < len(datas); i++ {
			if datas[i] != nil {
				datas[i].Delete("stdErr")
			}
			seedStr, ok := datas[i].Get("seed")
			if !ok {
				log.Println("seed not found")
				continue
			}
			seed, _ := strconv.Atoi(seedStr)
			if bestScore, ok := bestScores[seed]; ok {
				datas[i].Set("Best", fmt.Sprintf("%d", bestScore))
				if scoreStr, ok := datas[i].Get("Score"); ok {
					datas[i].Set("Score", scoreStr)
					score, _ := strconv.Atoi(scoreStr)
					if score >= 0 && bestScore >= 0 {
						ratio := float64(score) / float64(bestScore)
						datas[i].Set("Ratio", fmt.Sprintf("%.2f", ratio))
					}
				}
			}
		}
		err = DisplayTable(datas)
		if err != nil {
			log.Println("Error DisplayTable:", err)
			// 表示できないけど、結果は出力したいので、続行
		}
	}
	if len(errSeeds) > 0 {
		fmt.Fprintln(os.Stderr, "Errors:", errSeeds)
	}
	if len(zeroSeeds) > 0 {
		fmt.Fprintln(os.Stderr, "Zeros:", zeroSeeds)
	}
	if len(tleSeeds) > 0 {
		fmt.Fprintln(os.Stderr, "TLEs:", tleSeeds)
	}
	// timeがあれば、平均と最大を表示
	if _, exsit := datas[0].Get("time"); exsit {
		sumTime := 0.0
		maxTime := 0.0
		nums := 0
		for i := 0; i < len(datas); i++ {
			if datas[i] == nil {
				continue
			}
			if t, ok := datas[i].Get("time"); ok {
				timef, err := strconv.ParseFloat(t, 64)
				if err != nil {
					log.Println("Error time:", err)
					continue
				}
				sumTime += timef
				maxTime = math.Max(maxTime, timef)
				nums++
			}
		}
		fmt.Fprintf(os.Stderr, "(Time)avarage:%.2f  max:%.2f\n", sumTime/float64(nums), maxTime)
	}
	avarageScore := sumScore / (len(datas) - len(errSeeds))
	p := message.NewPrinter(language.English)
	p.Fprintf(os.Stderr, "(Score)avarage:%d sum:%d\n", avarageScore, sumScore)

	if jsonOutput != nil && *jsonOutput {
		err := JsonOutput(datas)
		if err != nil {
			log.Println("Error JsonOutput:", err)
		}
	}
	if csvOutput != nil && *csvOutput {
		err := CsvOutput(datas)
		if err != nil {
			log.Println("Error CsvOutput:", err)
		}
	}
	// update best score
	for i := 0; i < len(datas); i++ {
		if datas[i] == nil {
			continue
		}
		scoreString, ok := datas[i].Get("Score")
		if !ok {
			log.Println("Score not found")
			continue
		}
		seedStr, ok := datas[i].Get("seed")
		if !ok {
			log.Println("seed not found")
			continue
		}
		score, _ := strconv.Atoi(scoreString)
		seed, _ := strconv.Atoi(seedStr)
		UpdateBestScore(seed, score)
	}
}

const progressBarWidth = 40

// printProgress 進行度を表示する
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

// createDirIfNotExist 使用するディレクトリが存在しない時に作成する
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
