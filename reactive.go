package fj

import (
	"fmt"

	"github.com/elliotchance/orderedmap/v2"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(ctf *Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	rtn, err := reactiveRun(ctf, seed)
	if err != nil {
		return nil, err
	}
	//fmt.Fprintln(os.Stderr, mapString(rtn))
	//fmt.Println(rtn["Score"]) // ここだけ標準出力
	return &rtn, nil
}

func reactiveRun(ctf *Config, seed int) (orderedmap.OrderedMap[string, any], error) {
	err := createDirIfNotExist(ctf.OutfilePath)
	if err != nil {
		return orderedmap.OrderedMap[string, any]{}, err
	}
	out, err := reactiveRunCmd(ctf, seed)
	if err != nil {
		return orderedmap.OrderedMap[string, any]{}, err
	}
	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return orderedmap.OrderedMap[string, any]{}, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerDate, err := extractData((string(out)))
	if err != nil {
		return orderedmap.OrderedMap[string, any]{}, err
	}
	for k, v := range testerDate {
		pair.Set(k, v)
	}
	pair.Set("seed", seed)
	pair.Set("stdErr", out)
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
