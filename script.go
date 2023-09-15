package fj

import (
	"flag"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Fj() {
	Mode := flag.String("mode", "run", "select mode (init, run, gcloud)")
	seed := flag.Int("seed", 0, "seed for testcase")
	start := flag.Int("start", 0, "seed for start")
	end := flag.Int("end", 0, "seed for end")
	cmdArgs := flag.String("cmdArgs", "", "cmdArgs")
	flag.Parse()

	if *Mode == "init" {
		GenerateConfig()
		return
	}
	// load config.toml
	cnf, err := LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}
	// set flag value

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
		// １つのseedを実行
		if len(seeds) == 0 {
			err := RunVis(cnf, *seed)
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// 複数のseedを並列実行
			RunParallel(cnf, seeds)
		}
	case "init":
		// 設定ファイルの生成
		GenerateConfig()
	case "gcloud":
		// TODO
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
