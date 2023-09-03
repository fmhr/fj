package main

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os/exec"
	"sync"
)

// 毎回URLを変える必要がある
var GCR_URL = "https://go-cloud-run-sample-upxz55snfq-an.a.run.app/"

func runJob(jobID int, wg *sync.WaitGroup, results chan<- string) {
	defer wg.Done()

	requestURL := fmt.Sprintf("%s?seed=%d", GCR_URL, jobID)
	resp, err := http.Get(requestURL)
	if err != nil {
		results <- fmt.Sprintf("Job %d: ERROR: %s", jobID, err)
		return
	}
	defer resp.Body.Close()
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		results <- fmt.Sprintf("Job %d: ERROR: %s", jobID, err)
		return
	}
	results <- fmt.Sprintf("Job %d: %s", jobID, string(b))
}

func gcloud() {
	var wg sync.WaitGroup

	numJobs := 10
	results := make(chan string, numJobs)

	for i := 0; i < numJobs; i++ {
		wg.Add(1)
		go runJob(i, &wg, results)
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	for r := range results {
		fmt.Println(r)
	}
}

// コンテスト中はテスターが変わらないと仮定すれば、URLは変わらないので、いらないかも
func getCloudRunURL(serviceName, region string) (string, error) {
	cmd := exec.Command("gcloud", "run", "services", "describe", serviceName, "--region", region, "--format", "value(status.url)")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return out.String(), nil
}
