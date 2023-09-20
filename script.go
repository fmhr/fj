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
	// set cmdArgs
	if *cmdArgs != "" {
		cnf.Cmd = cnf.Cmd + " " + *cmdArgs
	}
	// set seeds
	seeds := make([]int, 0)
	for i := *start; i < *end; i++ {
		seeds = append(seeds, i)
	}
	// set jobs
	if *jobs > 0 {
		cnf.Jobs = *jobs
	}
	// mode select
	switch *mode {
	case "run":
		// １つのseedを実行
		if len(seeds) == 0 {
			var err error
			if cnf.Reactive {
				_, err = ReactiveRun(cnf, *seed)
			} else {
				_, err = RunVis(cnf, *seed)
			}
			if err != nil {
				log.Fatal(err)
			}
		} else {
			// 複数のseedを並列実行
			//if cnf.Reactive {
			//ReactiveRunParallel(cnf, seeds)
			//} else {
			RunParallel(cnf, seeds)
			//}
		}
	case "init":
		// 設定ファイルの生成
		GenerateConfig()
	case "gen":
		Gen(cnf, *seed)
	case "gcloud":
		// TODO
		gcloud()
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
