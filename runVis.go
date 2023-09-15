package fj

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"sort"
)

// RunVis runs the program with the given seed and visualize the result
func runVis(cnf *config, seed int) (map[string]float64, error) {
	out, err := Run(cnf, seed)
	if err != nil {
		return nil, err
	}
	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return pair, err
	}
	// vis
	infile := cnf.infilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := cnf.outfilePath + fmt.Sprintf("%04d.out", seed)
	outVis := vis(cnf, infile, outfile)
	// score
	sc, err := extractScore(string(string(outVis)))
	if err != nil {
		log.Fatal(err)
	}
	pair["TesterScore"] = float64(sc)
	pair["seed"] = float64(seed)
	return pair, nil
}

func RunVis(cnf *config, seed int) error {
	rtn, err := runVis(cnf, seed)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, mapString(rtn))
	fmt.Println(rtn["TesterScore"]) // ここだけ標準出力
	return nil
}

func mapString(data map[string]float64) string {
	var str string
	str += fmt.Sprintf("seed=%d ", int(data["seed"]))
	str += fmt.Sprintf("Score=%.2f ", data["TesterScore"])
	orderKey := make([]string, 0)
	for k := range data {
		if k != "seed" && k != "TesterScore" {
			orderKey = append(orderKey, k)
		}
	}
	sort.Strings(orderKey)
	for _, k := range orderKey {
		str += fmt.Sprintf("%s=%v ", k, data[k])
	}
	return str
}

func vis(cnf *config, infile, outfile string) []byte {
	cmdStr := fmt.Sprintf(cnf.visPath+" %s %s", infile, outfile)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err))
	}
	return out
}

func RunVis10(cnf *config) error {
	var sumScore int
	for seed := 0; seed < 10; seed++ {
		r, err := runVis(cnf, seed)
		if err != nil {
			return err
		}
		// fmt.Fprintln(os.Stderr, mapString(r))
		fmt.Fprintln(os.Stderr, mapString(r))
		sumScore += int(r["TesterScore"])
	}
	fmt.Fprintln(os.Stderr, "sumScore=", sumScore)
	fmt.Println(sumScore) // ここだけ標準出力
	return nil
}
