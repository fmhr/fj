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
	fj         = kingpin.New("fj", "fj is a command line tool for competitive programming.")
	debug      = fj.Flag("debug", "Enable debug mode.").Default("false").Bool()
	cloud      = fj.Flag("cloud", "Enable cloud mode.").Default("false").Bool()
	jsonOutput = fj.Flag("json", "Output json format.").Default("false").Bool()

	setup = fj.Command("init", "Generate config file.")
	force = setup.Flag("force", "Force generate config file.").Default("false").Bool()

	setupcloud = fj.Command("setupCloud", "Generate Dockerfile and gcloud build files for cloud mode.")

	test  = fj.Command("test", "Run test case.")
	seed  = test.Arg("seed", "Seed value.").Default("0").Int()
	args1 = test.Flag("args", "Command line arguments.").Strings()

	tests = fj.Command("tests", "Run test cases.")
	seed2 = tests.Arg("seed", "Seed value.").Int()
	args2 = tests.Flag("args", "Command line arguments.").Strings()
	start = tests.Flag("start", "Start seed value.").Default("0").Short('s').Int()
	end   = tests.Flag("end", "End seed value.").Default("10").Short('e').Int()
	jobs  = tests.Flag("jobs", "Number of parallel jobs.").Int()
	//cloud = tests.Flag("cloud", "Enable cloud mode.").Default("false").Bool()

)

func Fj() {
	if debug != nil && *debug {
		log.Println("debug mode")
	}

	switch kingpin.MustParse(fj.Parse(os.Args[1:])) {
	// Setup generate config file
	case setup.FullCommand():
		GenerateConfig()
	case setupcloud.FullCommand():
		mkDirCompilerBase()
		mkDirWorkerBase()
	// Test run test case
	case test.FullCommand(), tests.FullCommand():
		config, err := LoadConfigFile()
		if err != nil {
			log.Fatal(err)
		}
		// config を読み込む
		updateConfig(config)
		// cloud mode ならソースコードをアップロードしてバイナリを受け取る
		if config.Cloud {
			config.tmpBinary, err = CloudCompile(config)
			if err != nil {
				log.Fatal("Cloud mode Compile error:", err)
			}
		}
		// select run
		switch kingpin.MustParse(fj.Parse(os.Args[1:])) {
		case test.FullCommand():
			rtn, err := RunSelector(config, *seed)
			if err != nil {
				log.Fatal("Error: ", err)
			}
			fmt.Fprintln(os.Stdout, rtn)
		case tests.FullCommand():
			if seed2 != nil {
				*start = 0
				*end = *seed2
			}
			seeds := make([]int, *end-*start)
			for i := *start; i < *end; i++ {
				seeds[i-*start] = i
			}
			log.Println(*start, *end, seeds)
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

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
