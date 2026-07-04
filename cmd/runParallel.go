package cmd

import (
	"context"
	"fmt"
	"log"
	"math"
	"os"
	"os/signal"
	"slices"
	"strconv"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/fmhr/fj/cmd/setup"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
)

type runParallelMetrics struct {
	sumScore     int
	logScore     float64
	zeroSeeds    []int
	timeoutSeeds []int
	errSeeds     []int
	averageScore int
}

func runTestSeeds(cnf *setup.Config, seeds []int) (datas []SliceMap, errSeeds []int) {
	concurrentNum := 1
	if cnf.Jobs > 0 {
		concurrentNum = cnf.Jobs
	}
	if cnf.CloudMode {
		concurrentNum = max(1, cnf.ConcurrentRequests)
	}
	var wg sync.WaitGroup
	sem := make(chan struct{}, concurrentNum)
	datas = make([]SliceMap, 0, len(seeds))
	errorChan := make(chan string, len(seeds))
	errorSeedChan := make(chan int, len(seeds))

	var taskCompleted int32 = 0
	totalTask := len(seeds)

	var datasMutex sync.Mutex
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	printProgress(int(taskCompleted), totalTask) // プログレスバーの表示

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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
			time.Sleep(100 * time.Millisecond) // 100ms待つ エラー出力を見やすくするため
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
					datas = append(datas, data)
					datasMutex.Unlock()
					currentTaskCompleted := atomic.AddInt32(&taskCompleted, 1)
					printProgress(int(currentTaskCompleted), totalTask)
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

	for seed := range errorSeedChan {
		errSeeds = append(errSeeds, seed)
	}

	return datas, errSeeds
}

func calculateMetrics(datas []SliceMap, errSeeds []int) runParallelMetrics {
	var metrics runParallelMetrics
	metrics.errSeeds = errSeeds

	for i := 0; i < len(datas); i++ {
		if datas[i] == nil {
			continue
		}
		seedStr, ok := datas[i].Get("seed")
		if !ok {
			log.Println("seed not found")
			continue
		}
		v, ok := datas[i].Get("result")
		if ok && v == "TimeOut" {
			seed, err := strconv.Atoi(seedStr)
			if err == nil {
				metrics.timeoutSeeds = append(metrics.timeoutSeeds, seed)
			}
		}

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
		metrics.sumScore += score
		metrics.logScore += math.Max(0, math.Log(float64(score)))

		if score == 0 {
			seed, err := strconv.Atoi(seedStr)
			if err == nil {
				metrics.zeroSeeds = append(metrics.zeroSeeds, seed)
			}
		}
	}

	if len(datas) > len(errSeeds) {
		metrics.averageScore = metrics.sumScore / (len(datas) - len(errSeeds))
	}

	return metrics
}

func applyBestScores(datas []SliceMap) {
	bestScores, err := getBestScores()
	if err != nil {
		log.Println("Error GetBestScores:", err)
		return
	}
	for i := 0; i < len(datas); i++ {
		if datas[i] == nil {
			continue
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
}

func printTimeSummary(datas []SliceMap) {
	if len(datas) == 0 {
		return
	}
	if _, exist := datas[0].Get("time"); !exist {
		return
	}
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
	if nums > 0 {
		fmt.Fprintf(os.Stderr, "(Time)average:%.2f  max:%.2f\n", sumTime/float64(nums), maxTime)
	}
}

func printAverageOfNumericalFields(datas []SliceMap) {
	if len(datas) == 0 {
		return
	}
	keys := make([]string, 0)
	for _, data := range datas {
		for _, kvp := range data {
			if !slices.Contains(keys, kvp.Key) {
				keys = append(keys, kvp.Key)
			}
		}
	}
	for _, key := range keys {
		if key == "seed" || key == "result" {
			continue
		}
		sum := 0.0
		count := 0
		for i := 0; i < len(datas); i++ {
			if datas[i] == nil {
				continue
			}
			if v, ok := datas[i].Get(key); ok {
				f, err := strconv.ParseFloat(v, 64)
				if err == nil {
					sum += f
					count++
				}
			}
		}
		if count > 0 {
			fmt.Fprintf(os.Stderr, "(%s)average:%.4f\n", key, sum/float64(count))
		}
	}
}

func updateBestScoresWithResults(datas []SliceMap) {
	newScores := make(map[int]int)
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
		newScores[seed] = score
	}
	err := updateBestScores(newScores)
	if err != nil {
		log.Println("Error UpdateBestScores:", err)
	}
}

// runParallel 複数のシードに対して並列にテストを実行する
func runParallel(cnf *setup.Config, seeds []int) {
	datas, errSeeds := runTestSeeds(cnf, seeds)

	metrics := calculateMetrics(datas, errSeeds)

	if *bestScore {
		applyBestScores(datas)
	}

	err := DisplayTable(datas)
	if err != nil {
		log.Println("Error DisplayTable:", err)
	}

	if len(metrics.errSeeds) > 0 {
		fmt.Fprintln(os.Stderr, "Errors:", metrics.errSeeds)
	}
	if len(metrics.zeroSeeds) > 0 {
		fmt.Fprintln(os.Stderr, "Zeros:", metrics.zeroSeeds)
	}
	if len(metrics.timeoutSeeds) > 0 {
		fmt.Fprintln(os.Stderr, "TimeOuts:", metrics.timeoutSeeds)
	}

	printTimeSummary(datas)

	p := message.NewPrinter(language.English)
	p.Fprintf(os.Stderr, "(Score)average:%d sum:%d sum(log):%f\n", metrics.averageScore, metrics.sumScore, metrics.logScore)

	printAverageOfNumericalFields(datas)

	if jsonOutput != nil && *jsonOutput {
		err := writeJsonOutput(datas)
		if err != nil {
			log.Println("Error jsonOutput:", err)
		}
	}
	if csvOutput != nil && *csvOutput {
		err := writeCsvOutput(datas, "result")
		if err != nil {
			log.Println("Error csvOutput:", err)
		}
	}

	if *bestScore {
		updateBestScoresWithResults(datas)
	}

	if len(metrics.errSeeds) == 0 && len(metrics.zeroSeeds) == 0 && len(metrics.timeoutSeeds) == 0 {
		fmt.Fprintln(os.Stderr, "All tests passed")
	} else {
		fmt.Fprintln(os.Stderr, "Some tests failed")
	}
	fmt.Println(metrics.sumScore)
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


