package main

import (
	"fmt"
	"os"
	"os/exec"
)

// Run runs the program with the given seed
func Run(cnf *config, seed int) ([]byte, error) {
	if _, err := os.Stat(fmt.Sprintf("tools/in/%04d.txt", seed)); err != nil {
		return []byte{}, fmt.Errorf("tools/in/%04d.txt not found", seed)
	}
	cmdStr := fmt.Sprintf(cnf.Cmd+" < tools/in/%04d.txt > tmp/%04d.out", seed, seed)
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	return out, nil
}
