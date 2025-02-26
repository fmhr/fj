package cmd

import (
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"

	"github.com/elliotchance/orderedmap/v2"
	"github.com/fmhr/fj/cmd/setup"
)

func RunVis(cnf *setup.Config, seed int) (*orderedmap.OrderedMap[string, string], error) {
	return runVis(cnf, seed)
}

// runVis は指定された設定とシードに基づいてコマンドを実行して、
// その結果をvisに渡して、両方の結果を返す
// 通常の問題（reactive=false)で使う
func runVis(cnf *setup.Config, seed int) (pair *orderedmap.OrderedMap[string, string], err error) {
	out, _, err := normalRun(cnf, seed)
	if err != nil {
		//log.Println("Error: ", err)
		if len(out) > 0 {
			log.Println("out:", string(out))
		}
		return nil, WrapError(fmt.Errorf("%w", err))
	}
	pair = orderedmap.NewOrderedMap[string, string]()
	pair.Set("seed", fmt.Sprintf("%d", seed))

	keys, err := ExtractKeyValuePairs(pair, string(out))
	if err != nil {
		return pair, err
	}
	_ = keys // Ordermapを消す時に使う

	// メモリ使用量を正規表現で抽出
	// 数字の後に"maximum resident set size"が続くパターンを検索
	re := regexp.MustCompile(`\s*(\d+)\s+maximum resident set size`)
	// ステータス出力からメモリ使用量を抽出
	matches := re.FindStringSubmatch(string(out))
	if len(matches) > 1 {
		mb, err := strconv.Atoi(matches[1])
		if err != nil {
			return nil, err
		}
		mb /= 1024 * 1024
		pair.Set("Memory", fmt.Sprintf("%d", mb))
	} else {
		log.Println("Memory usage not found.")
	}

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
	if sc["Score"] == "0" {
		log.Println("Score=0 out:", string(outVis))
	}

	for k, v := range sc {
		pair.Set(k, v)
	}
	pair.Set("seed", fmt.Sprintf("%d", seed))
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
