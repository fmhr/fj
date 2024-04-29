package fj

import (
	"fmt"
	"log"

	"github.com/elliotchance/orderedmap/v2"
)

// ReactiveRun はreactive=trueのときに使う
func ReactiveRun(ctf *Config, seed int) (*orderedmap.OrderedMap[string, any], error) {
	return reactiveRun(ctf, seed)
}

func reactiveRun(ctf *Config, seed int) (pair *orderedmap.OrderedMap[string, any], err error) {
	err = createDirIfNotExist(ctf.OutfilePath)
	if err != nil {
		return pair, err
	}
	out, result, err := reactiveRunCmd(ctf, seed)
	if err != nil {
		log.Println(err)
	}
	pair = orderedmap.NewOrderedMap[string, any]()
	pair.Set("seed", seed)
	keys, err := ExtractKeyValuePairs(pair, string(out))
	_ = keys // ordermapを消す時に使う
	if err != nil {
		return pair, fmt.Errorf("failed to extract key-value pairs: %v, source: %s", err, string(out))
	}
	testerDate, err := extractData((string(out)))
	if err != nil {
		return pair, err
	}
	for k, v := range testerDate {
		pair.Set(k, v)
	}
	pair.Set("stdErr", out)
	pair.Set("result", result)
	if result == "TLE" {
		pair.Set("Score", 0)
		pair.Set("time", float64(ctf.TimeLimitMS/1000))
	}
	return pair, nil
}

// reactiveRunCmd はreactive=trueのときに使う
func reactiveRunCmd(ctf *Config, seed int) ([]byte, string, error) {
	cmd := LanguageSets[ctf.Language].ExeCmd
	infile := ctf.InfilePath + fmt.Sprintf("%04d.txt", seed)
	outfile := ctf.OutfilePath + fmt.Sprintf("%04d.txt", seed)
	setsArgs := setArgs(ctf.Args)
	cmdStr := fmt.Sprintf("%s %s %s < %s > %s", ctf.TesterPath, cmd, setsArgs, infile, outfile)
	cmdStrings := createCommand(cmdStr)
	out, result, err := runCommandWithTimeout(cmdStrings, int(ctf.TimeLimitMS))
	//ctx, cancel := context.WithTimeout(context.Background(), time.Duration(ctf.TimeLimitMS)*time.Millisecond)
	//defer cancel()
	//cmd := exec.CommandContext(ctx, cmdStrings[0], cmdStrings[1:]...)
	//out, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("Error: ", err, "command:", cmdStr)
		//return nil, fmt.Errorf("command [%q] failed with: %v out: %v", cmd, err, out)
	}
	return out, result, err
}

func setArgs(args []string) string {
	var str string
	for _, v := range args {
		str += v + " "
	}
	return str
}
