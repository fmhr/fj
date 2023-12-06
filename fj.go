package fj

import (
	"fmt"
	"log"
	"math"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
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

	result := kingpin.MustParse(fj.Parse(os.Args[1:]))
	switch result {
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
		switch result {
		case test.FullCommand():
			rtn, err := RunSelector(config, *seed)
			if err != nil {
				log.Fatal("Error: ", err)
			}
			fmt.Fprintln(os.Stdout, rtn)
			for k, v := range rtn {
				p := message.NewPrinter(language.English)
				if v == math.Floor(v) {
					p.Fprintf(os.Stdout, "%s:%d ", k, int(v))
				} else {
					p.Fprintf(os.Stdout, "%s:%f ", k, v)
				}
			}
			fmt.Println("")
		case tests.FullCommand():
			// seed2 が指定されていれば end=seed2
			if seed2 != nil && *seed2 != 0 {
				*start = 0
				*end = *seed2
			}
			// start, endが指定されていれば、その範囲のシードを並列実行
			seeds := make([]int, *end-*start)
			for i := *start; i < *end; i++ {
				seeds[i-*start] = i
			}
			RunParallel(config, seeds)
		}
	}

}

func updateConfig(config *Config) {
	if *args1 != nil && len(*args1) > 0 {
		config.Args = *args1
	}
	if args2 != nil && len(*args2) > 0 {
		config.Args = *args2
	}
	if jobs != nil && *jobs > 0 {
		config.Jobs = *jobs
	}
	config.Cloud = config.Cloud || *cloud
}

func maxInt(a, b int) int {
	if a > b {
		return a
	}
	return b
}
