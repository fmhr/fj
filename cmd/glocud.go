package main

import (
	"fmt"
	"io"
	"net/http"
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
