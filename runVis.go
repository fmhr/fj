package fj

import (
	"fmt"
	"path/filepath"
	"sort"
)

func RunVis(cnf *Config, seed int) (map[string]float64, error) {
	return runVis(cnf, seed)
}

// runVis は指定された設定とシードに基づいてコマンドを実行して、
// その結果をvisに渡して、両方の結果を返す
func runVis(cnf *Config, seed int) (map[string]float64, error) {
	out, err := normalRun(cnf, seed)
	if err != nil {
		return nil, fmt.Errorf("failed running with seed%d: %v", seed, err)
	}
	//log.Println("run:", string(out))

	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return pair, fmt.Errorf("ExtractKeyValuePairs failed: %v", err)
	}
	// vis
	infile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.out", seed))

	outVis, err := vis(cnf, infile, outfile)
	if err != nil {
		return nil, fmt.Errorf("vis failed: %v", err)
	}

	sc, err := extractData(string(outVis))
	if err != nil {
		return nil, fmt.Errorf("extractData failed: %v", err)
	}

	for k, v := range sc {
		pair[k] = v
	}
	pair["seed"] = float64(seed)
	return pair, nil
}

func mapString(data map[string]float64) string {
	var str string
	str += fmt.Sprintf("seed=%d ", int(data["seed"]))
	str += fmt.Sprintf("Score=%.2f ", data["Score"])
	orderKey := make([]string, 0)
	for k := range data {
		if k != "seed" && k != "Score" {
			orderKey = append(orderKey, k)
		}
	}
	sort.Strings(orderKey)
	for _, k := range orderKey {
		str += fmt.Sprintf("%s=%v ", k, data[k])
	}
	return str
}

func vis(cnf *Config, infile, outfile string) ([]byte, error) {
	cmdStr := fmt.Sprintf(cnf.VisPath+" %s %s", infile, outfile)
	cmd := createComand(cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err)
	}
	return out, nil
}
