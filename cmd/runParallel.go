package main

import (
	"log"
	"sync"
)

func RunParallel(seeds []int) {
	CORE := 4
	var wg sync.WaitGroup
	sem := make(chan struct{}, CORE)
	for _, seed := range seeds {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			_, err := tester(seed)
			if err != nil {
				log.Printf("seed=%d %v\n", seed, err)
			}
		}(seed)
	}
	wg.Wait()
}
