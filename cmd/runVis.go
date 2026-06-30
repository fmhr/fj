package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/fmhr/fj/cmd/setup"
)

func RunVis(cnf *setup.Config, seed int) (SliceMap, error) {
	return runVis(cnf, seed)
}

// runVis は指定された設定とシードに基づいてコマンドを実行して、
// その結果をvisに渡して、両方の結果を返す
// 通常の問題（reactive=false)で使う
func runVis(cnf *setup.Config, seed int) (kv SliceMap, err error) {
	out, timeout, maxRSS, err := normalRun(cnf, seed)
	if err != nil {
		//log.Println("Error: ", err)
		if len(out) > 0 {
			log.Println("out:", string(out))
		}
		return nil, fmt.Errorf("%w", err)
	}
	kv = NewSliceMap()
	kv.Set("seed", fmt.Sprintf("%d", seed))

	if timeout {
		kv.Set("result", "TLE")
		kv.Set("Score", "0")
		kv.Set("time", fmt.Sprintf("%d", cnf.TimeLimitMS))
		return kv, nil
	}

	keys, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return kv, err
	}
	for _, v := range keys {
		kv.Set(v.Key, v.Val)
	}

	// メモリ使用量を結果に追加
	if maxRSS > 0 {
		kv.Set("Memory", fmt.Sprintf("%d", maxRSS))
	}

	// vis
	infile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	outVis, err := vis(cnf, infile, outfile)
	if err != nil {
		return nil, fmt.Errorf("visの実行に失敗: %w", err)
	}

	// testerで要素を出力するように改造したときに、ここで値を取得する
	kvv, err := ExtractKeyValuePairs(string(outVis))
	if err != nil {
		return nil, fmt.Errorf("visの出力からキーと値のペアを抽出に失敗: %w", err)
	}
	for _, v := range kvv {
		kv.Set(v.Key, v.Val)
	}

	// Score = を抜き出す
	sc, err := extractData(string(outVis))
	if err != nil {
		return nil, err
	}

	// Score=0の場合は出力を表示
	if sc["Score"] == "0" {
		log.Println("Score=0 out:", string(outVis))
	}

	for k, v := range sc {
		kv.Set(k, v)
	}
	kv.Set("seed", fmt.Sprintf("%d", seed))
	return kv, nil
}

// vis is a wrapper for vis command
func vis(cnf *setup.Config, infile, outfile string) ([]byte, error) {
	// visコマンドが存在しない場合は、scoreコマンドを探す
	if !exists(cnf.VisPath) && exists("tools/target/release/score") {
		cnf.VisPath = "tools/target/release/score"
	}
	cmdStr := fmt.Sprintf(cnf.VisPath+" %s %s", infile, outfile)
	cmdStrings := createCommand(cmdStr)
	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed: %v\nout: %s cmd: %s\n", err, string(out), cmdStrings)
		return nil, fmt.Errorf("failed: %v\nout: %s cmd: %s", err, string(out), cmdStrings)
	}
	return out, nil
}

func exists(path string) bool {
	_, err := exec.LookPath(path)
	return err == nil
}
