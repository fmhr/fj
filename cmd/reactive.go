package cmd

import (
	"fmt"
	"log"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(ctf *setup.Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	return reactiveRun(ctf, seed)
}

func reactiveRun(ctf *setup.Config, seed int) (pair *orderedmap.OrderedMap[string, any], err error) {
	// 出力フォルダがない場合、作成
	err = createDirIfNotExist(ctf.OutfilePath)
	if err != nil {
		return pair, err
	}
	// 実行
	out, timeout, err := reactiveRunCmd(ctf, seed)
	if err != nil {
		log.Println(err)
	}
	// 出力のパース
	pair = orderedmap.NewOrderedMap[string, any]()
	pair.Set("seed", seed)
	keys, err := ExtractKeyValuePairs(pair, string(out))
	_ = keys // ordermapを消す時に使う
	if err != nil {
		return pair, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerDate, err := extractData((string(out)))
	if err != nil {
		return pair, err
	}
	for k, v := range testerDate {
		pair.Set(k, v)
	}
	pair.Set("stdErr", out)
	pair.Set("result", timeout)
	if timeout {
		pair.Set("Score", 0)
		pair.Set("time", float64(ctf.TimeLimitMS/1000))
	}
	return pair, nil
}

// reactiveRunCmd はreactive=trueのときに使う
func reactiveRunCmd(ctf *setup.Config, seed int) ([]byte, bool, error) {
	cmd := setup.LanguageSets[ctf.Language].ExeCmd
	infile := ctf.InfilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := ctf.OutfilePath + fmt.Sprintf("%04d.txt", seed)
	setsArgs := setArgs(ctf.Args) // コマンドオプションの追加
	cmdStr := fmt.Sprintf("%s %s %s < %s > %s", ctf.TesterPath, cmd, setsArgs, infile, outfile)
	cmdStrings := createCommand(cmdStr)
	out, timeout, err := runCommandWithTimeout(cmdStrings, int(ctf.TimeLimitMS))
	if err != nil {
		log.Println("Error: ", err, "command:", cmdStr)
	}
	return out, timeout, err
}

// setArgs return string
func setArgs(args []string) string {
	var str string
	for _, v := range args {
		str += v + " "
	}
	return str
}
