package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/fmhr/fj/cmd/setup"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(ctf *setup.Config, seed int) (kv SliceMap, err error) {
	// 出力フォルダがない場合、作成
	err = createDirIfNotExist(ctf.OutfilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to create output directory: %v", err)
	}
	// 実行
	out, timeout, err := reactiveRunCmd(ctf, seed)
	if err != nil {
		return nil, fmt.Errorf("reactiveRunCmdの実行に失敗: %w", err)
	}
	// 出力のパース
	kv.Set("seed", fmt.Sprintf("%d", seed))
	keys, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return kv, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerDate, err := extractData((string(out)))
	if err != nil {
		return kv, fmt.Errorf("extractDataの実行に失敗: %w", err)
	}
	for _, v := range keys {
		kv.Set(v.key, v.val)
	}

	for k, v := range testerDate {
		kv.Set(k, v)
	}
	if timeout {
		kv.Set("Score", "0")
		kv.Set("time", fmt.Sprintf("%v", ctf.TimeLimitMS/1000))
	}
	return kv, nil
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
	// outfilePathが存在するか確認
	_, err = fileExists(ctf.OutfilePath)
	if err != nil {
		return nil, false, fmt.Errorf("output directoryの確認に失敗: %w", err)
	}
	// testerPathが存在するか確認
	_, err = fileExists(ctf.TesterPath)
	if err != nil {
		return nil, false, fmt.Errorf("tester fileの確認に失敗: %w", err)
	}
	// コマンド作成
	setsArgs := setArgs(ctf.Args) // コマンドオプションの追加
	cmdStr := fmt.Sprintf("%s %s %s < %s > %s", ctf.TesterPath, cmd, setsArgs, infile, outfile)
	cmdStrings := createCommand(cmdStr)

	// 実行時間を計測
	startTime := time.Now()

	// コマンド実行
	out, timeout, err := runCommandWithTimeout(cmdStrings, int(ctf.TimeLimitMS))
	elapsed := time.Since(startTime)

	if err != nil {
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
