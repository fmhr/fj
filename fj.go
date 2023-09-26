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
	args2 = tests.Flag("args", "Command line arguments.").Strings()
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
