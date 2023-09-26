package fj

import (
	"flag"
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

func Fj() {
	seed := flag.Int("seed", 0, "テストケースのシード値を指定 未設定の場合は0 (0000.txt) が選択されます")
	start := flag.Int("start", 0, "開始シードの指定 デフォルト:０")
	end := flag.Int("end", 10, "終了シード デフォルト:10")
	cmdArgs := flag.String("cmdArgs", "", "`Cmd [cmdArgs]`として実行される")
	jobs := flag.Int("jobs", 1, "並列実行数 デフォルト:1")
	cloud := flag.Bool("cloud", false, "クラウドモード デフォルト:false")
	flag.Parse()

	if len(flag.Args()) < 1 {
		fmt.Println("Usage: fj [command] [options]")
		os.Exit(1)
	}

	cmd := flag.Args()[0]

	if cmd == "init" {
		GenerateConfig()
		return
	}
	config, err := LoadConfigFile()
	if err != nil {
		log.Fatal(err)
	}

	if config.Cmd == "" {
		fmt.Printf("%s dファイルで実行コマンドを設定してください\n", configFileName)
		os.Exit(1)
	}

	// check flags
	if isFlagSet("cmdArgs") {
		config.Args = []string{*cmdArgs}
	}
	if isFlagSet("jobs") {
		config.Jobs = *jobs
	}
	if isFlagSet("cloud") {
		config.Cloud = *cloud
	}
	log.Println(*seed)
	switch cmd {
	case "test":
		var rtn map[string]float64
		var err error
		if config.Reactive {
			rtn, err = ReactiveRun(config, *seed)
		} else {
			rtn, err = RunVis(config, *seed)
		}
		if err != nil {
			log.Fatal("Error: ", err)
		}
		fmt.Println(rtn)
	case "tests":
		seeds := make([]int, *end-*start)
		for i := *start; i < *end; i++ {
			seeds[i-*start] = i
		}
		RunParallel(config, seeds)
	}
}

func isFlagSet(flagName string) bool {
	flagSet := false

	flag.Visit(func(f *flag.Flag) {
		if f.Name == flagName {
			flagSet = true
		}
	})
	return flagSet
}
