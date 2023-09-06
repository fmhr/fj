package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
)

var CMD = "bin/main"
var TESTER = "./tools/target/release/tester"
var VIS = "./tools/target/release/vis"
var OUTFILE = "out.txt"
var INFILE_FOLDER = "tools/in/"

func Run(seed int) ([]byte, error) {
	if _, err := os.Stat(fmt.Sprintf("tools/in/%04d.txt", seed)); err != nil {
		return []byte{}, fmt.Errorf("tools/in/%04d.txt not found", seed)
	}
	cmdStr := fmt.Sprintf(CMD+" < tools/in/%04d.txt > out.txt", seed)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err)
	}
	return out, nil
}

func RunVis(seed int) ([]byte, error) {
	out, err := Run(seed)
	if err != nil {
		return nil, err
	}
	infile := INFILE_FOLDER + fmt.Sprintf("%04d.txt", seed)
	outfile := OUTFILE
	outVis := vis(infile, outfile)
	out = append(out, outVis...)
	//log.Print(string(outVis))
	pair, err := ExtractKeyValuePairs(string(out))
	if err != nil {
		return []byte{}, err
	}
	pair["seed"] = float64(seed)
	log.Println(pair)
	return out, nil
}

func vis(infile, outfile string) []byte {
	cmdStr := fmt.Sprintf(VIS+" %s %s", infile, outfile)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err))
	}
	return out
}

func tester10() error {
	result := make([][]byte, 10)
	for seed := 0; seed < 10; seed++ {
		r, err := RunVis(seed)
		if err != nil {
			return err
		}
		result[seed] = r
	}
	var sumScore int
	for i := 0; i < 10; i++ {
		data, err := ExtractKeyValuePairs(string(result[i]))
		if err != nil {
			return err
		}
		sumScore += int(data["score"])
	}
	//table := tablewriter.NewWriter(os.Stdout)
	//table.SetHeader([]string{"seed", "N", "L", "S", "Score", "Hit", "Miss", "Placement cost", "Measurement cost", "count"})
	//for _, v := range result {
	//table.Append([]string{fmt.Sprintf("%d", v[0]), fmt.Sprintf("%d", v[1]), fmt.Sprintf("%d", v[2]), fmt.Sprintf("%d", v[3]), fmt.Sprintf("%d", v[4]), fmt.Sprintf("%d", v[1]-v[5]), fmt.Sprintf("%d", v[5]), fmt.Sprintf("%d", v[6]), fmt.Sprintf("%d", v[7]), fmt.Sprintf("%d", v[8])})
	//sumScore += v[4]
	//}
	//table.Render()
	log.Println("sumScore=", sumScore)
	return nil
}
