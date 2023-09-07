package main

import (
	"fmt"
	"log"
	"sync"
)

func RunParallel(seeds []int) {
	CORE := 4
	var wg sync.WaitGroup
	sem := make(chan struct{}, CORE)
	datas := make([]map[string]float64, len(seeds))
	for _, seed := range seeds {
		wg.Add(1)
		sem <- struct{}{}
		go func(seed int) {
			defer func() {
				<-sem
				wg.Done()
			}()
			out, err := Run(seed)
			if err != nil {
				log.Printf("seed=%d %v\n", seed, err)
			}
			data, err := ExtractKeyValuePairs(string(out))
			if err != nil {
				log.Printf("seed=%d %v\n", seed, err)
			}
			data["seed"] = float64(seed)
			datas[seed] = data
		}(seed)
	}
	wg.Wait()
	sumScore := 0.0
	for i := 0; i < len(datas); i++ {
		log.Println(datas[i])
		sumScore += datas[i]["score"]
	}
	log.Printf("sumScore=%.2f\n", sumScore)
	fmt.Printf("%.2f\n", sumScore)
}
