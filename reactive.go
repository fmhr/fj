package fj

import (
	"fmt"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(ctf *Config, seed int) (map[string]float64, error) {
	rtn, err := reactiveRun(ctf, seed)
	if err != nil {
		return rtn, ErrorTrace("failed to run: %v", err)
	}
	//fmt.Fprintln(os.Stderr, mapString(rtn))
	//fmt.Println(rtn["Score"]) // ここだけ標準出力
	return rtn, nil
}

func reactiveRun(ctf *Config, seed int) (map[string]float64, error) {
	out, err := reactiveRunCmd(ctf, seed)
	if err != nil {
		return nil, err
	}
	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return nil, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerDate, err := extractData((string(out)))
	if err != nil {
		return nil, err
	}
	for k, v := range testerDate {
		pair[k] = v
	}
	pair["seed"] = float64(seed)
	return pair, nil
}

func reactiveRunCmd(ctf *Config, seed int) ([]byte, error) {
	infile := ctf.InfilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := ctf.OutfilePath + fmt.Sprintf("%04d.txt", seed)
	setsArgs := setArgs(ctf.Args)
	cmdStr := fmt.Sprintf("%s %s %s < %s > %s", ctf.TesterPath, ctf.Cmd, setsArgs, infile, outfile)
	cmd := createCommand(cmdStr)
	out, err := runCommandWithTimeout(cmd, ctf)
	if err != nil {
		return nil, fmt.Errorf("command [%q] failed with: %v out: %v", cmd, err, out)
	}
	return out, nil
}

func setArgs(args []string) string {
	var str string
	for _, v := range args {
		str += v + " "
	}
	return str
}
