package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fmhr/fj/cmd/setup"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(cnf *setup.Config, seed int) (kv SliceMap, err error) {
	// 出力フォルダがない場合、作成
	err = ensureDirExists(cnf.OutfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}
	// 実行
	out, timeout, maxRSS, err := reactiveRunCmd(cnf, seed)
	if err != nil {
		return nil, fmt.Errorf("reactiveRunCmdの実行に失敗: %w", err)
	}
	// 出力のパース
	kv.Set("seed", fmt.Sprintf("%d", seed))
	keys, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return kv, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerData, err := extractData((string(out)))
	if err != nil {
		return kv, fmt.Errorf("extractDataの実行に失敗: %w", err)
	}
	for _, v := range keys {
		kv.Set(v.Key, v.Val)
	}

	for k, v := range testerData {
		kv.Set(k, v)
	}
	if timeout {
		kv.Set("result", "TimeOut")
		kv.Set("Score", "0")
		kv.Set("time", fmt.Sprintf("%d", cnf.TimeLimitMS))
	}

	// メモリ使用量を結果に追加
	if maxRSS > 0 {
		kv.Set("Memory", fmt.Sprintf("%d", maxRSS))
	}

	return kv, nil
}

// reactiveRunCmd はreactive=trueのときに使う
func reactiveRunCmd(cnf *setup.Config, seed int) ([]byte, bool, int64, error) {
	cmd := cnf.ExecuteCmd
	infile := cnf.InfilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := cnf.OutfilePath + fmt.Sprintf("%04d.txt", seed)
	// in/0000.txtの存在確認
	_, err := fileExists(infile)
	if err != nil {
		return nil, false, 0, fmt.Errorf("input fileの確認に失敗: %w", err)
	}
	// outfilePathが存在するか確認
	_, err = fileExists(cnf.OutfilePath)
	if err != nil {
		return nil, false, 0, fmt.Errorf("output directoryの確認に失敗: %w", err)
	}
	// testerPathが存在するか確認
	_, err = fileExists(cnf.TesterPath)
	if err != nil {
		return nil, false, 0, fmt.Errorf("tester fileの確認に失敗: %w", err)
	}
	// コマンド作成
	setsArgs := setArgs(cnf.Args) // コマンドオプションの追加
	cmdStr := fmt.Sprintf("%s %s %s < %s > %s", cnf.TesterPath, cmd, setsArgs, infile, outfile)
	cmdStrings := createCommand(cmdStr)

	// 実行時間を計測
	startTime := time.Now()

	// コマンド実行
	out, timeout, maxRSS, err := runCommandWithTimeout(cmdStrings, int(cnf.TimeLimitMS))
	elapsed := time.Since(startTime)

	if err != nil {
		return out, timeout, maxRSS, fmt.Errorf("commandの実行に失敗: %w cmd: %s stderr: %s", err, cmdStr, string(out))
	}

	// 実行時間を出力に追加（ミリ秒単位）
	if !timeout {
		timeStr := fmt.Sprintf("time=%d", elapsed.Milliseconds())
		out = append(out, []byte("\n"+timeStr)...)
	}

	return out, timeout, maxRSS, nil
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
