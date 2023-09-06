package main

import (
	"fmt"
	"log"
	"os/exec"
)

var CMD = "bin/main"
var TESTER = "./tools/target/release/tester"
var VIS = "./tools/target/release/vis"
var OUTFILE = "out.txt"
var INFILE_FOLDER = "tools/in/"

func RunTester(seed int) ([]byte, error) {
	log.Println("seed=", seed)
	cmdStr := fmt.Sprintf(CMD+" < tools/in/%04d.txt > out.txt", seed)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err)
	}

	fmt.Printf("------------------------ seed=%d\n", seed)
	fmt.Print(string(out))
	fmt.Println("------------------------")
	//r := make([]int, len(regexStrs)+1)
	//r[0] = seed
	//for i, regexStr := range regexStrs {
	//r[i+1] = parseInt(string(out), regexStr.re, regexStr.str)
	//}
	infile := INFILE_FOLDER + fmt.Sprintf("%04d.txt", seed)
	outfile := OUTFILE
	vis(infile, outfile)
	return out, nil
}

func vis(infile, outfile string) {
	cmdStr := fmt.Sprintf(VIS+" %s %s", infile, outfile)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(fmt.Errorf("cmd.Run() for command %q failed with: %v", cmdStr, err))
	}
	fmt.Println(string(out))
}

func tester10() error {
	result := make([][]byte, 10)
	for seed := 0; seed < 10; seed++ {
		r, err := RunTester(seed)
		if err != nil {
			return err
		}
		result[seed] = r
	}
	var sumScore int
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
