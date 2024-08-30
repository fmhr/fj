package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
)

func RunVis(cnf *setup.Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	return runVis(cnf, seed)
}

// runVis は指定された設定とシードに基づいてコマンドを実行して、
// その結果をvisに渡して、両方の結果を返す
// 通常の問題（reactive=false)で使う
func runVis(cnf *setup.Config, seed int) (pair *orderedmap.OrderedMap[string, any], err error) {
	out, _, err := normalRun(cnf, seed)
	if err != nil {
		log.Println("normalRun error:", err)
		return nil, WrapError(fmt.Errorf("%w\nout: %s", err, string(out)))
	}
	pair = orderedmap.NewOrderedMap[string, any]()
	pair.Set("seed", seed)

	keys, err := ExtractKeyValuePairs(pair, string(out))
	if err != nil {
		return pair, err
	}
	_ = keys // Ordermapを消す時に使う
	// vis
	infile := filepath.Join(cnf.InfilePath, fmt.Sprintf("%04d.txt", seed))
	outfile := filepath.Join(cnf.OutfilePath, fmt.Sprintf("%04d.txt", seed))

	outVis, err := vis(cnf, infile, outfile)
	if err != nil {
		log.Println("vis error:", err)
		return nil, err
	}
	sc, err := extractData(string(outVis))
	if err != nil {
		return nil, err
	}

	// Score=0の場合は出力を表示
	if sc["Score"] == 0 {
		log.Println("Score=0 out:", string(outVis))
	}

	for k, v := range sc {
		pair.Set(k, v)
	}
	pair.Set("seed", seed)
	return pair, nil
}

// vis is a wrapper for vis command
func vis(cnf *setup.Config, infile, outfile string) ([]byte, error) {
	cmdStr := fmt.Sprintf(cnf.VisPath+" %s %s", infile, outfile)
	cmdStrings := createCommand(cmdStr)
	cmd := exec.Command(cmdStrings[0], cmdStrings[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("failed: %v\nout: %s cmd: %s\n", err, string(out), cmdStrings)
		return nil, WrapError(fmt.Errorf("failed: %v\nout: %s cmd: %s", err, string(out), cmdStrings))
	}
	return out, nil
}
