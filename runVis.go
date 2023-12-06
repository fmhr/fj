package fj

import (
	"fmt"
	"path/filepath"
)

func RunVis(cnf *Config, seed int) (map[string]float64, error) {
	return runVis(cnf, seed)
}

// runVis は指定された設定とシードに基づいてコマンドを実行して、
// その結果をvisに渡して、両方の結果を返す
// 通常の問題（reactive=false)で使う
func runVis(cnf *Config, seed int) (map[string]float64, error) {
	out, err := normalRun(cnf, seed)
	if err != nil {
		return nil, ErrorTrace("falied: normalRun", err)
	}

	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return pair, ErrorTrace("failed:ExtractKeyValuePairs", err)
	}
	// vis
	infile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	outVis, err := vis(cnf, infile, outfile)
	if err != nil {
		//return nil, TraceMsg(fmt.Errorf("failed: %v", err).Error())
		return nil, ErrorTrace("failed:vis", err)
	}
	sc, err := extractData(string(outVis))
	if err != nil {
		return nil, ErrorTrace("failed:extractData", err)
	}

	for k, v := range sc {
		pair[k] = v
	}
	pair["seed"] = float64(seed)
	return pair, nil
}

// vis is a wrapper for vis command
func vis(cnf *Config, infile, outfile string) ([]byte, error) {
	cmdStr := fmt.Sprintf(cnf.VisPath+" %s %s", infile, outfile)
	cmd := createCommand(cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, ErrorTrace(fmt.Sprintf("failed: command %s", cmdStr), err)
	}
	return out, nil
}
