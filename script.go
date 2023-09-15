package fj

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
var OUTFILE_FOLDER = "tmp/"

func Fj() {

	Mode := flag.String("mode", "run", "select mode (init, run, gcloud)")
	//reactive := flag.Bool("reactive", false, "reactive")
	seed := flag.Int("seed", 0, "seed for testcase")
	start := flag.Int("start", 0, "seed for start")
	end := flag.Int("end", 0, "seed for end")
	cmdArgs := flag.String("cmdArgs", "", "cmdArgs")
	flag.Parse()

	// load config.toml
	cnf, err := LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	if *cmdArgs != "" {
		cnf.Cmd = cnf.Cmd + " " + *cmdArgs
	}
	// set seeds
	seeds := make([]int, 0)
	for i := *start; i < *end; i++ {
		seeds = append(seeds, i)
	}
	// mode select
	switch *Mode {
	case "run":
		if len(seeds) == 0 {
			err := RunVis(cnf, *seed)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			RunParallel(cnf, seeds)
		}
	case "init":
		GenerateConfig()
	case "gcloud":
		gcloud()
	}

	args := os.Args
	if len(args) > 1 {
		if args[1] == "t" {
			for i := 0; i < 10; i++ {
				RunVis(cnf, i)
			}
			return
		} else if args[1] == "init" {
			GenerateConfig()
			return
		} else if args[1] == "run" {
			RunParallel(cnf, []int{1})
			return
		}
	}
}
