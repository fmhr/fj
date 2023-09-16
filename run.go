package fj

import (
	"fmt"
	"os"
	"os/exec"
)

// Run runs the program with the given seed
func Run(cnf *config, seed int) ([]byte, error) {
	inputfile := fmt.Sprintf("%s%04d.txt", cnf.InfilePath, seed)
	if _, err := os.Stat(inputfile); err != nil {
		return []byte{}, fmt.Errorf("input file [%s] does not exist", inputfile)
	}
	outputfile := fmt.Sprintf("%s%04d.out", cnf.OutfilePath, seed)
	cmdStr := cnf.Cmd + " < " + inputfile + " > " + outputfile
	cmd := exec.Command("sh", "-c", cmdStr)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return []byte{}, fmt.Errorf("cmd.Run() for command [%q] failed with: %v", cmd, err)
	}
	return out, nil
}
