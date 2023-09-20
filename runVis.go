package fj

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
)

func RunVis(cnf *Config, seed int) error {
	rtn, err := runVis(cnf, seed)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, mapString(rtn))
	fmt.Println(rtn["Score"]) // ここだけ標準出力
	return nil
}

// runVis は指定された設定とシードに基づいてコマンドを実行して、
// その結果をvisに渡して、両方の結果を返す
func runVis(cnf *Config, seed int) (map[string]float64, error) {
	out, err := Run(cnf, seed)
	if err != nil {
		return nil, err
	}
	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return pair, err
	}
	// vis
	infile := cnf.InfilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := cnf.OutfilePath + fmt.Sprintf("%04d.out", seed)
	outVis := vis(cnf, infile, outfile)
	// visの結果をpairに追加
	sc, err := extractData(string(outVis))
	if err != nil {
		log.Fatal(err)
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

func vis(cnf *Config, infile, outfile string) []byte {
	cmdStr := fmt.Sprintf(cnf.VisPath+" %s %s", infile, outfile)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err))
	}
	return out
}

func RunVis10(cnf *Config) error {
	var sumScore int
	for seed := 0; seed < 10; seed++ {
		r, err := runVis(cnf, seed)
		if err != nil {
			return err
		}
		// fmt.Fprintln(os.Stderr, mapString(r))
		fmt.Fprintln(os.Stderr, mapString(r))
		sumScore += int(r["Score"])
	}
	fmt.Fprintln(os.Stderr, "sumScore=", sumScore)
	fmt.Println(sumScore) // ここだけ標準出力
	return nil
}
