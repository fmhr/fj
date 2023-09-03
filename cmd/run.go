package main

import (
	"fmt"
	"log"
	"os/exec"
)

// a.exe < 0000.in > out.out
func Run(seed int) ([]int, error) {
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
	return []int{0}, nil
}
