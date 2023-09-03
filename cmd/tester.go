package main

import (
	"fmt"
	"log"
	"os/exec"
)

func tester(seed int) ([]int, error) {
	log.Println("seed=", seed)
	err := build()
	if err != nil {
		return nil, err
	}
	cmdStr := fmt.Sprintf("bin/a < tools/in/%04d.txt > out.txt", seed)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return nil, err
	}
	fmt.Printf("------------------------ seed=%d\n", seed)
	fmt.Print(string(out))
	fmt.Println("------------------------")
	//r := make([]int, len(regexStrs)+1)
	//r[0] = seed
	//for i, regexStr := range regexStrs {
	//r[i+1] = parseInt(string(out), regexStr.re, regexStr.str)
	//}
	return []int{0}, nil
}

func tester10() error {
	result := make([][]int, 10)
	for seed := 0; seed < 10; seed++ {
		r, err := tester(seed)
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
