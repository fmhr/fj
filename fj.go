package fj

import (
	"flag"
	"log"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Fj() {
	mode := flag.String("mode", "help", "モードの選択 [init(設定ファイルの生成), test(コマンドの実行), gen(fj -mode gen -seed 1 (シード1のテストwケースw生成) gcloud")
	seed := flag.Int("seed", 0, "テストケースのシード値を指定。未設定の場合は0 (0000.txt) が選択されます。-end オプションが指定されている場合、このオプションは無視されます。")
	start := flag.Int("start", 0, "開始シードの指定。デフォルトは０")
	end := flag.Int("end", 0, "[start, end) のシードケースがテストされる.")
	cmdArgs := flag.String("cmdArgs", "", "指定されたとき、 `Cmd [cmdArgs]`として実行される")
	jobs := flag.Int("jobs", -1, "同時実行するwジョブの上限w。。CPUのコア数を越えるとパフォーマンスが低下。")
	cloud := flag.Bool("cloud", false, "クラウドモード")
	flag.Parse()

	switch *mode {
	case "init":
		GenerateConfig()
		return
	case "test", "gen", "cloud":
		cnf, err := LoadConfigFile()
		if err != nil {
			log.Fatal(err)
		}
		if *cloud {
			cnf.Cloud = true
		}
		if *cmdArgs != "" {
			cnf.Cmd = cnf.Cmd + " " + *cmdArgs
		}
		seeds := []int{}
		if *end > *start {
			for i := *start; i < *end; i++ {
				seeds = append(seeds, i)
			}
		} else {
			seeds = append(seeds, *seed)
		}
		if *jobs > 0 {
			cnf.Jobs = *jobs
		}
		// modeによって処理を分岐
		switch *mode {
		case "test":
			RunParallel(cnf, seeds)
		case "gen":
			Gen(cnf, *seed)
		case "cloud":
			log.Println("cloud mode")
			rtn, err := CloudRun(cnf, *seed)
			if err != nil {
				log.Fatal(err)
			}
			log.Println(rtn)
		}
	default:
		flag.Usage()
	}

	if *mode == "init" {
		GenerateConfig()
		return
	}
}
