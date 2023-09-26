package fj

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	fj    = kingpin.New("fj", "fj is a command line tool for competitive programming.")
	debug = fj.Flag("debug", "Enable debug mode.").Default("false").Bool()

	setup = fj.Command("init", "Generate config file.")

	test  = fj.Command("test", "Run test case.")
	seed  = test.Arg("seed", "Seed value.").Default("0").Int()
	args1 = test.Flag("args", "Command line arguments.").Strings()

	tests = fj.Command("tests", "Run test cases.")
	args2 = tests.Arg("args", "Command line arguments.").Strings()
	start = tests.Flag("start", "Start seed value.").Default("0").Int()
	end   = tests.Flag("end", "End seed value.").Default("10").Int()
	jobs  = tests.Flag("jobs", "Number of parallel jobs.").Default("1").Int()
	cloud = tests.Flag("cloud", "Enable cloud mode.").Default("false").Bool()
)

func Fj() {
	if debug != nil && *debug {
		log.Println("debug mode")
	}

	switch kingpin.MustParse(fj.Parse(os.Args[1:])) {
	// Setup generate config file
	case setup.FullCommand():
		GenerateConfig()
	// Test run test case
	case test.FullCommand(), tests.FullCommand():
		config, err := LoadConfigFile()
		if err != nil {
			log.Fatal(err)
		}
		updateConfig(config)
		switch kingpin.MustParse(fj.Parse(os.Args[1:])) {
		case test.FullCommand():
			rtn, err := Run(config, *seed)
			if err != nil {
				log.Fatal("Error: ", err)
			}
			fmt.Fprintln(os.Stdout, rtn)
		case tests.FullCommand():
			seeds := make([]int, *end-*start)
			for i := *start; i < *end; i++ {
				seeds[i-*start] = i
			}
			RunParallel(config, seeds)
		}
	}

}

func updateConfig(config *Config) {
	if args1 != nil {
		config.Args = *args1
	}
	if args2 != nil {
		config.Args = *args2
	}
	if jobs != nil {
		config.Jobs = *jobs
	}
	if cloud != nil {
		config.Cloud = *cloud
	}
}

//func Fj() {
//seed := flag.Int("seed", 0, "テストケースのシード値を指定 未設定の場合は0 (0000.txt) が選択されます")
//start := flag.Int("start", 0, "開始シードの指定 デフォルト:０")
//end := flag.Int("end", 10, "終了シード デフォルト:10")
//cmdArgs := flag.String("cmdArgs", "", "`Cmd [cmdArgs]`として実行される")
//jobs := flag.Int("jobs", 1, "並列実行数 デフォルト:1")
//cloud := flag.Bool("cloud", false, "クラウドモード デフォルト:false")
//flag.Parse()

//if len(flag.Args()) < 1 {
//fmt.Println("Usage: fj [command] [options]")
//os.Exit(1)
//}

//cmd := flag.Args()[0]

//if cmd == "init" {
//GenerateConfig()
//return
//}
//config, err := LoadConfigFile()
//if err != nil {
//log.Fatal(err)
//}

//if config.Cmd == "" {
//fmt.Printf("%s dファイルで実行コマンドを設定してください\n", configFileName)
//os.Exit(1)
//}

//// check flags
//if isFlagSet("cmdArgs") {
//config.Args = []string{*cmdArgs}
//}
//if isFlagSet("jobs") {
//config.Jobs = *jobs
//}
//if isFlagSet("cloud") {
//config.Cloud = *cloud
//}
//log.Println(*seed)
//switch cmd {
//case "test":
//var rtn map[string]float64
//var err error
//if config.Reactive {
//rtn, err = ReactiveRun(config, *seed)
//} else {
//rtn, err = RunVis(config, *seed)
//}
//if err != nil {
//log.Fatal("Error: ", err)
//}
//fmt.Println(rtn)
//case "tests":
//seeds := make([]int, *end-*start)
//for i := *start; i < *end; i++ {
//seeds[i-*start] = i
//}
//RunParallel(config, seeds)
//}
//}

//func isFlagSet(flagName string) bool {
//flagSet := false

//flag.Visit(func(f *flag.Flag) {
//if f.Name == flagName {
//flagSet = true
//}
//})
//return flagSet
//}
