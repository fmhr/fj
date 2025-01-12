package cmd

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/alecthomas/kingpin/v2"
	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/fmhr/fj/cmd/download"
	"github.com/fmhr/fj/cmd/setup"
)

const (
	FJ_DIRECTORY = "fj/"
)

func init() {
	// flagのパース前にinitが実行されることに注意
}

var (
	fj         = kingpin.New("fj", "fj is a command line tool for AtCoder Heuristic Contest.")
	debug      = fj.Flag("debug", "Enable debug mode.").Default("false").Bool()
	cloud      = fj.Flag("cloud", "Enable cloud mode.").Default("false").Bool()
	jsonOutput = fj.Flag("json", "Output json format.").Default("false").Bool()
	csvOutput  = fj.Flag("csv", "Output csv format.").Default("false").Bool()

	// init command
	setupCmd = fj.Command("init", "Generate config file.")
	minimax  = setupCmd.Flag("minimax", "BestScore is mini or max").Default("max").String()

	setupcloud = fj.Command("setupCloud", "Generate Dockerfile and gcloud build files for cloud mode.")

	// test command
	test     = fj.Command("test", "Run test case.").Alias("t")
	cmd      = test.Arg("cmd", "Exe Cmd.").Required().String()
	seed     = test.Flag("seed", "Set Seed. default : 0.").Short('s').Default("0").Int()
	count    = test.Flag("count", "Number of test cases.").Short('n').Default("1").Int()
	parallel = test.Flag("parallel", "Number of parallel jobs.").Short('p').Default("1").Int()
	//updateBest = test.Flag("updateBest", "Update best score.").Default("false").Bool()

	// downloadcmd tester file from URL
	downloadcmd = fj.Command("download", "Download tester file from URL.").Alias("d")
	testerURL   = downloadcmd.Arg("url", "Tester file URL.").Required().String()
	// login to atcoder
	login    = fj.Command("login", "Log in to AtCoder.").Alias("l")
	username = login.Flag("username", "Username.").Required().Short('u').String()
	password = login.Flag("password", "Password.").Required().Short('p').String()
	loginurl = login.Arg("url", "URL.").Default("https://atcoder.jp/login?").String()
	// logout
	logout = fj.Command("logout", "Log out from AtCoder.")
	// info
	info = fj.Command("info", "Show info.")
)

func Execute() error {
	// Parse command line arguments
	result := kingpin.MustParse(fj.Parse(os.Args[1:]))

	if debug != nil && *debug {
		log.Println("debug mode")
		log.SetFlags(log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	switch result {
	// Setup generate config file
	case setupCmd.FullCommand():
		// fj ディレクトリを作成
		if _, err := os.Stat(FJ_DIRECTORY); os.IsNotExist(err) {
			if err := os.Mkdir(FJ_DIRECTORY, 0755); err != nil {
				return err
			}
		}
		// fj/best_score.json の作成
		if *minimax == "min" {
			err := createNewBestScorejson(-1)
			if err != nil {
				return err
			}
		} else {
			err := createNewBestScorejson(1)
			if err != nil {
				return err
			}
		}

	case setupcloud.FullCommand():
		err := mkDirCompilerBase()
		if err != nil {
			return err
		}
		err = mkDirWorkerBase()
		if err != nil {
			return err
		}
	// Test run test case
	// test と　tests 時の共通処理
	case test.FullCommand():
		config, err := setup.SetConfig()
		if err != nil {
			return err
		}
		// config を読み込む
		updateConfig(config)
		// cloud mode ならソースコードをアップロードしてバイナリを受け取る
		if config.CloudMode {
			err = CloudCompile(config)
			if err != nil {
				log.Println("Cloud mode Compile error:", err)
				return err
			}
		}
		if *count == 1 {
			rtn, err := RunSelector(config, *seed)
			if err != nil {
				return err
			}
			r, ok := rtn.Get("result")
			if ok {
				if r == "TLE" {
					log.Println("TLE")
				}
			}
			for _, k := range rtn.Keys() {
				v, ok := rtn.Get(k)
				if !ok {
					continue
				}
				p := message.NewPrinter(language.English)
				p.Fprintf(os.Stderr, "%s:%s ", k, v)
			}
			fmt.Fprintln(os.Stderr, "")
			Score, _ := rtn.Get("Score")
			scoreInt, err := strconv.Atoi(Score)
			if err != nil {
				return err
			}
			UpdateBestScore(*seed, scoreInt)
		} else {
			startSeed := *seed
			config.Jobs = *parallel
			seeds := make([]int, *count)
			for i := startSeed; i < startSeed+*count; i++ {
				seeds[i-startSeed] = i
			}
			RunParallel(config, seeds)
		}
	case downloadcmd.FullCommand():
		return download.Download(*testerURL)
	case login.FullCommand():
		return download.Login(*loginurl, *username, *password)
	case logout.FullCommand():
		download.Logout()
	case info.FullCommand():
		fmt.Println("isReactive:", download.IsReactive())
		cDir, _ := os.UserCacheDir()
		fmt.Println("cacheDirectory:", cDir+"/fmhr-judge-tools")
	}
	return nil
}

// updateConfig はコマンドライン引数でconfigを更新する
func updateConfig(config *setup.Config) {
	config.CloudMode = config.CloudMode || *cloud
	config.ExecuteCmd = *cmd
}
