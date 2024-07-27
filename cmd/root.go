package fj

import (
	"fmt"
	"log"
	"os"

	"github.com/alecthomas/kingpin/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/fmhr/fj/cmd/download"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

var (
	fj         = kingpin.New("fj", "fj is a command line tool for competitive programming.")
	debug      = fj.Flag("debug", "Enable debug mode.").Default("false").Bool()
	cloud      = fj.Flag("cloud", "Enable cloud mode.").Default("false").Bool()
	jsonOutput = fj.Flag("json", "Output json format.").Default("false").Bool()
	csvOutput  = fj.Flag("csv", "Output csv format.").Default("false").Bool()

	setup = fj.Command("init", "Generate config file.")

	setupcloud = fj.Command("setupCloud", "Generate Dockerfile and gcloud build files for cloud mode.")

	test  = fj.Command("test", "Run test case.")
	seed  = test.Arg("seed", "Seed value.").Default("0").Int()
	args1 = test.Flag("args", "Command line arguments.").Strings()

	tests        = fj.Command("tests", "Run test cases.")
	seed2        = tests.Arg("seed", "Seed value.").Int()
	args2        = tests.Flag("args", "Command line arguments.").Strings()
	start        = tests.Flag("start", "Start seed value.").Default("0").Short('s').Int()
	end          = tests.Flag("end", "End seed value.").Default("10").Short('e').Int()
	jobs         = tests.Flag("jobs", "Number of parallel jobs.").Int()
	displayTable = tests.Flag("table", "Output table format.").Default("true").Bool()
	Logscore     = tests.Flag("logscore", "Output score log.").Default("false").Bool()

	// downloadcmd tester file from URL
	downloadcmd = fj.Command("download", "Download tester file from URL.").Alias("d")
	testerURL   = downloadcmd.Arg("url", "Tester file URL.").Required().String()
	// login to atcoder
	login    = fj.Command("login", "Login to fj.").Alias("l")
	username = login.Flag("username", "Username.").Required().Short('u').String()
	password = login.Flag("password", "Password.").Required().Short('p').String()
	loginurl = login.Arg("url", "URL.").Default("https://atcoder.jp/login?").String()
	// check reactive 開発テスト用
	checkReactive = fj.Command("checkReactive", "Check if tester is reactive.")
)

func Execute() {
	if debug != nil && *debug {
		log.Println("debug mode")
	}

	result := kingpin.MustParse(fj.Parse(os.Args[1:]))
	switch result {
	// Setup generate config file
	case setup.FullCommand():
		err := generateConfig()
		if err != nil {
			log.Fatal(err)
		}
	case setupcloud.FullCommand():
		mkDirCompilerBase()
		mkDirWorkerBase()
	// Test run test case
	case test.FullCommand(), tests.FullCommand():
		config, err := setConfig()
		if err != nil {
			log.Fatal(err)
		}
		// config を読み込む
		updateConfig(config)
		// cloud mode ならソースコードをアップロードしてバイナリを受け取る
		if config.CloudMode {
			err = CloudCompile(config)
			if err != nil {
				log.Fatal("Cloud mode Compile error:", err)
			}
		}
		// select run
		switch result {
		case test.FullCommand():
			rtn, err := RunSelector(config, *seed)
			if err != nil {
				log.Fatal(err)
			}
			r, ok := rtn.Get("result")
			if ok {
				if r == "TLE" {
					log.Println("TLE")
				}
			}
			//fmt.Fprintln(os.Stdout, rtn)
			for _, k := range rtn.Keys() {
				v, ok := rtn.Get(k)
				if !ok {
					continue
				}
				p := message.NewPrinter(language.English)
				switch v := v.(type) {
				case int:
					p.Fprintf(os.Stderr, "%s:%d ", k, v)
				case float64:
					if v == float64(int(v)) {
						p.Fprintf(os.Stderr, "%s:%d ", k, int(v))
					} else {
						p.Fprintf(os.Stderr, "%s:%f ", k, v)
					}
				}
			}
			fmt.Fprintln(os.Stderr, "")
			Score, _ := rtn.Get("Score")
			fmt.Println(Score)
			stderr, ok := rtn.Get("stdErr")
			if ok {
				log.Print("StdErr:------->\n", string(stderr.([]byte)))
				log.Println("ここまで<-----")
			}
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
	case downloadcmd.FullCommand():
		download.Download(*testerURL)
	case login.FullCommand():
		download.Login(*loginurl, *username, *password)
	case checkReactive.FullCommand():
		fmt.Println("isReactive:", download.IsReactive())
	}
}

// updateConfig はコマンドライン引数でconfigを更新する
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
	config.CloudMode = config.CloudMode || *cloud
}
