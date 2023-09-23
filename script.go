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
	mode := flag.String("mode", "help", "Chose mode [init(generage fj_config.toml), run(execute cmd), gen(need -seed option) gcloud(TODO))")
	seed := flag.Int("seed", 0, "Seed for testcase seed. If unset, the default is 0 (0000.txt)run 0 seeds(0000.txt). If -end is set, thid value is ignored.")
	start := flag.Int("start", 0, "Specifies the starting seed number. Default is 0.")
	end := flag.Int("end", 0, "If set, run for seeds in the range: [start, end).")
	cmdArgs := flag.String("cmdArgs", "", "If provides, , will run `Cmd [cmdArgs]`")
	jobs := flag.Int("jobs", -1, "Set the limit for concurrently executed jobs. if Above the number of CPU cores may decrease performance.")
	cloud := flag.Bool("cloud", false, "If set, run on cloud mode.")
	flag.Parse()
	if *mode == "init" {
		GenerateConfig()
		return
	}

	// load config.toml
	cnf, err := LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	if cloud != nil && *cloud {
		cnf.Cloud = true
	}
	// set cmdArgs
	if *cmdArgs != "" {
		cnf.Cmd = cnf.Cmd + " " + *cmdArgs
	}
	// set seeds
	seeds := make([]int, 0)
	for i := *start; i < *end; i++ {
		seeds = append(seeds, i)
	}
	if len(seeds) == 0 {
		seeds = append(seeds, *seed)
	}
	// set jobs
	if *jobs > 0 {
		cnf.Jobs = *jobs
	}
	// mode select
	switch *mode {
	case "run":
		RunParallel(cnf, seeds)
	case "init":
		// 設定ファイルの生成
		GenerateConfig()
	case "gen":
		// テストケースの生成
		Gen(cnf, *seed)
	case "cloud":
		log.Println("cloud mode")
		rtn, err := CloudRun(cnf, *seed)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(rtn)
	default:
		flag.Usage()
	}

	// oj mode
	args := os.Args
	if len(args) > 1 {
		if args[1] == "t" {
			for i := 0; i < 10; i++ {
				RunVis(cnf, i)
			}
			return
		}
	}
}
