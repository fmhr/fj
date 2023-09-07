package main

import (
	"flag"
	"fmt"
	"log"
)

func main() {
	log.SetFlags(log.Lshortfile)
	script()
}

func script() {
	app := flag.String("app", "", "app name")
	seed := flag.Int("seed", 0, "seed for testcase")
	start := flag.Int("start", 0, "seed for start")
	end := flag.Int("end", 0, "seed for end")
	flag.Parse()

	var seeds []int
	if start != nil && end != nil {
		for i := *start; i < *end; i++ {
			seeds = append(seeds, i)
		}
	}
	//log.Println(args, *seed)
	switch *app {
	case "tester":
		_, err := RunVis(*seed)
		if err != nil {
			log.Fatal(err)
		}
	case "tester10":
		err := tester10()
		if err != nil {
			log.Fatal(err)
		}
	case "run":
		fmt.Printf("start=%d end=%d\n", *start, *end)
		RunParallel(seeds)
	case "gcloud":
		gcloud()
	case "seedSearch":
		seedSorting()
	default:
		log.Fatal("invalid command")
	}
}
