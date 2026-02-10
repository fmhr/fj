package cmd

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(ctf *setup.Config, seed int) (pair *orderedmap.OrderedMap[string, string], err error) {
	// 出力フォルダがない場合、作成
	err = createDirIfNotExist(ctf.OutfilePath)
	if err != nil {
		return pair, fmt.Errorf("failed to create output directory: %v", err)
	}
	// 実行
	out, timeout, err := reactiveRunCmd(ctf, seed)
	if err != nil {
		return nil, fmt.Errorf("reactiveRunCmdの実行に失敗: %w", err)
	}
	// 出力のパース
	pair = orderedmap.NewOrderedMap[string, string]()
	pair.Set("seed", fmt.Sprintf("%d", seed))
	keys, err := ExtractKeyValuePairs(pair, string(out))
	_ = keys // ordermapを消す時に使う
	if err != nil {
		return pair, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerDate, err := extractData((string(out)))
	if err != nil {
		return pair, fmt.Errorf("extractDataの実行に失敗: %w", err)
	}
	for k, v := range testerDate {
		pair.Set(k, v)
	}
	if timeout {
		pair.Set("Score", "0")
		pair.Set("time", fmt.Sprintf("%v", ctf.TimeLimitMS/1000))
	}
	return pair, nil
}

// reactiveRunCmd はreactive=trueのときに使う
func reactiveRunCmd(ctf *setup.Config, seed int) ([]byte, bool, error) {
	cmd := ctf.ExecuteCmd
	infile := ctf.InfilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := ctf.OutfilePath + fmt.Sprintf("%04d.txt", seed)
	// in/0000.txtの存在確認
	_, err := fileExists(infile)
	if err != nil {
		return nil, false, fmt.Errorf("input fileの確認に失敗: %w", err)
	}
	setsArgs := setArgs(ctf.Args) // コマンドオプションの追加
	cmdStr := fmt.Sprintf("%s %s %s < %s > %s", ctf.TesterPath, cmd, setsArgs, infile, outfile)
	cmdStrings := createCommand(cmdStr)

	// 実行時間を計測
	startTime := time.Now()
	out, timeout, err := runCommandWithTimeout(cmdStrings, int(ctf.TimeLimitMS))
	elapsed := time.Since(startTime)

	if err != nil {
		log.Println("Command output:", string(out))
		return out, timeout, fmt.Errorf("commandの実行に失敗: %w cmd: %s stderr: %s", err, cmdStr, string(out))
	}

	// 実行時間を出力に追加（秒単位）
	if !timeout {
		timeStr := fmt.Sprintf("time=%.3f", elapsed.Seconds())
		out = append(out, []byte("\n"+timeStr)...)
	}

	return out, timeout, nil
}

// setArgs return string
func setArgs(args []string) string {
	var str string
	for _, v := range args {
		str += v + " "
	}
	return str
}

func fileExists(filename string) (bool, error) {
	_, err := os.Stat(filename)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, fmt.Errorf("file does not exist: %s", filename)
	}
	return false, err
}
