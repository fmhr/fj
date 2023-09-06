package main

import (
	"flag"
	"log"
	"os/exec"
)

func main() {
	log.SetFlags(log.Lshortfile)
	script()
}

func script() {
	app := flag.String("app", "", "app name")
	seed := flag.Int("seed", 0, "seed for testcase")
	startSeed := flag.Int("startSeed", 0, "seed for start")
	endSeed := flag.Int("endSeed", 0, "seed for end")
	flag.Parse()

	var seeds []int
	if startSeed != nil && endSeed != nil {
		for i := *startSeed; i <= *endSeed; i++ {
			seeds = append(seeds, i)
		}
	}
	//log.Println(args, *seed)
	switch *app {
	case "tester":
		_, err := tester(*seed)
		if err != nil {
			log.Fatal(err)
		}
	case "tester10":
		err := tester10()
		if err != nil {
			log.Fatal(err)
		}
	case "seedSearch":
		seedSorting()
	case "gcloud":
		gcloud()
	case "run":
		RunParallel(seeds)
	default:
		log.Fatal("invalid command")
	}
}

func build() error {
	cmd := exec.Command("go", "build", "-o", "main", "src/main.go")
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}

func buildTester() error {

	cmd2 := exec.Command("cargo", "build", "--manifest-path", "tools/Cargo.toml", "--release", "--bin", "tester")
	err := cmd2.Run()
	if err != nil {
		return err
	}
	return nil
}
