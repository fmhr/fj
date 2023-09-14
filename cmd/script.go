package main

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var TESTER = "./tools/target/release/tester"
var VIS = "./tools/target/release/vis"
var OUTFILE = "out.txt"
var INFILE_FOLDER = "tools/in/"

func main() {
	args := os.Args
	if len(args) > 1 && args[1] == "init" {
		GenerateConfig()
		return
	}

	app := flag.String("app", "", "app name")
	//reactive := flag.Bool("reactive", false, "reactive")
	seed := flag.Int("seed", 0, "seed for testcase")
	start := flag.Int("start", 0, "seed for start")
	end := flag.Int("end", 0, "seed for end")
	cmd := flag.String("cmd", "", "cmd")
	cmdArgs := flag.String("cmdArgs", "", "cmdArgs")
	flag.Parse()

	cnf, err := LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	if *cmd != "" {
		cnf.Cmd = *cmd
	}
	if *cmdArgs != "" {
		cnf.Cmd = cnf.Cmd + " " + *cmdArgs
	}
	// ---------------------------------------
	if len(args) > 1 && args[1] == "t" {
		seeds := make([]int, 10)
		for i := 0; i < 10; i++ {
			seeds[i] = 1
		}
		RunVis10(cnf)
		log.Println("t")
		return
	}
	// ---------------------------------------

	var seeds []int
	if start != nil && end != nil {
		for i := *start; i < *end; i++ {
			seeds = append(seeds, i)
		}
	} else if end != nil {
		for i := 0; i < *end; i++ {
			seeds = append(seeds, i)
		}
	} else if seed != nil {
		seeds = append(seeds, *seed)
	}

	//log.Println(args, *seed)
	switch *app {
	case "init":
		GenerateConfig()
	case "runVis":
		err := RunVis(cnf, *seed)
		if err != nil {
			log.Fatal(err)
		}
	case "runVis10":
		err := RunVis10(cnf)
		if err != nil {
			log.Fatal(err)
		}
	case "run":
		log.Printf("start=%d end=%d\n", *start, *end)
		RunParallel(cnf, seeds)
	case "gcloud":
		gcloud()
	case "seedSearch":
		seedSorting()
	default:
		log.Fatal("invalid command")
	}
}
